package handlers

import (
	"log/slog"
	"net/http"
)

// ResetPoll resets the given poll, removing all votes.
func (hc *HandlerContext) ResetPoll(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := hc.store.Load(id)
	if err != nil {
		slog.Error("ResetPoll: Failed to load poll", "error", err)
		http.Error(rw, "failed to load poll", http.StatusInternalServerError)
		return
	}
	user, _ := hc.UserFromSession(r)
	if p.Owner() != user.Username {
		slog.Error("ResetPoll: User is not the owner", "user", user.Username, "owner", p.Owner())
		http.Error(rw, "user is not the owner", http.StatusForbidden)
		return
	}
	p.Reset()
	if err := hc.store.Save(p); err != nil {
		slog.Error("ResetPoll: Failed to save poll", "error", err)
		http.Error(rw, "failed to save poll", http.StatusInternalServerError)
		return
	}
	rw.Header().Add("HX-Refresh", "true")
}
