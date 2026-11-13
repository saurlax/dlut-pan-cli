package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileStoreRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	store := FileStore{}
	session := Session{Token: "fixture-session-token", BaseURL: BaseURL, Username: "student", Password: "fixture-password"}
	if err := store.Save(session); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".dlut-pan", "session.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		dir, err := os.Stat(filepath.Dir(path))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm()&0o077 != 0 || dir.Mode().Perm()&0o077 != 0 {
			t.Fatal("session or session directory is accessible to other users")
		}
	}
	got, err := store.Load()
	if err != nil || got != session {
		t.Fatal("saved session did not round trip")
	}
	session = Session{Token: "replacement-fixture-token", BaseURL: BaseURL}
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
		if err := (FileStore{}).Save(session); err == nil {
			t.Fatal("accepted incomplete session")
		}
	}
}
