package server

import (
	"testing"
)

func TestCommandManager(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Testa execução de comando
	result, err := manager.ExecuteCommand(sessionInfo.ID, "echo 'test'")
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
	result, err = manager.ExecuteCommand(sessionInfo.ID, "invalid_command")
	if err == nil {
		t.Error("Erro esperado ao executar comando inválido")
	}
}

func TestCommandManagerConcurrent(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Testa execução concorrente de comandos
	done := make(chan bool)
	concurrent := 10

	for i := 0; i < concurrent; i++ {
		go func(id int) {
			cmd := "echo 'test'"
			result, err := manager.ExecuteCommand(sessionInfo.ID, cmd)
			if err != nil {
				t.Errorf("Erro ao executar comando %d: %v", id, err)
				done <- true
				return
			}
			if result == nil {
				t.Errorf("Resultado do comando %d é nil", id)
				done <- true
				return
			}
			if result.Output != "test\n" {
				t.Errorf("Saída incorreta do comando %d: %s != %s", id, result.Output, "test\n")
			}
			if result.Status != 0 {
				t.Errorf("Status incorreto do comando %d: %d != 0", id, result.Status)
			}
			done <- true
		}(i)
	}

	// Aguarda todas as goroutines terminarem
	for i := 0; i < concurrent; i++ {
		<-done
	}
}

func TestCommandManagerHistory(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Executa vários comandos
	commands := []string{"echo 'test1'", "echo 'test2'", "echo 'test3'"}
	for _, cmd := range commands {
		result, err := manager.ExecuteCommand(sessionInfo.ID, cmd)
		if err != nil {
			t.Fatalf("Erro ao executar comando: %v", err)
		}
		if result == nil {
			t.Fatal("Resultado é nil")
		}
	}

	// Obtém o histórico de comandos
	history, err := manager.GetCommandHistory(sessionInfo.ID)
	if err != nil {
		t.Fatalf("Erro ao obter histórico de comandos: %v", err)
	}

	// Verifica se todos os comandos estão no histórico
	if len(history) != len(commands) {
		t.Errorf("Tamanho do histórico incorreto: %d != %d", len(history), len(commands))
	}

	for i, cmd := range commands {
		if history[i].Command != cmd {
			t.Errorf("Comando incorreto no histórico: %s != %s", history[i].Command, cmd)
		}
	}
}

func TestCommandManagerSystemInfo(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Obtém informações do sistema
	info, err := manager.GetSystemInfo(sessionInfo.ID)
	if err != nil {
		t.Fatalf("Erro ao obter informações do sistema: %v", err)
	}

	// Verifica se as informações necessárias estão presentes
	requiredFields := []string{"id", "user", "hostname", "os", "arch", "connected"}
	for _, field := range requiredFields {
		if value, ok := info[field]; !ok || value == "" {
			t.Errorf("Campo %s está vazio ou ausente", field)
		}
	}
}

func TestCommandManagerInvalidSession(t *testing.T) {
	manager := NewSessionManager()

	// Testa execução de comando em sessão inválida
	result, err := manager.ExecuteCommand("invalid_session", "echo 'test'")
	if err == nil {
		t.Error("Erro esperado ao executar comando em sessão inválida")
	}
	if result != nil {
		t.Error("Resultado não deveria existir para sessão inválida")
	}

	// Testa obtenção de histórico em sessão inválida
	history, err := manager.GetCommandHistory("invalid_session")
	if err == nil {
		t.Error("Erro esperado ao obter histórico de sessão inválida")
	}
	if history != nil {
		t.Error("Histórico não deveria existir para sessão inválida")
	}

	// Testa obtenção de informações do sistema em sessão inválida
	info, err := manager.GetSystemInfo("invalid_session")
	if err == nil {
		t.Error("Erro esperado ao obter informações de sessão inválida")
	}
	if info != nil {
		t.Error("Informações não deveriam existir para sessão inválida")
	}
}
