package views

import (
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// ListPollsPage is the page to list all polls owned by the user.
func (pr *PageCollection) ListPollsPage(polls []models.Poll) g.Node {
	return Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] items-center space-y-4"),
		H1(Class("text-3xl text-sky-700"), g.Text("My polls")),
		components.PollList(polls),
	)
}
