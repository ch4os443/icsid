package server

import (
	"io"
	"net"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/icsid/icsid/internal/config"
)

func TestServerCreation(t *testing.T) {
	cfg := &config.Config{}
	cfg.Server.SSH.Port = 2222
	cfg.Server.SSH.HostKey = "test_key"
	cfg.Server.SSH.Username = "test"
	cfg.Server.SSH.Password = "test"
	cfg.Server.Web.Port = 8443
	cfg.Server.Web.CertFile = "test_cert.pem"
	cfg.Server.Web.KeyFile = "test_key.pem"

	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar servidor: %v", err)
	}

	if srv == nil {
		t.Fatal("Servidor é nil")
	}

	if srv.config != cfg {
		t.Error("Configuração não foi corretamente atribuída")
	}

	if srv.sessions == nil {
		t.Error("Gerenciador de sessões é nil")
	}
}

func TestServerInitialization(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			SSH struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"ssh"`
			Web struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			} `yaml:"web"`
		}{
			SSH: struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Port:     2222,
				HostKey:  "testdata/host_key",
				Username: "test",
				Password: "test123",
			},
			Web: struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			}{
				Port:     8080,
				CertFile: "testdata/cert.pem",
				KeyFile:  "testdata/key.pem",
			},
		},
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar servidor: %v", err)
	}
	if server == nil {
		t.Fatal("Servidor é nil")
	}
	if server.ssh == nil {
		t.Fatal("Servidor SSH é nil")
	}
	if server.sessions == nil {
		t.Fatal("Gerenciador de sessões é nil")
	}
	if server.logger == nil {
		t.Fatal("Logger é nil")
	}
	if server.rateLimiter == nil {
		t.Fatal("Rate limiter é nil")
	}
}

func TestSessionManagement(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			SSH struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"ssh"`
			Web struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			} `yaml:"web"`
		}{
			SSH: struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Port:     2222,
				HostKey:  "testdata/host_key",
				Username: "test",
				Password: "test123",
			},
			Web: struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			}{
				Port:     8080,
				CertFile: "testdata/cert.pem",
				KeyFile:  "testdata/key.pem",
			},
		},
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar servidor: %v", err)
	}

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := server.sessions.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Verifica se a sessão foi adicionada
	session, ok := server.sessions.GetSession(sessionInfo.ID)
	if !ok {
		t.Fatal("Sessão não encontrada")
	}
	if session.ID != sessionInfo.ID {
		t.Errorf("ID da sessão incorreto: %s != %s", session.ID, sessionInfo.ID)
	}

	// Remove a sessão
	server.sessions.RemoveSession(sessionInfo.ID)
	session, ok = server.sessions.GetSession(sessionInfo.ID)
	if ok {
		t.Error("Sessão não foi removida")
	}
}

func TestCommandExecution(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			SSH struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"ssh"`
			Web struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			} `yaml:"web"`
		}{
			SSH: struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Port:     2222,
				HostKey:  "testdata/host_key",
				Username: "test",
				Password: "test123",
			},
			Web: struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			}{
				Port:     8080,
				CertFile: "testdata/cert.pem",
				KeyFile:  "testdata/key.pem",
			},
		},
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar servidor: %v", err)
	}

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := server.sessions.AddSession(mockSession)
	defer server.sessions.RemoveSession(sessionInfo.ID)

	// Testa execução de comando
	result, err := server.sessions.ExecuteCommand(sessionInfo.ID, "echo 'test'")
	if err != nil {
		t.Fatalf("Erro ao executar comando: %v", err)
	}
	if result == nil {
		t.Fatal("Resultado é nil")
	}
	if result.Output != "test\n" {
		t.Errorf("Saída incorreta: %s != %s", result.Output, "test\n")
	}
	if result.Status != 0 {
		t.Errorf("Status incorreto: %d != 0", result.Status)
	}

	// Testa comando inválido
	result, err = server.sessions.ExecuteCommand(sessionInfo.ID, "invalid_command")
	if err == nil {
		t.Error("Erro esperado ao executar comando inválido")
	}
}

func TestSystemInfo(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			SSH struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"ssh"`
			Web struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			} `yaml:"web"`
		}{
			SSH: struct {
				Port     int    `yaml:"port"`
				HostKey  string `yaml:"host_key"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Port:     2222,
				HostKey:  "testdata/host_key",
				Username: "test",
				Password: "test123",
			},
			Web: struct {
				Port     int    `yaml:"port"`
				CertFile string `yaml:"cert_file"`
				KeyFile  string `yaml:"key_file"`
			}{
				Port:     8080,
				CertFile: "testdata/cert.pem",
				KeyFile:  "testdata/key.pem",
			},
		},
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar servidor: %v", err)
	}

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := server.sessions.AddSession(mockSession)
	defer server.sessions.RemoveSession(sessionInfo.ID)

	// Obtém informações do sistema
	info, err := server.sessions.GetSystemInfo(sessionInfo.ID)
	if err != nil {
		t.Fatalf("Erro ao obter informações do sistema: %v", err)
	}
	if info == nil {
		t.Fatal("Informações do sistema são nil")
	}

	// Verifica campos obrigatórios
	requiredFields := []string{"hostname", "os", "architecture"}
	for _, field := range requiredFields {
		if value, ok := info[field]; !ok || value == "" {
			t.Errorf("Campo %s está vazio ou ausente", field)
		}
	}
}

// Mock para testes
type mockSession struct {
	id   string
	user string
	addr string
}

func (m *mockSession) ID() string {
	return m.id
}

func (m *mockSession) User() string {
	return m.user
}

func (m *mockSession) RemoteAddr() net.Addr {
	return nil
}

func (m *mockSession) ReadCommand() (string, error) {
	return "", nil
}

func (m *mockSession) WriteString(s string) error {
	return nil
}

func (m *mockSession) Close() error {
	return nil
}

func (m *mockSession) Break(c chan<- bool) {
}

func (m *mockSession) Context() ssh.Context {
	return nil
}

func (m *mockSession) SetContext(ctx interface{}) {
}

func (m *mockSession) Pty() (ssh.Pty, <-chan ssh.Window, bool) {
	return ssh.Pty{}, nil, false
}

func (m *mockSession) SetPty(term string, width, height int) error {
	return nil
}

func (m *mockSession) Signal(signal string) error {
	return nil
}

func (m *mockSession) Exit(code int) error {
	return nil
}

func (m *mockSession) CloseWrite() error {
	return nil
}

func (m *mockSession) CloseRead() error {
	return nil
}

func (m *mockSession) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (m *mockSession) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockSession) Command() []string {
	return []string{}
}

func (m *mockSession) Environ() []string {
	return []string{}
}

func (m *mockSession) LocalAddr() net.Addr {
	return nil
}

func (m *mockSession) Permissions() ssh.Permissions {
	return ssh.Permissions{}
}

func (m *mockSession) PublicKey() ssh.PublicKey {
	return nil
}

func (m *mockSession) RawCommand() string {
	return ""
}

func (m *mockSession) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return false, nil
}

func (m *mockSession) Signals(c chan<- ssh.Signal) {
}

func (m *mockSession) Stderr() io.ReadWriter {
	return nil
}

func (m *mockSession) Subsystem() string {
	return ""
}
