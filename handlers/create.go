package handlers

import (
	"net/http"
)

// CreatePage serves the create page.
func (hc *HandlerContext) CreatePage(rw http.ResponseWriter, _ *http.Request) {
	_= hc.pages.CreatePage().Render(rw)
}


