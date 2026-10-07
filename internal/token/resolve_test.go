package token

import (
	"errors"
	"testing"
)

type errStore struct{ err error }

func (s errStore) Get() (string, error) { return "", s.err }
func (s errStore) Set(string) error     { return nil }

func TestResolvePrefersStoredToken(t *testing.T) {
	store := &Memory{}
	if err := store.Set("stored"); err != nil {
		t.Fatal(err)
	}
	got, err := Resolve(store, "from-env")
	if err != nil {
		t.Fatal(err)
	}
	if got != "stored" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveUsesDevTokenWhenMissing(t *testing.T) {
	got, err := Resolve(&Memory{}, "from-env")
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-env" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveEmptyWhenNothingStored(t *testing.T) {
	got, err := Resolve(&Memory{}, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveFallsBackToDevWhenStoreFails(t *testing.T) {
	got, err := Resolve(errStore{err: errors.New("dbus down")}, "from-env")
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-env" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveReturnsStoreErrorWithoutDevToken(t *testing.T) {
	_, err := Resolve(errStore{err: errors.New("dbus down")}, "")
	if err == nil || err.Error() != "dbus down" {
		t.Fatalf("got %v", err)
	}
}
