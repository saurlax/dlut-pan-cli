package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoSession = errors.New("no saved login session")

// Session stores the access token, its API endpoint and optional account credentials.
type Session struct {
	Token    string `json:"token"`
	BaseURL  string `json:"base_url"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SessionStore interface {
	Load() (Session, error)
	Save(Session) error
	Clear() error
}

// FileStore saves a single session in the current user's home directory.
type FileStore struct{}

func sessionPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || !filepath.IsAbs(home) {
		return "", errors.New("cannot determine the home directory for the login session")
	}
	return filepath.Join(home, ".dlut-pan", "session.json"), nil
}

func (FileStore) Load() (Session, error) {
	path, err := sessionPath()
	if err != nil {
		return Session{}, err
	}
	value, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Session{}, ErrNoSession
	}
	if err != nil {
		return Session{}, err
	}
	var session Session
	if json.Unmarshal(value, &session) != nil || !session.valid() {
		return Session{}, errors.New("saved login session is invalid; run 'dlut-pan auth login' again")
	}
	return session, nil
}

func (FileStore) Save(session Session) error {
	if !session.valid() {
		return errors.New("cannot save an empty login session")
	}
	path, err := sessionPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".session-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	// Stage the new credentials before replacing the previous session.
	writeErr := json.NewEncoder(file).Encode(session)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func (FileStore) Clear() error {
	path, err := sessionPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s Session) valid() bool {
	return strings.TrimSpace(s.Token) != "" && strings.TrimSpace(s.BaseURL) != ""
}
