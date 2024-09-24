package handlers

import "github.com/knumor/qpoll/storage"

// HandlerContext is a common context struct for handler methods
type HandlerContext struct {
	store storage.Storage
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(store storage.Storage) *HandlerContext {
	return &HandlerContext{store: store}
}
