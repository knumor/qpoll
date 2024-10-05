package views

import (
	"github.com/knumor/qpoll/models"
)

// PageCollection is a struct with common data and methods to render the various pages.
type PageCollection struct {
	AuthenticatedUser models.User
}
