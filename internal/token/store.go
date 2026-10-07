package token

import "errors"

// ErrNotFound means the keychain has no token stored.
var ErrNotFound = errors.New("token not found")

// Store reads and writes the bearer token. Tests pass a fake.
// The page never writes the token to disk; only this process does.
type Store interface {
	Get() (string, error)
	Set(token string) error
}
