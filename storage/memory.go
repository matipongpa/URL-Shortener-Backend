package storage

import (
	"context"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string]string)}
}

func (m *Memory) Save(ctx context.Context, code, url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[code]; ok {
		return ErrCodeExists
	}
	m.data[code] = url

	return nil
}

func (m *Memory) Get(ctx context.Context, code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.data[code]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}
