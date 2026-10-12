package config

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

// Use an isolated service; never read or replace the user's real login session.
func TestSystemStoreRoundTrip(t *testing.T) {
	store := SystemStore{serviceName: "dlut-pan-cli-test-" + strconv.FormatInt(time.Now().UnixNano(), 10)}
	session := Session{Token: "fixture-session-token", BaseURL: BaseURL}
	if err := store.Save(session); err != nil {
		t.Skip("OS keyring unavailable in this environment:", err)
	}
	t.Cleanup(func() {
		if err := store.Clear(); err != nil {
			t.Error(err)
		}
	})
	got, err := store.Load()
	if err != nil || got != session {
		t.Fatal("saved session did not round trip")
	}
	session.Token = "replacement-fixture-token"
	if err := store.Save(session); err != nil {
		t.Fatal(err)
	}
	got, err = store.Load()
	if err != nil || got != session {
		t.Fatal("session replacement did not round trip")
	}
	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(); !errors.Is(err, ErrNoSession) {
		t.Fatal("logout did not remove saved session")
	}
	if err := store.Clear(); err != nil {
		t.Fatal("logout must be idempotent:", err)
	}
}

func TestRefuseEmptySession(t *testing.T) {
	for _, session := range []Session{{}, {Token: "token"}, {BaseURL: BaseURL}} {
		if err := (SystemStore{}).Save(session); err == nil {
			t.Fatal("accepted incomplete session")
		}
	}
}
