package server

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Cria um diretório temporário para os logs
	tmpDir := "testdata/logs"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Cria um novo logger
	logger, err := NewLogger(filepath.Join(tmpDir, "test.log"), DEBUG)
	if err != nil {
		t.Fatalf("Erro ao criar logger: %v", err)
	}
	defer logger.Close()

	// Testa diferentes níveis de log
	logger.Debug("Mensagem de debug")
	logger.Info("Mensagem de info")
	logger.Warn("Mensagem de aviso")
	logger.Error("Mensagem de erro")

	// Testa logs específicos
	logger.LogConnection("127.0.0.1", "test", true)
	logger.LogCommand("session1", "test", "ls", true)
	logger.LogError(errors.New("Erro de teste"), "Detalhes do erro")
	logger.LogSecurityEvent("Evento de segurança", "Detalhes do evento")

	// Aguarda um momento para garantir que os logs foram escritos
	time.Sleep(100 * time.Millisecond)

	// Verifica se o arquivo de log foi criado
	logFile := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("Arquivo de log não foi criado: %v", err)
	}

	// Lê o conteúdo do arquivo de log
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Erro ao ler arquivo de log: %v", err)
	}

	// Verifica se todas as mensagens estão presentes
	expectedMessages := []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"Connection",
		"Command",
		"Error",
		"Security",
	}

	for _, msg := range expectedMessages {
		if !bytes.Contains(content, []byte(msg)) {
			t.Errorf("Mensagem '%s' não encontrada no log", msg)
		}
	}
}

func TestLoggerLevels(t *testing.T) {
	tmpDir := "testdata/logs"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Testa diferentes níveis de log
	levels := []LogLevel{DEBUG, INFO, WARN, ERROR}
	for _, level := range levels {
		logger, err := NewLogger(filepath.Join(tmpDir, "test.log"), level)
		if err != nil {
			t.Fatalf("Erro ao criar logger com nível %v: %v", level, err)
		}

		// Registra mensagens em todos os níveis
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Warn("Warn message")
		logger.Error("Error message")

		logger.Close()

		// Verifica o conteúdo do log
		content, err := os.ReadFile(filepath.Join(tmpDir, "test.log"))
		if err != nil {
			t.Fatalf("Erro ao ler arquivo de log: %v", err)
		}

		// Verifica se as mensagens foram registradas de acordo com o nível
		switch level {
		case DEBUG:
			if !bytes.Contains(content, []byte("DEBUG")) {
				t.Error("Mensagem DEBUG não encontrada no nível DEBUG")
			}
			fallthrough
		case INFO:
			if !bytes.Contains(content, []byte("INFO")) {
				t.Error("Mensagem INFO não encontrada no nível INFO")
			}
			fallthrough
		case WARN:
			if !bytes.Contains(content, []byte("WARN")) {
				t.Error("Mensagem WARN não encontrada no nível WARN")
			}
			fallthrough
		case ERROR:
			if !bytes.Contains(content, []byte("ERROR")) {
				t.Error("Mensagem ERROR não encontrada no nível ERROR")
			}
		}
	}
}

func TestLoggerFileRotation(t *testing.T) {
	tmpDir := "testdata/logs"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logger, err := NewLogger(filepath.Join(tmpDir, "test.log"), DEBUG)
	if err != nil {
		t.Fatalf("Erro ao criar logger: %v", err)
	}
	defer logger.Close()

	// Gera muitos logs para forçar a rotação
	for i := 0; i < 1000; i++ {
		logger.Info("Mensagem de teste %d", i)
	}

	// Aguarda um momento para garantir que os logs foram escritos
	time.Sleep(100 * time.Millisecond)

	// Verifica se o arquivo de log foi criado
	logFile := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("Arquivo de log não foi criado: %v", err)
	}

	// Verifica se o arquivo de backup foi criado
	backupFile := filepath.Join(tmpDir, "test.log.bak")
	if _, err := os.Stat(backupFile); err != nil {
		t.Fatalf("Arquivo de backup não foi criado: %v", err)
	}
}

func TestLoggerConcurrent(t *testing.T) {
	tmpDir := "testdata/logs"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logger, err := NewLogger(filepath.Join(tmpDir, "test.log"), DEBUG)
	if err != nil {
		t.Fatalf("Erro ao criar logger: %v", err)
	}
	defer logger.Close()

	// Testa logs concorrentes
	done := make(chan bool)
	concurrent := 10

	for i := 0; i < concurrent; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("Mensagem concorrente %d-%d", id, j)
			}
			done <- true
		}(i)
	}

	// Aguarda todas as goroutines terminarem
	for i := 0; i < concurrent; i++ {
		<-done
	}

	// Aguarda um momento para garantir que os logs foram escritos
	time.Sleep(100 * time.Millisecond)

	// Verifica se o arquivo de log foi criado
	logFile := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("Arquivo de log não foi criado: %v", err)
	}

	// Lê o conteúdo do arquivo de log
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Erro ao ler arquivo de log: %v", err)
	}

	// Verifica se todas as mensagens foram registradas
	for i := 0; i < concurrent; i++ {
		for j := 0; j < 100; j++ {
			msg := fmt.Sprintf("Mensagem concorrente %d-%d", i, j)
			if !bytes.Contains(content, []byte(msg)) {
				t.Errorf("Mensagem '%s' não encontrada no log", msg)
			}
		}
	}
}
