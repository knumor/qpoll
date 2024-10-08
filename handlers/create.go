package handlers

import (
	"net/http"

	"github.com/knumor/qpoll/views"
)

// CreatePage serves the create page.
func (hc *HandlerContext) CreatePage(rw http.ResponseWriter, r *http.Request) {
	user, _ := hc.UserFromSession(r)
	views.Page("Create a new poll", false, user, hc.pages.CreatePage()).Render(rw)
}
