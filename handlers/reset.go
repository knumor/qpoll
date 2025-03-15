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
	p.Reset()
	if err := hc.store.Save(p); err != nil {
		slog.Error("ResetPoll: Failed to save poll", "error", err)
		http.Error(rw, "failed to save poll", http.StatusInternalServerError)
		return
	}
	rw.Header().Add("HX-Refresh", "true")
}
