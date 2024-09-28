package storage

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/knumor/qpoll/handlers"
	"github.com/knumor/qpoll/models"
)

type memstore struct {
	sync.RWMutex
	pollsByID   map[string]models.Poll
	pollsByCode map[string]models.Poll
}

// NewMemStore creates a new in-memory storage.
func NewMemStore() handlers.Storage {
	return &memstore{
		pollsByID:   make(map[string]models.Poll),
		pollsByCode: make(map[string]models.Poll),
	}
}

// Save saves the poll to the storage.
func (ms *memstore) Save(p models.Poll) error {
	ms.Lock()
	defer ms.Unlock()
	slog.Info("Save", "id", p.ID())
	slog.Info("Save", "code", p.Code())
	ms.pollsByID[p.ID()] = p
	ms.pollsByCode[p.Code()] = p
	return nil
}

// Load loads the poll from the storage.
func (ms *memstore) Load(id string) (models.Poll, error) {
	ms.RLock()
	defer ms.RUnlock()
	if p, ok := ms.pollsByID[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("poll with id %s not found", id)
}

// LoadByCode loads the poll by code from the storage.
func (ms *memstore) LoadByCode(code string) (models.Poll, error) {
	ms.RLock()
	defer ms.RUnlock()
	slog.Info("LoadByCode", "code", code)
	slog.Info("LoadByCode", "pollsByCodeLen", len(ms.pollsByCode))
	for k, v := range ms.pollsByCode {
		slog.Info("LoadByCode", "key", k)
		slog.Info("LoadByCode", "value", v)
	}
	if p, ok := ms.pollsByCode[code]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("poll with code %s not found", code)
}

func (ms *memstore) Close() {
}
