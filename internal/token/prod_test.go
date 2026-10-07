//go:build production

package token

import "testing"

func TestDevTokenIgnoredInProduction(t *testing.T) {
	t.Setenv("DESKTOP_DEV_TOKEN", "dev-secret")
	if got := DevToken(); got != "" {
		t.Fatalf("production build returned %q", got)
	}
}
