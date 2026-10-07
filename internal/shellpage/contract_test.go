package shellpage

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func dist(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "frontend", "dist")
}

func TestShellMarkup(t *testing.T) {
	html, err := os.ReadFile(filepath.Join(dist(t), "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(html)
	for _, want := range []string{
		`id="app"`,
		`id="panel-main"`,
		`id="panel-side"`,
		`id="shell-status"`,
		`id="shell-signin"`,
		`id="shell-retry"`,
		`id="shell-token"`,
		`ready: false`,
		`phase: 'connecting'`,
		`src="datastar.js"`,
		`src="shell.js"`,
		`href="shell.css"`,
		`openWhenHidden: true`,
		`retry: "auto"`,
		`retryMaxCount: 5`,
		`desktop.headers()`,
		`/desktop/bootstrap`,
		`evt.detail.argsRaw.status === '401'`,
		`retries-failed`,
	} {
		if !strings.Contains(page, want) && !containsInShell(t, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(page, "cdn.") || strings.Contains(page, "jsdelivr") {
		t.Fatal("shell must not load a CDN")
	}
	if strings.Contains(page, "$token") || strings.Contains(page, "token:") {
		t.Fatal("token must not be a signal in the shell markup")
	}
	if strings.Contains(page, `data-bind`) {
		t.Fatal("token field must not be data-bound into a signal")
	}
}

func containsInShell(t *testing.T, want string) bool {
	t.Helper()
	js, err := os.ReadFile(filepath.Join(dist(t), "shell.js"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Contains(string(js), want)
}

func TestShellJSKeepsTokenOffSignals(t *testing.T) {
	js, err := os.ReadFile(filepath.Join(dist(t), "shell.js"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(js)
	if strings.Contains(src, "$token") {
		t.Fatal("shell.js assigns a token signal")
	}
	if !strings.Contains(src, "Authorization") {
		t.Fatal("shell.js must set the Authorization header")
	}
}

func TestDatastarPin(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(dist(t), "datastar.js"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(b), "// Datastar v1.0.4\n") {
		t.Fatal("unexpected datastar bundle header")
	}
	sum := sha256.Sum256(b)
	// Vendored from starfederation/datastar v1.0.4, source map line removed.
	const want = "27fb94cd95af4bc2d5039223df237debb8e5828e881e8334903e391851c50d23"
	got := hex(sum[:])
	if got != want {
		t.Fatalf("datastar.js sha256 %s, want %s", got, want)
	}
}

func hex(b []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2] = digits[c>>4]
		out[i*2+1] = digits[c&0x0f]
	}
	return string(out)
}
