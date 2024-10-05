package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
)

// CreateWordCloudPage serves the create word cloud page.
func (hc *HandlerContext) CreateWordCloudPage(rw http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	_ = hc.pages.CreateWordCloudPage(csrfToken).Render(rw)
}

// CreateWordCloud creates a word cloud.
func (hc *HandlerContext) CreateWordCloud(rw http.ResponseWriter, r *http.Request) {
	q := r.FormValue("question")
	wc := models.NewWordCloud(q)
	hc.store.Save(wc)
	http.Redirect(rw, r, "/present/"+wc.ID(), http.StatusSeeOther)
}

// GetWordCloud generates and servers a word cloud.
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
