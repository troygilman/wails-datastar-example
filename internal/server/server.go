// Package server is the remote Datastar process. The desktop binary does not import it.
package server

import (
	"crypto/subtle"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/starfederation/datastar-go/datastar"
)

const (
	pathBootstrap = "/desktop/bootstrap"
	pathPing      = "/desktop/ping"
)

// Origins the desktop webview actually sends. Wails v2 uses a custom scheme
// on Linux and macOS, and http://wails.localhost on Windows.
var allowedOrigins = map[string]struct{}{
	"http://wails.localhost":  {},
	"https://wails.localhost": {},
	"wails://wails":           {},
}

// Server serves the shell's bootstrap frame and one follow-up action.
type Server struct {
	token string

	mu    sync.Mutex
	pings int
}

func New(token string) *Server {
	return &Server{token: token}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		log.Printf("%s %s origin %s", r.Method, r.URL.Path, origin)
	}
	if origin != "" {
		if _, ok := allowedOrigins[origin]; !ok {
			http.Error(w, "origin not allowed", http.StatusForbidden)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Datastar-Request, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == pathBootstrap:
		s.bootstrap(w, r)
	case r.Method == http.MethodPost && r.URL.Path == pathPing:
		s.ping(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) bootstrap(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	w.Header().Set("X-Accel-Buffering", "no")
	sse := datastar.NewSSE(w, r)
	if err := sse.PatchElements(frame("Hello from the server.", "Connected.")); err != nil {
		log.Printf("bootstrap patch: %v", err)
		return
	}
	if err := sse.MarshalAndPatchSignals(map[string]any{"ready": true}); err != nil {
		log.Printf("bootstrap signals: %v", err)
	}
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	s.mu.Lock()
	s.pings++
	n := s.pings
	s.mu.Unlock()

	w.Header().Set("X-Accel-Buffering", "no")
	sse := datastar.NewSSE(w, r)
	note := fmt.Sprintf(`<p id="panel-main-note">Pings: %d</p>`, n)
	if err := sse.PatchElements(note); err != nil {
		log.Printf("ping patch: %v", err)
	}
}

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) bool {
	got := bearerFrom(r.Header.Get("Authorization"))
	if subtle.ConstantTimeCompare([]byte(got), []byte(s.token)) == 1 {
		return true
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
	return false
}

func bearerFrom(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func frame(note, side string) string {
	return `<div id="app">` +
		`<section id="panel-main">` +
		`<p id="panel-main-note">` + html.EscapeString(note) + `</p>` +
		`<button id="ping" type="button" data-on:click="@post($api + '/desktop/ping', {headers: desktop.headers()})">Ping</button>` +
		`</section>` +
		`<aside id="panel-side">` +
		`<p id="panel-side-note">` + html.EscapeString(side) + `</p>` +
		`</aside>` +
		`</div>`
}
