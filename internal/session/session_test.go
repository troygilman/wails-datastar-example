package session

import (
	"testing"

	"wails-datastar-example/internal/token"
)

func TestLoadDefaultAPI(t *testing.T) {
	t.Setenv("DESKTOP_API_BASE", "")
	t.Setenv("DESKTOP_DEV_TOKEN", "")
	got, err := Load(&token.Memory{})
	if err != nil {
		t.Fatal(err)
	}
	if got.APIBase != "http://127.0.0.1:8080" {
		t.Fatalf("api %q", got.APIBase)
	}
	if got.Token != "" {
		t.Fatalf("token %q", got.Token)
	}
}

func TestLoadRejectsBadAPI(t *testing.T) {
	t.Setenv("DESKTOP_API_BASE", "not a url")
	if _, err := Load(&token.Memory{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveRejectsEmpty(t *testing.T) {
	if err := Save(&token.Memory{}, "  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveThenLoad(t *testing.T) {
	t.Setenv("DESKTOP_DEV_TOKEN", "from-env")
	store := &token.Memory{}
	if err := Save(store, "pas\nted"); err == nil {
		t.Fatal("newline should be rejected")
	}
	if err := Save(store, " pasted "); err != nil {
		t.Fatal(err)
	}
	got, err := Load(store)
	if err != nil {
		t.Fatal(err)
	}
	if got.Token != "pasted" {
		t.Fatalf("stored token should win, got %q", got.Token)
	}
}
