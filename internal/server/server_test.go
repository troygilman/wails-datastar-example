package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBootstrapRequiresToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/desktop/bootstrap", nil)
	req.Header.Set("Origin", "wails://wails")
	rec := httptest.NewRecorder()
	New("dev").ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "wails://wails" {
		t.Fatalf("cors %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestBootstrapPatchesShell(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/desktop/bootstrap", nil)
	req.Header.Set("Origin", "http://wails.localhost")
	req.Header.Set("Authorization", "Bearer desktop-test-token")
	rec := httptest.NewRecorder()
	New("desktop-test-token").ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("content-type %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("X-Accel-Buffering") != "no" {
		t.Fatal("missing X-Accel-Buffering")
	}
	body := rec.Body.String()
	for _, want := range []string{
		"event: datastar-patch-elements",
		`id="app"`,
		`id="panel-main"`,
		`id="panel-side"`,
		`id="ping"`,
		"event: datastar-patch-signals",
		`"ready":true`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q\n%s", want, body)
		}
	}
	if strings.Contains(body, "desktop-test-token") {
		t.Fatal("response includes the token")
	}
}

func TestRejectsOtherOrigins(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/desktop/bootstrap", nil)
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Authorization", "Bearer dev")
	rec := httptest.NewRecorder()
	New("dev").ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("reflected a foreign origin")
	}
}

func TestPreflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/desktop/bootstrap", nil)
	req.Header.Set("Origin", "wails://wails")
	rec := httptest.NewRecorder()
	New("dev").ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "Authorization") || !strings.Contains(got, "Datastar-Request") {
		t.Fatalf("headers %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "GET") || !strings.Contains(got, "POST") {
		t.Fatalf("methods %q", got)
	}
}

func TestPingRoundTrip(t *testing.T) {
	s := New("dev")
	for i, want := range []string{"Pings: 1", "Pings: 2"} {
		req := httptest.NewRequest(http.MethodPost, "/desktop/ping", nil)
		req.Header.Set("Authorization", "Bearer dev")
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("ping %d status %d", i, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("ping %d body %s", i, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), `id="panel-main-note"`) {
			t.Fatal("ping replaced an unknown id")
		}
	}
}

func TestNoOriginStillServes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/desktop/bootstrap", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	New("secret").ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Body)
	if rec.Code != http.StatusOK || !strings.Contains(string(body), `id="app"`) {
		t.Fatalf("status %d body %s", rec.Code, body)
	}
}
