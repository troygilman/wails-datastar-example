//go:build production

package session

import (
	"testing"

	"wails-datastar-example/internal/token"
)

func TestLoadIgnoresDevTokenInProduction(t *testing.T) {
	t.Setenv("DESKTOP_DEV_TOKEN", "dev-secret")
	t.Setenv("DESKTOP_API_BASE", "")
	got, err := Load(&token.Memory{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Token != "" {
		t.Fatalf("production build returned %q", got.Token)
	}
}
