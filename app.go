package main

import (
	"wails-datastar-example/internal/session"
	"wails-datastar-example/internal/token"
)

// App is the bound shell. The page calls Session once, then talks to
// the remote Datastar server directly. This process does not proxy that traffic.
type App struct {
	store token.Store
}

func NewApp(store token.Store) *App {
	if store == nil {
		store = token.Keyring{}
	}
	return &App{store: store}
}

// Session returns the API base URL and the bearer token.
// DESKTOP_API_BASE overrides the default http://127.0.0.1:8080.
// DESKTOP_DEV_TOKEN is read only when this binary was not built with the production tag.
func (a *App) Session() (session.Session, error) {
	return session.Load(a.store)
}

// SaveToken stores a pasted bearer token in the OS keychain.
func (a *App) SaveToken(raw string) error {
	return session.Save(a.store, raw)
}
