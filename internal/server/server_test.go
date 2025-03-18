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

func TestSessionManagement(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("Gerenciador de sessões é nil")
	}

	// Testa adição de sessão
	session := &mockSession{
		id:   "test1",
		user: "test",
		addr: "localhost:2222",
	}

	s := sm.AddSession(session)
	if s == nil {
		t.Fatal("Sessão não foi criada")
	}

	if s.ID != session.id {
		t.Errorf("ID da sessão incorreto: %s != %s", s.ID, session.id)
	}

	// Testa recuperação de sessão
	s, ok := sm.GetSession(session.id)
	if !ok {
		t.Fatal("Sessão não encontrada")
	}

	if s.ID != session.id {
		t.Errorf("ID da sessão incorreto: %s != %s", s.ID, session.id)
	}

	// Testa remoção de sessão
	sm.RemoveSession(session.id)
	_, ok = sm.GetSession(session.id)
	if ok {
		t.Error("Sessão não foi removida")
	}
}

func TestCommandExecution(t *testing.T) {
	sm := NewSessionManager()
	session := &mockSession{
		id:   "test1",
		user: "test",
		addr: "localhost:2222",
	}

	s := sm.AddSession(session)
	if s == nil {
		t.Fatal("Sessão não foi criada")
	}

	// Testa execução de comando
	result, err := sm.ExecuteCommand(session.id, "echo test")
	if err != nil {
		t.Fatalf("Erro ao executar comando: %v", err)
	}

	if result == nil {
		t.Fatal("Resultado é nil")
	}

	if result.Command != "echo test" {
		t.Errorf("Comando incorreto: %s != %s", result.Command, "echo test")
	}

	if result.Output != "test\n" {
		t.Errorf("Saída incorreta: %s != %s", result.Output, "test\n")
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
