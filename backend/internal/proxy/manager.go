package proxy

import (
	"errors"
	"os/exec"
	"sync"
)

type Manager struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	running bool
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	// TODO: replace with real CLIProxyAPIPlus command and args on production setup.
	m.cmd = exec.Command("cmd", "/C", "ping", "127.0.0.1", "-n", "99999")
	if err := m.cmd.Start(); err != nil {
		return err
	}

	m.running = true
	go func() {
		_ = m.cmd.Wait()
		m.mu.Lock()
		m.running = false
		m.mu.Unlock()
	}()

	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cmd == nil || m.cmd.Process == nil {
		return nil
	}

	if err := m.cmd.Process.Kill(); err != nil {
		return err
	}
	m.running = false
	return nil
}

func (m *Manager) Running() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

var ErrNotImplemented = errors.New("proxy adapter is not fully wired to CLIProxyAPIPlus yet")
