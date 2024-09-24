package handlers

import (
	"net/http"

	"github.com/knumor/qpoll/views"
)

// CreatePage serves the create page.
func CreatePage(rw http.ResponseWriter, _ *http.Request) {
	_= views.CreatePage().Render(rw)
}


