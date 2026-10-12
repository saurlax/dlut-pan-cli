package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProtocol(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch r.URL.Path {
		case "/v1/auth/sign_in":
			if r.Method != "POST" || r.ContentLength != 0 || q.Get("username") != "student" || q.Get("password") != "p&+?#中文" || q.Get("ldap_name") != "device" {
				t.Errorf("incorrect login protocol: method=%s length=%d", r.Method, r.ContentLength)
			}
			fmt.Fprint(w, `{"token":"test-token","user_id":"private-id"}`)
		case "/v1/fileops/list_folder":
			if r.Method != "POST" || r.ContentLength != 0 || q.Get("token") != "test-token" || q.Get("path") != "/课程/+ &" || q.Get("root_id") != "space" {
				t.Error("incorrect list protocol")
			}
			fmt.Fprint(w, `{"entries":[{"id":"f1","name":"课程.pdf"}]}`)
		default:
			t.Error("unexpected path")
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	c := &Client{BaseURL: server.URL + "/v1", AuthURL: server.URL + "/v1/auth/sign_in", Timeout: time.Second}
	result, err := c.Login(context.Background(), "student", "p&+?#中文", "device")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(fmt.Sprint(result), "test-token") || strings.Contains(fmt.Sprint(result), "private-id") {
		t.Fatal("login metadata leaked")
	}
	_, err = c.Call(context.Background(), "POST", "/fileops/list_folder", url.Values{"root_id": {"space"}, "path": {"/课程/+ &"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestErrorRedactionAndAuthRedirect(t *testing.T) {
	var targetHit bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetHit = true }))
	defer target.Close()
	for _, status := range []int{401, 403, 500, 302} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", target.URL)
				w.WriteHeader(status)
				fmt.Fprint(w, "password=hunter2&token=private-token")
			}))
			c := &Client{BaseURL: server.URL, Token: "private-token", Timeout: time.Second}
			_, err := c.Call(context.Background(), "GET", "/test", nil, nil)
			server.Close()
			if err == nil || strings.Contains(err.Error(), "hunter2") || strings.Contains(err.Error(), "private-token") {
				t.Fatalf("unsafe or missing error: %v", err)
			}
		})
	}
	if targetHit {
		t.Fatal("non-download request followed redirect")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	endpoint := server.URL
	server.Close()
	c := &Client{BaseURL: endpoint, Token: "private-token", Timeout: time.Second}
	_, err := c.Call(context.Background(), "GET", "/test", nil, nil)
	if err == nil || strings.Contains(err.Error(), "private-token") || strings.Contains(err.Error(), endpoint) {
		t.Fatalf("URL leaked in transport error: %v", err)
	}
}

func TestDownloadAtomicityAndRedirect(t *testing.T) {
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != "" {
			t.Error("token forwarded across origin")
		}
		if r.URL.Query().Get("signature") != "keep" {
			t.Error("storage signature lost")
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprint(w, "file contents")
	}))
	defer storage.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, storage.URL+"?token=private-token&signature=keep", 302)
	}))
	defer server.Close()
	c := &Client{BaseURL: server.URL, Token: "private-token", Timeout: time.Second}
	dir := t.TempDir()
	destination := filepath.Join(dir, "lecture.pdf")
	n, err := c.Download(context.Background(), "root", "file", destination, false)
	if err != nil || n != 13 {
		t.Fatalf("download: n=%d err=%v", n, err)
	}
	if _, err = c.Download(context.Background(), "root", "file", destination, false); err == nil {
		t.Fatal("overwrote without consent")
	}
	if err = os.WriteFile(destination, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		fmt.Fprint(w, "short")
	}))
	defer broken.Close()
	c.BaseURL = broken.URL
	if _, err = c.Download(context.Background(), "root", "file", destination, true); err == nil {
		t.Fatal("accepted truncated body")
	}
	b, _ := os.ReadFile(destination)
	if string(b) != "original" {
		t.Fatal("interruption damaged original")
	}
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Fatal("partial file left behind")
	}
	c.BaseURL = server.URL
	if _, err = c.Download(context.Background(), "root", "file", destination, true); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(destination)
	if string(b) != "file contents" {
		t.Fatal("overwrite did not publish complete file")
	}
}

func TestLoginMissingTokenAndInvalidJSON(t *testing.T) {
	for _, body := range []string{`{"error":"denied"}`, `[]`, `<html>login</html>`, `{"token":"ok"} garbage`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, body) }))
		c := &Client{AuthURL: server.URL, Timeout: time.Second}
		if _, err := c.Login(context.Background(), "student", "secret", "device"); err == nil {
			t.Errorf("accepted invalid login response %q", body)
		}
		server.Close()
	}
}

func TestRedact(t *testing.T) {
	v := map[string]any{"name": "file", "nested": []any{map[string]any{"access_token": "private", "password": "hidden"}}}
	got := fmt.Sprint(Redact(v))
	if strings.Contains(got, "private") || strings.Contains(got, "hidden") || !strings.Contains(got, "file") {
		t.Fatalf("incorrect redaction: %s", got)
	}
}
