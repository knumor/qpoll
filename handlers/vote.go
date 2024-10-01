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
	switch p.Type() {
	case models.WordCloudPoll:
		_ = views.WordsVotePage(p.ID(), p.Question(), false).Render(rw)
	case models.MultipleChoicePoll:
		mc := p.(*models.MultipleChoice)
		_ = views.MultipleChoiceVotePage(mc.ID(), mc.Question(), mc.GetOptions()).Render(rw)
	default:
		slog.Error("Invalid poll type", "type", p.Type())
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
	}
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
	switch p.Type() {
	case models.WordCloudPoll:
		words := r.PostForm["words"]
		slog.Info("Voted for words", "words", words)
		wc := p.(*models.WordCloud)
		var cnt int
		for _, w := range words {
			if w == "" {
				continue
			}
			wc.AddWord(w)
			cnt++
		}
		wc.AddVote(cnt)
	_ = hc.store.Save(wc)
	_ = views.WordsVotePage(p.ID(), p.Question(), true).Render(rw)
	case models.MultipleChoicePoll:
		mc := p.(*models.MultipleChoice)
		choice := r.PostForm.Get("choice")
		idx := int(choice[0] - '0')
		slog.Info("Voted for choice", "choice", idx)
		mc.AddVoteForOption(idx)
		_ = hc.store.Save(mc)
		_ = views.ThankYouPage().Render(rw)
	default:
		slog.Error("VoteSubmit: Invalid poll type", "type", fmt.Sprintf("%T", p))
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
		return
	}
}
