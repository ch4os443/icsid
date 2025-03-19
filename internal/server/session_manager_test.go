package server

import (
	"testing"
	"time"
)

func TestSessionManager(t *testing.T) {
	// Cria um novo gerenciador de sessões
	manager := NewSessionManager()

	// Cria uma sessão mock
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Testa recuperação de sessão
	session, ok := manager.GetSession(sessionInfo.ID)
	if !ok {
		t.Fatal("Sessão não encontrada")
	}
	if session.ID != sessionInfo.ID {
		t.Errorf("ID da sessão incorreto: %s != %s", session.ID, sessionInfo.ID)
	}

	// Testa remoção de sessão
	manager.RemoveSession(sessionInfo.ID)
	session, ok = manager.GetSession(sessionInfo.ID)
	if ok {
		t.Error("Sessão não foi removida")
	}
}

func TestSessionManagerConcurrent(t *testing.T) {
	manager := NewSessionManager()
	done := make(chan bool)
	concurrent := 10

	// Testa operações concorrentes
	for i := 0; i < concurrent; i++ {
		go func(id int) {
			// Cria sessão
			mockSession := &ExtendedSession{}
			sessionInfo := manager.AddSession(mockSession)
			if sessionInfo.ID == "" {
				t.Errorf("ID da sessão %d está vazio", id)
				done <- true
				return
			}

			// Recupera sessão
			session, ok := manager.GetSession(sessionInfo.ID)
			if !ok {
				t.Errorf("Sessão %d não encontrada", id)
				done <- true
				return
			}
			if session.ID != sessionInfo.ID {
				t.Errorf("ID da sessão incorreto: %s != %s", session.ID, sessionInfo.ID)
				done <- true
				return
			}

			// Remove sessão
			manager.RemoveSession(sessionInfo.ID)
			session, ok = manager.GetSession(sessionInfo.ID)
			if ok {
				t.Errorf("Sessão %d não foi removida", id)
			}

			done <- true
		}(i)
	}

	// Aguarda todas as goroutines terminarem
	for i := 0; i < concurrent; i++ {
		<-done
	}
}

func TestSessionManagerCleanup(t *testing.T) {
	manager := NewSessionManager()

	// Cria várias sessões
	sessions := make([]string, 5)
	for i := 0; i < 5; i++ {
		mockSession := &ExtendedSession{}
		sessionInfo := manager.AddSession(mockSession)
		sessions[i] = sessionInfo.ID
	}

	// Remove algumas sessões
	manager.RemoveSession(sessions[1])
	manager.RemoveSession(sessions[3])

	// Verifica se as sessões foram removidas corretamente
	if _, ok := manager.GetSession(sessions[1]); ok {
		t.Error("Sessão 1 não foi removida")
	}
	if _, ok := manager.GetSession(sessions[3]); ok {
		t.Error("Sessão 3 não foi removida")
	}

	// Verifica se as outras sessões ainda existem
	if _, ok := manager.GetSession(sessions[0]); !ok {
		t.Error("Sessão 0 não existe mais")
	}
	if _, ok := manager.GetSession(sessions[2]); !ok {
		t.Error("Sessão 2 não existe mais")
	}
	if _, ok := manager.GetSession(sessions[4]); !ok {
		t.Error("Sessão 4 não existe mais")
	}
}

func TestSessionManagerTimeout(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Simula passagem de tempo
	time.Sleep(2 * time.Second)

	// Verifica se a sessão ainda existe
	session, ok := manager.GetSession(sessionInfo.ID)
	if !ok {
		t.Error("Sessão expirou muito cedo")
	}
	if session.ID != sessionInfo.ID {
		t.Errorf("ID da sessão incorreto: %s != %s", session.ID, sessionInfo.ID)
	}
}

func TestSessionManagerCommandExecution(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão
	mockSession := &ExtendedSession{}
	sessionInfo := manager.AddSession(mockSession)
	if sessionInfo.ID == "" {
		t.Fatal("ID da sessão está vazio")
	}

	// Executa um comando válido
	result, err := manager.ExecuteCommand(sessionInfo.ID, "ls")
	if err != nil {
		t.Fatalf("Erro ao executar comando: %v", err)
	}
	if result == nil {
		t.Fatal("Resultado é nil")
	}
	if result.Output == "" {
		t.Error("Saída do comando está vazia")
	}
	if result.Status != 0 {
		t.Errorf("Status do comando incorreto: %d", result.Status)
	}

	// Executa um comando inválido
	result, err = manager.ExecuteCommand(sessionInfo.ID, "comando_invalido")
	if err == nil {
		t.Error("Erro esperado ao executar comando inválido")
	}
}

func TestSessionManagerSystemInfo(t *testing.T) {
	manager := NewSessionManager()

	// Cria uma sessão
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
