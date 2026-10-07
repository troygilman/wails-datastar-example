//go:build !production

package session

import (
	"testing"

	"wails-datastar-example/internal/token"
)

func TestLoadUsesDevTokenWhenStoreEmpty(t *testing.T) {
	t.Setenv("DESKTOP_DEV_TOKEN", "dev-secret")
	t.Setenv("DESKTOP_API_BASE", "http://api.test/root/")
	got, err := Load(&token.Memory{})
	if err != nil {
		t.Fatal(err)
	}
	if got.APIBase != "http://api.test/root" {
		t.Fatalf("api %q", got.APIBase)
	}
	if got.Token != "dev-secret" {
		t.Fatalf("token %q", got.Token)
	}
}
