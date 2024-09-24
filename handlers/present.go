package handlers

import (
	"log/slog"
	"net/http"

	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// PresentPoll serves the poll page.
func (hc *HandlerContext) PresentPoll(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := hc.store.Load(id)
	if err != nil {
		slog.Error("PresentPoll: Failed to load poll", "error", err)
		http.Error(rw, "failed to load poll", http.StatusInternalServerError)
		return
	}
	switch p.Type() {
	case models.WordCloudPoll:
		_ = views.ShowWordCloudPage(p.(*models.WordCloud)).Render(rw)
	case models.MultipleChoicePoll:
		http.Error(rw, "multiple choice polls not implemented", http.StatusNotImplemented)
	default:
		slog.Error("PresentPoll: Invalid poll type", "type", p.Type())
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
	}
}
