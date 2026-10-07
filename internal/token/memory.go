package token

import "sync"

// Memory is an in-process store for tests.
type Memory struct {
	mu  sync.Mutex
	val string
	ok  bool
}

func (m *Memory) Get() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.ok {
		return "", ErrNotFound
	}
	return m.val, nil
}

func (m *Memory) Set(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.val = token
	m.ok = true
	return nil
}
