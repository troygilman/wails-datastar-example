//go:build !production

package token

import "testing"

func TestDevTokenReadsEnv(t *testing.T) {
	t.Setenv("DESKTOP_DEV_TOKEN", "dev-secret")
	if got := DevToken(); got != "dev-secret" {
		t.Fatalf("got %q", got)
	}
}
