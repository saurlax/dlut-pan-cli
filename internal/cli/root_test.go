package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saurlax/dlut-pan-cli/internal/config"
)

type memoryStore struct {
	session                    config.Session
	loadErr, saveErr, clearErr error
}

func (s *memoryStore) Load() (config.Session, error) {
	if s.loadErr != nil {
		return config.Session{}, s.loadErr
	}
	if s.session.Token == "" {
		return config.Session{}, config.ErrNoSession
	}
	return s.session, nil
}
func (s *memoryStore) Save(session config.Session) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.session = session
	return nil
}
func (s *memoryStore) Clear() error {
	if s.clearErr != nil {
		return s.clearErr
	}
	s.session = config.Session{}
	return nil
}

func run(t *testing.T, input string, values ...string) (int, string, string) {
	return runStore(t, &memoryStore{}, input, values...)
}

func runStore(t *testing.T, store config.SessionStore, input string, values ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := execute(context.Background(), values, strings.NewReader(input), &out, &errOut, "test", store)
	return code, out.String(), errOut.String()
}

func cleanEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"TOKEN", "USERNAME", "PASSWORD", "ROOT_ID", "BASE_URL", "AUTH_URL"} {
		t.Setenv("DLUT_PAN_"+key, "")
	}
}

func TestHelpAndOfflineCommands(t *testing.T) {
	cleanEnv(t)
	for _, command := range [][]string{{}, {"auth"}, {"file"}, {"search"}, {"group"}, {"space"}, {"file", "download", "--help"}, {"auth", "login", "--help"}, {"completion", "powershell"}, {"--version"}} {
		code, out, err := run(t, "", command...)
		if code != 0 || out == "" || err != "" {
			t.Fatalf("%v: code=%d out=%s err=%s", command, code, out, err)
		}
	}
	code, out, _ := run(t, "", "space", "naosi", "--json")
	if code != 0 || !strings.Contains(out, "ak7an4chdvtq") {
		t.Fatal(out)
	}
}

func TestUsageErrorsBeforeNetwork(t *testing.T) {
	cleanEnv(t)
	for _, command := range [][]string{
		{"search", "files"}, {"search", "files", "test", "--limit", "0"}, {"search", "files", "test", "--offset", "-1"},
		{"file", "list", "relative"}, {"file", "download", "id"}, {"file", "download", "../id", "-o", "test"},
		{"group", "join", "id", "--user-id", "u"}, {"--unknown"}, {"unknown"}, {"file", "--timeout", "0s"},
		{"file", "list", "--base-url", "https://example.com?token=bad"},
	} {
		code, _, err := run(t, "", command...)
		if code != 2 {
			t.Errorf("%v returned %d: %s", command, code, err)
		}
	}
}

func TestDryRunProtocolAndRedaction(t *testing.T) {
	cleanEnv(t)
	t.Setenv("DLUT_PAN_TOKEN", "never-show-token")
	t.Setenv("DLUT_PAN_PASSWORD", "never-show-password")
	for _, command := range [][]string{
		{"file", "list", "/课程"}, {"search", "files", "test"}, {"search", "full-text", "test"}, {"search", "groups", "NAOSI"},
		{"group", "join", "gid", "--user-id", "uid"}, {"file", "download", "fid", "--output", "x.pdf"}, {"auth", "login"},
	} {
		command = append(command, "--dry-run")
		code, out, err := run(t, "", command...)
		if code != 0 || !json.Valid([]byte(out)) || strings.Contains(out, "never-show") {
			t.Fatalf("%v: code=%d out=%s err=%s", command, code, out, err)
		}
		if command[0] == "search" && command[1] == "groups" && strings.Contains(out, "root_id") {
			t.Fatal("group search must not send root ID")
		}
	}
}

func TestAuthenticatedCommands(t *testing.T) {
	cleanEnv(t)
	for _, spec := range []struct {
		command      []string
		method, path string
	}{
		{[]string{"file", "list", "/课程"}, "POST", "/v1/fileops/list_folder"},
		{[]string{"search", "files", "关键词", "--path", "/课程", "--offset", "10"}, "GET", "/v1/search/files"},
		{[]string{"search", "full-text", "关键词"}, "GET", "/v1/search/full_text_files"},
		{[]string{"search", "groups", "NAOSI"}, "GET", "/v1/search/groups"},
		{[]string{"group", "join", "gid", "--user-id", "uid", "--yes", "--remarks", "学习"}, "POST", "/v1/users/uid/groups"},
	} {
		t.Run(spec.path, func(t *testing.T) {
			var hit bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hit = true
				if r.Method != spec.method || r.URL.Path != spec.path || r.URL.Query().Get("token") != "stdin-token" {
					t.Error("request mismatch")
				}
				if strings.Contains(spec.path, "users") {
					b, _ := io.ReadAll(r.Body)
					if string(b) != `{"remarks":"学习"}` {
						t.Errorf("join body: %s", b)
					}
				}
				fmt.Fprint(w, `{"total":1,"entries":[{"id":"file-id","name":"课程","size":{"bytes":123},"is_dir":false}],"token":"should-hide"}`)
			}))
			defer server.Close()
			values := append(append([]string{}, spec.command...), "--json")
			store := &memoryStore{session: config.Session{Token: "stdin-token", BaseURL: server.URL + "/v1"}}
			code, out, err := runStore(t, store, "", values...)
			if code != 0 || !hit || strings.Contains(out, "should-hide") {
				t.Fatalf("code=%d out=%s err=%s", code, out, err)
			}
		})
	}
}

func TestLoginPersistsAcrossInvocationsAndLogout(t *testing.T) {
	cleanEnv(t)
	t.Setenv("DLUT_PAN_TOKEN", "stale")
	var loginHits, listHits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sign-in" {
			loginHits++
			if r.Method != "POST" || r.URL.Query().Get("password") != "secret +&" {
				t.Error("login mismatch")
			}
			fmt.Fprint(w, `{"token":"fresh"}`)
			return
		}
		listHits++
		if r.URL.Query().Get("token") != "fresh" {
			t.Error("did not use fresh token")
		}
		fmt.Fprint(w, `{"entries":[]}`)
	}))
	defer server.Close()
	store := &memoryStore{}
	code, out, err := runStore(t, store, "secret +&\n", "auth", "login", "--username", "student", "--password-stdin", "--auth-url", server.URL+"/sign-in", "--base-url", server.URL)
	if code != 0 || loginHits != 1 || store.session.Token != "fresh" || strings.Contains(out, "fresh") || strings.Contains(out, "secret") {
		t.Fatalf("code=%d out=%s err=%s", code, out, err)
	}
	// A new command tree represents a later CLI process. No flags or input credentials.
	for i := 0; i < 2; i++ {
		code, out, err = runStore(t, store, "", "file", "list")
		if code != 0 || !strings.Contains(out, "No entries") {
			t.Fatalf("reuse: code=%d out=%s err=%s", code, out, err)
		}
	}
	if loginHits != 1 || listHits != 2 {
		t.Fatal("data command reauthenticated instead of reusing session")
	}
	code, out, err = runStore(t, store, "", "auth", "status")
	if code != 0 || !strings.Contains(out, `"session_saved": true`) || strings.Contains(out, "fresh") {
		t.Fatal(out, err)
	}
	code, out, err = runStore(t, store, "", "auth", "logout")
	if code != 0 || store.session.Token != "" {
		t.Fatal(out, err)
	}
	code, _, err = runStore(t, store, "", "file", "list")
	if code != 1 || !strings.Contains(err, "auth login") || listHits != 2 {
		t.Fatal("logout did not prevent remote operation", err)
	}
}

func TestTableAndUnknownSchema(t *testing.T) {
	cleanEnv(t)
	for _, response := range []string{`{"entries":[{"id":"id","name":"test\u001b[31m","size":{"display_value":"1 KB"},"is_dir":true}]}`, `{"other_shape":{"key":123}}`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, response) }))
		store := &memoryStore{session: config.Session{Token: "test", BaseURL: server.URL}}
		code, out, err := runStore(t, store, "", "file", "list")
		server.Close()
		if code != 0 || strings.Contains(out, "\x1b") {
			t.Fatalf("code=%d out=%s err=%s", code, out, err)
		}
		if strings.Contains(response, "other_shape") && !json.Valid([]byte(out)) {
			t.Fatal("unknown schema was not preserved")
		}
	}
}

func TestRequiresSavedLoginEvenWithEnvironmentCredentials(t *testing.T) {
	cleanEnv(t)
	t.Setenv("DLUT_PAN_TOKEN", "ignored-token")
	t.Setenv("DLUT_PAN_USERNAME", "student")
	t.Setenv("DLUT_PAN_PASSWORD", "secret")
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++; w.WriteHeader(500) }))
	defer server.Close()
	for _, command := range [][]string{{"file", "list"}, {"search", "files", "test"}, {"search", "groups", "NAOSI"}, {"file", "download", "fid", "-o", "test.pdf"}, {"group", "join", "gid", "--user-id", "uid", "--yes"}} {
		command = append(command, "--base-url", server.URL)
		code, _, err := run(t, "", command...)
		if code != 1 || !strings.Contains(err, "auth login") {
			t.Fatalf("%v code=%d err=%s", command, code, err)
		}
	}
	if hits != 0 {
		t.Fatal("unauthenticated command contacted server")
	}
	for _, command := range [][]string{{"file", "list", "--username", "student"}, {"search", "files", "test", "--token", "secret"}, {"file", "list", "--token-stdin"}} {
		code, _, _ := run(t, "", command...)
		if code != 2 {
			t.Fatal("data command still accepts login credentials")
		}
	}
}

func TestSessionImportAndFailedLoginPreservesSession(t *testing.T) {
	cleanEnv(t)
	store := &memoryStore{}
	code, out, err := runStore(t, store, "browser-token\n", "auth", "login", "--token-stdin")
	if code != 0 || store.session.Token != "browser-token" || strings.Contains(out, "browser-token") {
		t.Fatal(out, err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(401) }))
	defer server.Close()
	code, _, _ = runStore(t, store, "wrong-password", "auth", "login", "--username", "student", "--password-stdin", "--auth-url", server.URL)
	if code != 1 || store.session.Token != "browser-token" {
		t.Fatal("failed login replaced previous session")
	}
	code, _, err = runStore(t, store, "", "file", "list", "--base-url", server.URL)
	if code != 1 || !strings.Contains(err, "another API endpoint") {
		t.Fatal("endpoint change reused wrong token", err)
	}
}

func TestExpiredSessionAndPermissionFailure(t *testing.T) {
	cleanEnv(t)
	for _, status := range []int{401, 403} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status) }))
		store := &memoryStore{session: config.Session{Token: "private-token", BaseURL: server.URL}}
		code, _, err := runStore(t, store, "", "file", "list")
		server.Close()
		if code != 1 || strings.Contains(err, "private-token") {
			t.Fatal(err)
		}
		if (status == 401) != strings.Contains(err, "auth login") {
			t.Fatal("incorrect relogin advice", err)
		}
		if store.session.Token == "" {
			t.Fatal("request error deleted session")
		}
	}
}

func TestStoreErrorsAndDryRun(t *testing.T) {
	cleanEnv(t)
	store := &memoryStore{loadErr: errors.New("credential store locked")}
	code, _, err := runStore(t, store, "", "file", "list")
	if code != 1 || !strings.Contains(err, "locked") {
		t.Fatal(err)
	}
	code, _, err = runStore(t, store, "", "file", "list", "--dry-run")
	if code != 0 {
		t.Fatal(err)
	}
	store = &memoryStore{saveErr: errors.New("cannot save session")}
	code, out, err := runStore(t, store, "imported-token", "auth", "login", "--token-stdin")
	if code != 1 || strings.Contains(out, "session_saved") || store.session.Token != "" {
		t.Fatal("save failure reported success", out, err)
	}
	code, _, err = runStore(t, store, "ignored", "auth", "login", "--token-stdin", "--dry-run")
	if code != 0 || store.session.Token != "" {
		t.Fatal("dry-run saved session", err)
	}
	store = &memoryStore{session: config.Session{Token: "existing", BaseURL: config.BaseURL}}
	code, _, err = runStore(t, store, "", "auth", "logout", "--dry-run")
	if code != 0 || store.session.Token != "existing" {
		t.Fatal("dry-run cleared session", err)
	}
	store.clearErr = errors.New("cannot remove session")
	code, _, err = runStore(t, store, "", "auth", "logout")
	if code != 1 || store.session.Token != "existing" {
		t.Fatal("clear failure lost session", err)
	}
	code, out, err = run(t, "", "auth", "status")
	if code != 0 || !strings.Contains(out, `"session_saved": false`) {
		t.Fatal(out, err)
	}
}
