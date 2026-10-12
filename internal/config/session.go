package config

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

var ErrNoSession = errors.New("no saved login session")

// Session stores only the access token and the endpoint it belongs to, never a password.
type Session struct {
	Token   string `json:"token"`
	BaseURL string `json:"base_url"`
}

type SessionStore interface {
	Load() (Session, error)
	Save(Session) error
	Clear() error
}

type SystemStore struct{ serviceName string }

const service = "dlut-pan-cli"
const account = "session"

func (s SystemStore) name() string {
	if s.serviceName != "" {
		return s.serviceName
	}
	return service
}

func (s SystemStore) Load() (Session, error) {
	value, err := keyring.Get(s.name(), account)
	if errors.Is(err, keyring.ErrNotFound) {
		return Session{}, ErrNoSession
	}
	if err != nil {
		return Session{}, errors.New("cannot read login session from system credential store; unlock or enable the OS keyring")
	}
	var session Session
	if json.Unmarshal([]byte(value), &session) != nil || strings.TrimSpace(session.Token) == "" || strings.TrimSpace(session.BaseURL) == "" {
		return Session{}, errors.New("saved login session is invalid; run 'dlut-pan auth login' again")
	}
	return session, nil
}

func (s SystemStore) Save(session Session) error {
	if strings.TrimSpace(session.Token) == "" || strings.TrimSpace(session.BaseURL) == "" {
		return errors.New("cannot save an empty login session")
	}
	value, err := json.Marshal(session)
	if err != nil {
		return errors.New("cannot encode login session")
	}
	if err = keyring.Set(s.name(), account, string(value)); err != nil {
		return errors.New("cannot save login session to system credential store; unlock or enable the OS keyring")
	}
	return nil
}

func (s SystemStore) Clear() error {
	err := keyring.Delete(s.name(), account)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return errors.New("cannot remove login session from system credential store; unlock or enable the OS keyring")
	}
	return nil
}
