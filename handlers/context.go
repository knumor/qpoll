package handlers

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// Storage is an interface for poll storage.
type Storage interface {
	Save(p models.Poll) error
	Load(id string) (models.Poll, error)
	LoadByCode(code string) (models.Poll, error)
	LoadAllByUser(username string) ([]models.Poll, error)
	Close()
}

// HandlerContext is a common context struct for handler methods
type HandlerContext struct {
	store    Storage
	sessions *scs.SessionManager
	pages    *views.PageCollection
}

// NewHandlerContext creates a new handler context
func NewHandlerContext(store Storage, sessions *scs.SessionManager) *HandlerContext {
	return &HandlerContext{store: store, sessions: sessions, pages: &views.PageCollection{}}
}

// EnsureClientID is a middleware that ensures a client ID is set in the session
func (hc *HandlerContext) EnsureClientID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientID := hc.sessions.GetString(r.Context(), "clientID")
		if clientID == "" {
			clientID = uuid.NewString()
			slog.Info("Generated new client ID", "clientID", clientID)
			hc.sessions.Put(r.Context(), "clientID", clientID)
		}
		next.ServeHTTP(rw, r)
	})
}
