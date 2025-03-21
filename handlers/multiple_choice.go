package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// CreateMultipleChoicePage serves the create multiple choice page.
func (hc *HandlerContext) CreateMultipleChoicePage(rw http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	user, _ := hc.UserFromSession(r)
	views.Page("Create Multiple Choice", false, user, hc.pages.CreateMultipleChoicePage(csrfToken)).Render(rw)
}

// CreateMultipleChoice creates a multiple choice poll.
func (hc *HandlerContext) CreateMultipleChoice(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.PostForm.Get("question")
	givenOpts := r.PostForm["options"]

	var opts []string
	for _, o := range givenOpts {
		if o == "" {
			continue
		}
		opts = append(opts, o)
	}
	slog.Info("CreateMultipleChoice", "question", q, "options", opts)
	user, err := hc.UserFromSession(r)
	if err != nil {
		http.Error(rw, "Error getting user when creating multiple choice poll", http.StatusBadRequest)
		return
	}
	mc := models.NewMultipleChoice(q, user.Username, opts)
	hc.store.Save(mc)
	http.Redirect(rw, r, "/present/"+mc.ID(), http.StatusSeeOther)
}

// GetMultipleChoice generates and serves multiple choice poll results.
func (hc *HandlerContext) GetMultipleChoice(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, _ := hc.store.Load(id)
	wc, ok := p.(*models.MultipleChoice)
	if !ok {
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
		return
	}
	opts := wc.GetOptions()
	_ = components.PollCounter(-1, wc.VoteCount()).Render(rw)
	_ = components.MultipleChoiceResults(id, opts).Render(rw)
}
