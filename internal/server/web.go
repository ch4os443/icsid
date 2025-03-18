package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/icsid/icsid/internal/config"
)

type WebServer struct {
	config   *config.Config
	sessions *SessionManager
	router   *mux.Router
}

func NewWebServer(cfg *config.Config, sessions *SessionManager) *WebServer {
	ws := &WebServer{
		config:   cfg,
		sessions: sessions,
		router:   mux.NewRouter(),
	}

	// Configura as rotas
	ws.setupRoutes()

	return ws
}

func (ws *WebServer) setupRoutes() {
	// Servir arquivos estáticos
	ws.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Servir a página principal
	ws.router.HandleFunc("/", ws.handleIndex)

	// APIs
	ws.router.HandleFunc("/api/sessions", ws.handleListSessions).Methods("GET")
	ws.router.HandleFunc("/api/system-info/{id}", ws.handleSystemInfo).Methods("GET")
	ws.router.HandleFunc("/api/execute", ws.handleExecuteCommand).Methods("POST")
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "index.html"))
}

func (ws *WebServer) handleListSessions(w http.ResponseWriter, r *http.Request) {
	sessions := ws.sessions.ListSessions()
	json.NewEncoder(w).Encode(sessions)
}

func (ws *WebServer) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	info, err := ws.sessions.GetSystemInfo(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(info)
}

func (ws *WebServer) handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
		Command   string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := ws.sessions.ExecuteCommand(req.SessionID, req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (ws *WebServer) Start() error {
	addr := fmt.Sprintf(":%d", ws.config.Server.Web.Port)
	log.Printf("Iniciando servidor web na porta %d", ws.config.Server.Web.Port)
	return http.ListenAndServeTLS(addr, ws.config.Server.Web.CertFile, ws.config.Server.Web.KeyFile, ws.router)
}
