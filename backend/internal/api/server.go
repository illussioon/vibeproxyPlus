package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"vibeproxyplus/backend/internal/config"
	"vibeproxyplus/backend/internal/proxy"
)

type App struct {
	mu      sync.Mutex
	store   *config.Store
	proxy   *proxy.Manager
	state   config.State
	authDir string
}

func NewServer() (*http.Server, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	baseDir = filepath.Join(baseDir, "VibeProxyPlus")
	store, err := config.NewStore(baseDir)
	if err != nil {
		return nil, err
	}
	state, err := store.Load()
	if err != nil {
		return nil, err
	}
	authDir := filepath.Join(baseDir, "auth")
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		return nil, err
	}

	app := &App{store: store, proxy: proxy.NewManager(), state: state, authDir: authDir}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/state", app.handleState)
	mux.HandleFunc("/api/server/start", app.handleStart)
	mux.HandleFunc("/api/server/stop", app.handleStop)
	mux.HandleFunc("/api/settings/launch-at-login", app.handleLaunchAtLogin)
	mux.HandleFunc("/api/auth/open", app.handleOpenAuth)

	return &http.Server{Addr: "127.0.0.1:7890", Handler: withCORS(mux)}, nil
}

func (a *App) handleState(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	a.state.Running = a.proxy.Running()
	state := a.state
	a.mu.Unlock()
	writeJSON(w, state)
}

func (a *App) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := a.proxy.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	a.state.Running = true
	_ = a.store.Save(a.state)
	a.mu.Unlock()
	writeJSON(w, map[string]bool{"ok": true})
}

func (a *App) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := a.proxy.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	a.state.Running = false
	_ = a.store.Save(a.state)
	a.mu.Unlock()
	writeJSON(w, map[string]bool{"ok": true})
}

func (a *App) handleLaunchAtLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.mu.Lock()
	a.state.LaunchAtLogin = body.Enabled
	_ = a.store.Save(a.state)
	a.mu.Unlock()
	writeJSON(w, map[string]bool{"ok": true})
}

func (a *App) handleOpenAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if runtime.GOOS == "windows" {
		_ = exec.Command("explorer", a.authDir).Start()
	}
	writeJSON(w, map[string]string{"path": a.authDir})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("encode: %v", err), http.StatusInternalServerError)
	}
}
