package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	Running       bool `json:"running"`
	LaunchAtLogin bool `json:"launchAtLogin"`
}

type Store struct {
	path string
}

func NewStore(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(baseDir, "state.json")}, nil
}

func (s *Store) Load() (State, error) {
	var st State
	bytes, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	if err := json.Unmarshal(bytes, &st); err != nil {
		return st, err
	}
	return st, nil
}

func (s *Store) Save(state State) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0o644)
}
