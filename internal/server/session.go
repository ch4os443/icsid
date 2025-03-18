package server

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/gliderlabs/ssh"
)

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

type Session struct {
	ID        string
	User      string
	Hostname  string
	OS        string
	Arch      string
	Connected bool
	Commands  []Command
}

type Command struct {
	ID        string
	Command   string
	Output    string
	Status    int
	Timestamp int64
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) AddSession(session ssh.Session) *Session {
	s := &Session{
		ID:        session.ID(),
		User:      session.User(),
		Hostname:  session.RemoteAddr().String(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Connected: true,
		Commands:  make([]Command, 0),
	}

	sm.mu.Lock()
	sm.sessions[s.ID] = s
	sm.mu.Unlock()

	return s
}

func (sm *SessionManager) RemoveSession(id string) {
	sm.mu.Lock()
	delete(sm.sessions, id)
	sm.mu.Unlock()
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	s, ok := sm.sessions[id]
	sm.mu.RUnlock()
	return s, ok
}

func (sm *SessionManager) ListSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, s := range sm.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

func (sm *SessionManager) ExecuteCommand(sessionID, command string) (*Command, error) {
	s, ok := sm.GetSession(sessionID)
	if !ok {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	cmdResult := Command{
		ID:        fmt.Sprintf("%d", len(s.Commands)+1),
		Command:   command,
		Output:    string(output),
		Status:    cmd.ProcessState.ExitCode(),
		Timestamp: cmd.ProcessState.SystemTime().Milliseconds(),
	}

	s.Commands = append(s.Commands, cmdResult)
	return &cmdResult, nil
}

func (sm *SessionManager) GetSystemInfo(sessionID string) (map[string]interface{}, error) {
	s, ok := sm.GetSession(sessionID)
	if !ok {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	info := map[string]interface{}{
		"id":        s.ID,
		"user":      s.User,
		"hostname":  s.Hostname,
		"os":        s.OS,
		"arch":      s.Arch,
		"connected": s.Connected,
	}

	return info, nil
}

func (sm *SessionManager) GetCommandHistory(sessionID string) ([]Command, error) {
	s, ok := sm.GetSession(sessionID)
	if !ok {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	return s.Commands, nil
}
