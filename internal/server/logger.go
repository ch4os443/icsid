package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	file     *os.File
	mu       sync.Mutex
	level    LogLevel
	filePath string
}

func NewLogger(path string, level LogLevel) (*Logger, error) {
	// Cria o diretório de logs se não existir
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de logs: %v", err)
	}

	// Abre o arquivo de log
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}

	return &Logger{
		file:     file,
		level:    level,
		filePath: path,
	}, nil
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	levelStr := "DEBUG"
	switch level {
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, levelStr, message)

	if _, err := l.file.WriteString(logLine); err != nil {
		fmt.Printf("Erro ao escrever log: %v\n", err)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) LogConnection(addr, user string, success bool) {
	status := "falhou"
	if success {
		status = "sucesso"
	}
	l.Info("Tentativa de conexão - Endereço: %s, Usuário: %s, Status: %s", addr, user, status)
}

func (l *Logger) LogCommand(sessionID, user, cmd string, success bool) {
	status := "falhou"
	if success {
		status = "sucesso"
	}
	l.Info("Comando executado - Sessão: %s, Usuário: %s, Comando: %s, Status: %s", sessionID, user, cmd, status)
}

func (l *Logger) LogSecurityEvent(event, details string) {
	l.Warn("Evento de segurança - %s: %s", event, details)
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
