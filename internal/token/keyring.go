package token

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "wails-datastar-example"
	keyringUser    = "bearer"
)

// Keyring stores the bearer token in the OS keychain.
type Keyring struct{}

func (Keyring) Get() (string, error) {
	val, err := keyring.Get(keyringService, keyringUser)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (Keyring) Set(token string) error {
	return keyring.Set(keyringService, keyringUser, token)
}
