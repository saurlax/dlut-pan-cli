// Package api implements the documented DLUT Pan protocol without logging secrets.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	BaseURL, AuthURL, Token string
	Timeout                 time.Duration
}

// HTTPError exposes a status without retaining request URLs or response bodies.
type HTTPError struct {
	StatusCode int
	Hint       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("API returned HTTP %d: %s", e.StatusCode, e.Hint)
}

func ValidateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("API URL must be an absolute HTTP(S) URL without credentials, query or fragment")
	}
	return nil
}

func (c *Client) endpoint(path string, query url.Values) string {
	q := url.Values{}
	for k, v := range query {
		q[k] = append([]string(nil), v...)
	}
	if c.Token != "" {
		q.Set("token", c.Token)
	}
	return strings.TrimRight(c.BaseURL, "/") + path + "?" + q.Encode()
}

func (c *Client) send(ctx context.Context, method, endpoint string, body io.Reader, download bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, errors.New("cannot construct API request")
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "dlut-pan-cli")
	client := &http.Client{Timeout: c.Timeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if !download {
			return http.ErrUseLastResponse
		}
		if len(via) >= 10 {
			return errors.New("too many download redirects")
		}
		if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
			return errors.New("unsupported redirect scheme")
		}
		if via[len(via)-1].URL.Scheme == "https" && req.URL.Scheme == "http" {
			return errors.New("download redirect downgrades HTTPS")
		}
		if req.URL.Host != via[0].URL.Host {
			q := req.URL.Query()
			for _, k := range []string{"token", "password", "username"} {
				q.Del(k)
			}
			req.URL.RawQuery = q.Encode()
		}
		return nil
	}}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// net/url errors contain the complete request URL, including credentials.
		return nil, errors.New("API request failed (network, timeout or redirect); check connectivity and endpoint")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		hint := "check the API endpoint and request parameters"
		switch resp.StatusCode {
		case 401:
			hint = "token expired or credentials rejected"
		case 403:
			hint = "check root ID, group membership and permissions"
		case 404:
			hint = "resource or endpoint not found"
		case 429:
			hint = "rate limited; retry later"
		}
		return nil, &HTTPError{StatusCode: resp.StatusCode, Hint: hint}
	}
	return resp, nil
}

func decode(resp *http.Response) (any, error) {
	defer resp.Body.Close()
	const max = 16 << 20
	b, err := io.ReadAll(io.LimitReader(resp.Body, max+1))
	if err != nil {
		return nil, errors.New("cannot read API response")
	}
	if len(b) > max {
		return nil, errors.New("API JSON exceeds 16 MiB; narrow the search")
	}
	var value any
	d := json.NewDecoder(strings.NewReader(string(b)))
	d.UseNumber()
	if err = d.Decode(&value); err != nil {
		return nil, errors.New("API did not return valid JSON")
	}
	if d.Decode(new(any)) != io.EOF {
		return nil, errors.New("API returned trailing data after JSON")
	}
	return value, nil
}

func (c *Client) Call(ctx context.Context, method, path string, query url.Values, body io.Reader) (any, error) {
	resp, err := c.send(ctx, method, c.endpoint(path, query), body, false)
	if err != nil {
		return nil, err
	}
	return decode(resp)
}

func (c *Client) Login(ctx context.Context, username, password, device string) (any, error) {
	q := url.Values{"username": {username}, "password": {password}, "link_device": {device}, "link_name": {device}, "ldap_name": {device}}
	resp, err := c.send(ctx, "POST", c.AuthURL+"?"+q.Encode(), nil, false)
	if err != nil {
		return nil, err
	}
	value, err := decode(resp)
	if err != nil {
		return nil, err
	}
	m, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("login response is not an object")
	}
	token, _ := m["token"].(string)
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("login response contains no token; authentication failed or API changed")
	}
	c.Token = token
	// Login metadata may contain identity details; only expose success.
	return map[string]any{"authenticated": true, "credentials_saved": false}, nil
}

// Download stages bytes next to the destination, then publishes a complete file.
// Without overwrite, a hard link atomically refuses an existing destination.
func (c *Client) Download(ctx context.Context, root, id, destination string, overwrite bool) (int64, error) {
	if !overwrite {
		if _, err := os.Lstat(destination); err == nil {
			return 0, errors.New("destination exists; choose another --output or pass --overwrite")
		} else if !os.IsNotExist(err) {
			return 0, err
		}
	}
	resp, err := c.send(ctx, "GET", c.endpoint("/roots/"+url.PathEscape(root)+"/files/"+url.PathEscape(id), nil), nil, true)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	media := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(media, "text/html") {
		return 0, errors.New("download returned HTML, possibly a login page; no file saved")
	}
	f, err := os.CreateTemp(filepath.Dir(destination), ".dlut-pan-download-*")
	if err != nil {
		return 0, err
	}
	defer os.Remove(f.Name())
	n, copyErr := io.Copy(f, resp.Body)
	syncErr := f.Sync()
	closeErr := f.Close()
	if copyErr != nil {
		return 0, errors.New("download interrupted; partial file discarded")
	}
	if syncErr != nil {
		return 0, syncErr
	}
	if closeErr != nil {
		return 0, closeErr
	}
	if resp.ContentLength >= 0 && n != resp.ContentLength {
		return 0, errors.New("download length mismatch; partial file discarded")
	}
	if overwrite {
		err = os.Rename(f.Name(), destination)
	} else {
		err = os.Link(f.Name(), destination)
	}
	if err != nil {
		return 0, fmt.Errorf("cannot publish downloaded file: %w", err)
	}
	return n, nil
}

// Redact suppresses secret fields even when an API changes its response shape.
func Redact(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, item := range v {
			lower := strings.ToLower(k)
			if strings.Contains(lower, "token") || strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "cookie") {
				out[k] = "[REDACTED]"
			} else {
				out[k] = Redact(item)
			}
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = Redact(item)
		}
		return out
	default:
		return value
	}
}
