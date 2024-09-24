package handlers

import (
	"fmt"
	"net/http"

	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/views"
)

// JoinPollPage serves the join poll page.
func JoinPollPage(rw http.ResponseWriter, _ *http.Request) {
	_= views.JoinPage("", "").Render(rw)
}

// JoinExistingPoll tries to join an existing poll.
func (hc *HandlerContext) JoinExistingPoll(rw http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	p, err := hc.store.LoadByCode(code)
	if err != nil {
		_= components.JoinForm(code, err.Error()).Render(rw)
	}
	if p == nil {
		panic("store.LoadByCode returned nil poll")
	}
	rw.Header().Set("HX-Location", fmt.Sprintf("/vote/%s/?code=true", p.ID()))
	rw.WriteHeader(http.StatusNoContent)
}
