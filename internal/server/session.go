package server

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type SessionInfo struct {
	ID       string
	User     string
	Hostname string
	Start    time.Time
	LastCmd  time.Time
}

type CommandResult struct {
	Output string
	Error  error
}

type SessionManager struct {
	sessions map[string]*SessionInfo
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*SessionInfo),
	}
}

func (sm *SessionManager) AddSession(session *ExtendedSession) *SessionInfo {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	info := &SessionInfo{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		User:     session.User(),
		Hostname: session.RemoteAddr().String(),
		Start:    time.Now(),
		LastCmd:  time.Now(),
	}

	sm.sessions[info.ID] = info
	return info
}

func (sm *SessionManager) RemoveSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)
}

func (sm *SessionManager) GetSession(id string) *SessionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.sessions[id]
}

func (sm *SessionManager) ListSessions() []*SessionInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*SessionInfo, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func (sm *SessionManager) ExecuteCommand(sessionID string, cmd string) (*CommandResult, error) {
	session := sm.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	session.LastCmd = time.Now()

	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}

	output, err := command.CombinedOutput()
	if err != nil {
		return &CommandResult{
			Output: string(output),
			Error:  err,
		}, nil
	}

	return &CommandResult{
		Output: string(output),
		Error:  nil,
	}, nil
}

func (sm *SessionManager) GetSystemInfo(sessionID string) (map[string]interface{}, error) {
	session := sm.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	info := make(map[string]interface{})

	// Informações do sistema operacional
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH
	info["hostname"] = session.Hostname

	// Informações de hardware
	if runtime.GOOS == "windows" {
		// Windows
		if output, err := exec.Command("wmic", "cpu", "get", "name").Output(); err == nil {
			info["cpu"] = string(output)
		}
		if output, err := exec.Command("wmic", "memorychip", "get", "capacity").Output(); err == nil {
			info["memory"] = string(output)
		}
	} else {
		// Linux/macOS
		if output, err := exec.Command("lscpu").Output(); err == nil {
			info["cpu"] = string(output)
		}
		if output, err := exec.Command("free", "-h").Output(); err == nil {
			info["memory"] = string(output)
		}
	}

	// Informações de rede
	if output, err := exec.Command("ip", "addr").Output(); err == nil {
		info["network"] = string(output)
	}

	// Informações de disco
	if runtime.GOOS == "windows" {
		if output, err := exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption").Output(); err == nil {
			info["disk"] = string(output)
		}
	} else {
		if output, err := exec.Command("df", "-h").Output(); err == nil {
			info["disk"] = string(output)
		}
	}

	return info, nil
}

func (sm *SessionManager) GetSessionStats(sessionID string) (map[string]interface{}, error) {
	session := sm.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("sessão não encontrada: %s", sessionID)
	}

	stats := make(map[string]interface{})
	stats["id"] = session.ID
	stats["user"] = session.User
	stats["hostname"] = session.Hostname
	stats["start_time"] = session.Start
	stats["last_command"] = session.LastCmd
	stats["duration"] = time.Since(session.Start)

	return stats, nil
}
