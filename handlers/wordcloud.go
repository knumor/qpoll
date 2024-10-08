package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	"github.com/knumor/qpoll/views"
)

// CreateWordCloudPage serves the create word cloud page.
func (hc *HandlerContext) CreateWordCloudPage(rw http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	user, _ := hc.UserFromSession(r)
	views.Page("Create a Word Cloud", false, user, hc.pages.CreateWordCloudPage(csrfToken)).Render(rw)
}

// CreateWordCloud creates a word cloud.
func (hc *HandlerContext) CreateWordCloud(rw http.ResponseWriter, r *http.Request) {
	q := r.FormValue("question")
	user, err := hc.UserFromSession(r)
	if err != nil {
		http.Error(rw, "Error getting user when creating wordcloud", http.StatusBadRequest)
		return
	}
	wc := models.NewWordCloud(q, user.Username)
	hc.store.Save(wc)
	http.Redirect(rw, r, "/present/"+wc.ID(), http.StatusSeeOther)
}

// GetWordCloud generates and serves a word cloud.
func (hc *HandlerContext) GetWordCloud(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, _ := hc.store.Load(id)
	wc, ok := p.(*models.WordCloud)
	if !ok {
		http.Error(rw, "invalid poll type", http.StatusBadRequest)
		return
	}
	words := wc.GetWords()
	_ = components.PollCounter(wc.ResponseCount(), wc.VoteCount()).Render(rw)
	_ = components.WordCloud(id, words).Render(rw)
}
