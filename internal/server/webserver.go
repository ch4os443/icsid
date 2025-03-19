package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ch4os443/icsid/internal/config"
	"github.com/gorilla/mux"
)

type WebServer struct {
	config   *config.Config
	sessions *SessionManager
	router   *mux.Router
	server   *http.Server
	mu       sync.RWMutex
}

func NewWebServer(cfg *config.Config, sessions *SessionManager) *WebServer {
	ws := &WebServer{
		config:   cfg,
		sessions: sessions,
		router:   mux.NewRouter(),
	}

	ws.setupRoutes()
	return ws
}

func (ws *WebServer) setupRoutes() {
	// Rotas da API
	api := ws.router.PathPrefix("/api").Subrouter()
	api.Use(ws.authMiddleware)

	// Listar sessões
	api.HandleFunc("/sessions", ws.handleListSessions).Methods("GET")

	// Obter informações de uma sessão
	api.HandleFunc("/sessions/{id}", ws.handleGetSession).Methods("GET")

	// Executar comando
	api.HandleFunc("/sessions/{id}/execute", ws.handleExecuteCommand).Methods("POST")

	// Obter informações do sistema
	api.HandleFunc("/sessions/{id}/system", ws.handleGetSystemInfo).Methods("GET")

	// Obter estatísticas da sessão
	api.HandleFunc("/sessions/{id}/stats", ws.handleGetSessionStats).Methods("GET")

	// Servir arquivos estáticos
	ws.router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/static")))
}

func (ws *WebServer) Start() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", ws.config.Server.Web.Port),
		Handler: ws.router,
	}

	return ws.server.ListenAndServeTLS(
		ws.config.Server.Web.CertFile,
		ws.config.Server.Web.KeyFile,
	)
}

func (ws *WebServer) Stop() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.server != nil {
		return ws.server.Close()
	}
	return nil
}

func (ws *WebServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="ICSID"`)
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		if username != ws.config.Server.SSH.Username ||
			password != ws.config.Server.SSH.Password {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ws *WebServer) handleListSessions(w http.ResponseWriter, r *http.Request) {
	sessions := ws.sessions.ListSessions()
	json.NewEncoder(w).Encode(sessions)
}

func (ws *WebServer) handleGetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	session := ws.sessions.GetSession(vars["id"])
	if session == nil {
		http.Error(w, "Sessão não encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func (ws *WebServer) handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cmd struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Comando inválido", http.StatusBadRequest)
		return
	}

	result, err := ws.sessions.ExecuteCommand(vars["id"], cmd.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (ws *WebServer) handleGetSystemInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	info, err := ws.sessions.GetSystemInfo(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (ws *WebServer) handleGetSessionStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stats, err := ws.sessions.GetSessionStats(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
