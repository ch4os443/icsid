package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/icsid/icsid/internal/config"
)

type Server struct {
	config   *config.Config
	ssh      *ssh.Server
	web      *WebServer
	sessions *SessionManager
	wg       sync.WaitGroup
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config:   cfg,
		sessions: NewSessionManager(),
	}

	// Gera ou carrega a chave SSH
	hostKey, err := s.loadOrGenerateHostKey()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar/gerar chave SSH: %v", err)
	}

	// Configura o servidor SSH
	s.ssh = &ssh.Server{
		Addr:            fmt.Sprintf(":%d", cfg.Server.SSH.Port),
		HostSigners:     []ssh.Signer{hostKey},
		PasswordHandler: s.handlePassword,
		Handler:         s.handleSession,
	}

	// Configura o servidor web
	s.web = NewWebServer(cfg, s.sessions)

	return s, nil
}

func (s *Server) Start() error {
	// Inicia o servidor SSH
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Iniciando servidor SSH na porta %d", s.config.Server.SSH.Port)
		if err := s.ssh.ListenAndServe(); err != nil {
			log.Printf("Erro no servidor SSH: %v", err)
		}
	}()

	// Inicia o servidor web
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Iniciando servidor web na porta %d", s.config.Server.Web.Port)
		if err := s.web.Start(); err != nil {
			log.Printf("Erro no servidor web: %v", err)
		}
	}()

	// Aguarda todos os servidores
	s.wg.Wait()
	return nil
}

func (s *Server) Stop() error {
	// Encerra o servidor SSH
	if err := s.ssh.Close(); err != nil {
		log.Printf("Erro ao encerrar servidor SSH: %v", err)
	}

	// TODO: Implementar encerramento do servidor web

	return nil
}

func (s *Server) handlePassword(ctx ssh.Context, password string) bool {
	return ctx.User() == s.config.Server.SSH.Username &&
		password == s.config.Server.SSH.Password
}

func (s *Server) handleSession(session ssh.Session) {
	// Registra a sessão
	sessionInfo := s.sessions.AddSession(session)
	defer s.sessions.RemoveSession(sessionInfo.ID)

	// Envia informações do sistema
	systemInfo, err := s.sessions.GetSystemInfo(sessionInfo.ID)
	if err != nil {
		log.Printf("Erro ao obter informações do sistema: %v", err)
		return
	}

	infoJSON, err := json.Marshal(systemInfo)
	if err != nil {
		log.Printf("Erro ao serializar informações do sistema: %v", err)
		return
	}

	io.WriteString(session, fmt.Sprintf("Sistema: %s\n", string(infoJSON)))

	// Processa comandos
	for {
		cmd, err := session.ReadCommand()
		if err != nil {
			if err != io.EOF {
				log.Printf("Erro ao ler comando: %v", err)
			}
			return
		}

		result, err := s.sessions.ExecuteCommand(sessionInfo.ID, cmd)
		if err != nil {
			io.WriteString(session, fmt.Sprintf("Erro ao executar comando: %v\n", err))
			continue
		}

		io.WriteString(session, fmt.Sprintf("Resultado: %s\n", result.Output))
	}
}

func (s *Server) loadOrGenerateHostKey() (ssh.Signer, error) {
	keyPath := s.config.Server.SSH.HostKey

	// Tenta carregar a chave existente
	if keyData, err := os.ReadFile(keyPath); err == nil {
		key, err := ssh.ParsePrivateKey(keyData)
		if err == nil {
			return key, nil
		}
	}

	// Gera uma nova chave se não existir
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Cria o diretório se não existir
	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return nil, err
	}

	// Salva a chave
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	// Converte a chave para formato PEM
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	if _, err := keyFile.Write(keyPEM); err != nil {
		return nil, err
	}

	// Retorna a chave como ssh.Signer
	return ssh.NewSignerFromKey(key)
}
