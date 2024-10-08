package views

import (
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
)

// PageCollection is a struct with common data and methods to render the various pages.
type PageCollection struct {
}

// Page is a wrapper for a page view
func Page(title string, hideHeader bool, user models.User, pageNode g.Node) g.Node {
	return components.Page(
		title,
		hideHeader,
		user,
		pageNode,
	)
}
