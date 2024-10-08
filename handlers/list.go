package handlers

import (
	"log/slog"
	"net/http"

	"github.com/knumor/qpoll/views"
)

// ListPollsPage serves the list polls page.
func (hc *HandlerContext) ListPollsPage(rw http.ResponseWriter, r *http.Request) {
	user, err := hc.UserFromSession(r)
	if err != nil {
		slog.Error("ListPollsPage: Error getting user from session", "error", err)
		http.Error(rw, "Error getting user from session", http.StatusBadRequest)
		return
	}
	polls, err := hc.store.LoadAllByUser(user.Username)
	if err != nil {
		slog.Error("ListPollsPage: Error loading polls", "error", err)
		http.Error(rw, "Error loading polls", http.StatusInternalServerError)
		return
	}
	views.Page("My Polls", false, user, hc.pages.ListPollsPage(polls)).Render(rw)
}

