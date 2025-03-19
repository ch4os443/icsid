package server

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ch4os443/icsid/internal/config"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type Server struct {
	config      *config.Config
	ssh         *ssh.Server
	web         *WebServer
	sessions    *SessionManager
	logger      *Logger
	rateLimiter *RateLimiter
	wg          sync.WaitGroup
}

// ExtendedSession adiciona funcionalidades ao ssh.Session
type ExtendedSession struct {
	ssh.Session
	reader *bufio.Reader
}

func (s *ExtendedSession) ReadCommand() (string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config:      cfg,
		sessions:    NewSessionManager(),
		rateLimiter: NewRateLimiter(5*time.Minute, 5), // 5 tentativas em 5 minutos
	}

	// Inicializa o logger
	logger, err := NewLogger("logs/ssh.log", INFO)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar logger: %v", err)
	}
	s.logger = logger

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

	// Inicia a limpeza do rate limiter
	go s.cleanupRateLimiter()

	return s, nil
}

func (s *Server) cleanupRateLimiter() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.rateLimiter.Cleanup()
	}
}

func (s *Server) Start() error {
	s.logger.Info("Iniciando servidor ICSID na porta %d", s.config.Server.SSH.Port)

	// Inicia o servidor SSH
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.ssh.ListenAndServe(); err != nil {
			s.logger.Error("Erro no servidor SSH: %v", err)
		}
	}()

	// Inicia o servidor web
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.web.Start(); err != nil {
			s.logger.Error("Erro no servidor web: %v", err)
		}
	}()

	s.wg.Wait()
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("Encerrando servidor ICSID")
	if err := s.ssh.Close(); err != nil {
		s.logger.Error("Erro ao encerrar servidor SSH: %v", err)
	}
	if err := s.logger.Close(); err != nil {
		s.logger.Error("Erro ao encerrar logger: %v", err)
	}
	return nil
}

func (s *Server) handlePassword(ctx ssh.Context, password string) bool {
	addr := ctx.RemoteAddr().String()
	if !s.rateLimiter.Allow(addr) {
		s.logger.LogSecurityEvent("Rate limit excedido", fmt.Sprintf("Endereço: %s", addr))
		return false
	}

	success := ctx.User() == s.config.Server.SSH.Username &&
		password == s.config.Server.SSH.Password

	s.logger.LogConnection(addr, ctx.User(), success)
	return success
}

func (s *Server) handleSession(session ssh.Session) {
	// Cria uma sessão estendida
	extSession := &ExtendedSession{
		Session: session,
		reader:  bufio.NewReader(session),
	}

	// Registra a sessão
	sessionInfo := s.sessions.AddSession(extSession)
	defer s.sessions.RemoveSession(sessionInfo.ID)

	s.logger.Info("Nova sessão SSH - ID: %s, Usuário: %s, Endereço: %s",
		sessionInfo.ID, sessionInfo.User, sessionInfo.Hostname)

	// Envia informações do sistema
	systemInfo, err := s.sessions.GetSystemInfo(sessionInfo.ID)
	if err != nil {
		s.logger.Error("Erro ao obter informações do sistema: %v", err)
		return
	}

	infoJSON, err := json.Marshal(systemInfo)
	if err != nil {
		s.logger.Error("Erro ao serializar informações do sistema: %v", err)
		return
	}

	io.WriteString(session, fmt.Sprintf("Sistema: %s\n", string(infoJSON)))

	// Processa comandos
	for {
		cmd, err := extSession.ReadCommand()
		if err != nil {
			if err != io.EOF {
				s.logger.Error("Erro ao ler comando: %v", err)
			}
			return
		}

		result, err := s.sessions.ExecuteCommand(sessionInfo.ID, cmd)
		if err != nil {
			s.logger.Error("Erro ao executar comando: %v", err)
			io.WriteString(session, fmt.Sprintf("Erro ao executar comando: %v\n", err))
			continue
		}

		s.logger.LogCommand(sessionInfo.ID, sessionInfo.User, cmd, true)
		io.WriteString(session, fmt.Sprintf("Resultado: %s\n", result.Output))
	}
}

func (s *Server) loadOrGenerateHostKey() (ssh.Signer, error) {
	keyPath := s.config.Server.SSH.HostKey

	// Tenta carregar a chave existente
	if keyData, err := os.ReadFile(keyPath); err == nil {
		key, err := gossh.ParsePrivateKey(keyData)
		if err == nil {
			s.logger.Info("Chave SSH carregada com sucesso")
			return key, nil
		}
	}

	s.logger.Info("Gerando nova chave SSH")

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

	// Converte para o formato SSH
	signer, err := gossh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Nova chave SSH gerada e salva com sucesso")
	return signer, nil
}
