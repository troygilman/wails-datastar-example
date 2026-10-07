package token

import "errors"

// Resolve returns the token to send as Authorization.
// A stored token wins. dev is used only when the store is empty,
// or when the store itself fails and a dev token was compiled in.
// Production builds pass an empty dev string.
func Resolve(store Store, dev string) (string, error) {
	tok, err := store.Get()
	switch {
	case err == nil && tok != "":
		return tok, nil
	case err == nil || errors.Is(err, ErrNotFound):
		return dev, nil
	case dev != "":
		return dev, nil
	default:
		return "", err
	}
}
