package session

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"wails-datastar-example/internal/token"
)

const defaultAPIBase = "http://127.0.0.1:8080"

// Session is what the page reads once, before bootstrap.
// The page must keep Token out of Datastar signals. Signals are sent
// to the server on every request; the token travels only as a header.
type Session struct {
	APIBase string `json:"apiBase"`
	Token   string `json:"token"`
}

// Load reads the API base URL and the bearer token.
func Load(store token.Store) (Session, error) {
	base, err := APIBase()
	if err != nil {
		return Session{}, err
	}
	tok, err := token.Resolve(store, token.DevToken())
	if err != nil {
		return Session{}, err
	}
	return Session{APIBase: base, Token: tok}, nil
}

// Save writes a pasted token into the store.
func Save(store token.Store, raw string) error {
	tok, err := CleanToken(raw)
	if err != nil {
		return err
	}
	return store.Set(tok)
}

// CleanToken trims a pasted bearer token and rejects values that
// cannot sit in an Authorization header.
func CleanToken(raw string) (string, error) {
	tok := strings.TrimSpace(raw)
	if tok == "" {
		return "", errors.New("token is empty")
	}
	if strings.ContainsAny(tok, "\r\n") {
		return "", errors.New("token is invalid")
	}
	return tok, nil
}

// APIBase is DESKTOP_API_BASE, or http://127.0.0.1:8080 when unset.
func APIBase() (string, error) {
	raw := strings.TrimSpace(os.Getenv("DESKTOP_API_BASE"))
	if raw == "" {
		raw = defaultAPIBase
	}
	raw = strings.TrimRight(raw, "/")
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("DESKTOP_API_BASE must be an absolute http(s) URL")
	}
	return u.Scheme + "://" + u.Host + u.Path, nil
}
