package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func run(t *testing.T, input string, values ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := Execute(context.Background(), values, strings.NewReader(input), &out, &errOut, "test")
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
			values := append(append([]string{}, spec.command...), "--base-url", server.URL+"/v1", "--token-stdin", "--json")
			code, out, err := run(t, "stdin-token\n", values...)
			if code != 0 || !hit || strings.Contains(out, "should-hide") {
				t.Fatalf("code=%d out=%s err=%s", code, out, err)
			}
		})
	}
}

func TestLoginThenOperationInOneProcess(t *testing.T) {
	cleanEnv(t)
	t.Setenv("DLUT_PAN_TOKEN", "stale")
	var loginHit, listHit bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sign-in" {
			loginHit = true
			if r.Method != "POST" || r.URL.Query().Get("password") != "secret +&" {
				t.Error("login mismatch")
			}
			fmt.Fprint(w, `{"token":"fresh"}`)
			return
		}
		listHit = true
		if r.URL.Query().Get("token") != "fresh" {
			t.Error("did not use fresh token")
		}
		fmt.Fprint(w, `{"entries":[]}`)
	}))
	defer server.Close()
	code, out, err := run(t, "secret +&\n", "file", "list", "--username", "student", "--password-stdin", "--auth-url", server.URL+"/sign-in", "--base-url", server.URL)
	if code != 0 || !loginHit || !listHit || !strings.Contains(out, "No entries") {
		t.Fatalf("code=%d out=%s err=%s", code, out, err)
	}
}

func TestTableAndUnknownSchema(t *testing.T) {
	cleanEnv(t)
	for _, response := range []string{`{"entries":[{"id":"id","name":"test\u001b[31m","size":{"display_value":"1 KB"},"is_dir":true}]}`, `{"other_shape":{"key":123}}`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, response) }))
		code, out, err := run(t, "", "file", "list", "--token", "test", "--base-url", server.URL)
		server.Close()
		if code != 0 || strings.Contains(out, "\x1b") {
			t.Fatalf("code=%d out=%s err=%s", code, out, err)
		}
		if strings.Contains(response, "other_shape") && !json.Valid([]byte(out)) {
			t.Fatal("unknown schema was not preserved")
		}
	}
}
