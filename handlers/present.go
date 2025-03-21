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
	user, _ := hc.UserFromSession(r)
	switch p.Type() {
	case models.WordCloudPoll:
		views.Page(
			"Word Cloud",
			false,
			user,
			hc.pages.ShowWordCloudPage(p.(*models.WordCloud)),
		).Render(rw)
	case models.MultipleChoicePoll:
		views.Page(
			"Multiple Choice",
			false,
			user,
			hc.pages.ShowMultipleChoicePage(p.(*models.MultipleChoice)),
		).Render(rw)
	default:
		slog.Error("PresentPoll: Invalid poll type", "type", p.Type())
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
	}
}
