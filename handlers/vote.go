package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// VotePage serves the vote page.
func (hc *HandlerContext) VotePage(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := hc.store.Load(id)
	if err != nil {
		slog.Error("Failed to load poll", "error", err)
		http.Error(rw, "failed to load poll", http.StatusInternalServerError)
		return
	}
	_ = views.VotePage(p.ID(), p.Question(), false).Render(rw)
}

// VoteSubmit submits a vote.
func (hc *HandlerContext) VoteSubmit(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	p, err := hc.store.Load(id)
	if err != nil {
		slog.Error("VoteSubmit: Failed to load poll", "error", err)
		http.Error(rw, "failed to load poll", http.StatusInternalServerError)
		return
	}
	wc, ok := p.(*models.WordCloud)
	if !ok {
		slog.Error("VoteSubmit: Invalid poll type", "type", fmt.Sprintf("%T", p))
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
		return
	}
	words := r.PostForm["words"]
	slog.Info("Voted for words", "words", words)
	for _, w := range words {
		if w == "" {
			continue
		}
		wc.AddWord(w)
	}
	_ = hc.store.Save(wc)
	_ = views.VotePage(p.ID(), p.Question(), true).Render(rw)
}
