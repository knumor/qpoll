package handlers

import (
	"github.com/knumor/qpoll/models"
)

// Storage is an interface for poll storage.
type Storage interface {
	Save(p models.Poll) error
	Load(id string) (models.Poll, error)
	LoadByCode(code string) (models.Poll, error)
	Close()
}

// HandlerContext is a common context struct for handler methods
type HandlerContext struct {
	store Storage
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(store Storage) *HandlerContext {
	return &HandlerContext{store: store}
}
