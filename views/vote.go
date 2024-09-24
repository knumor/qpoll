package views

import (
	"github.com/knumor/qpoll/components"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// VotePage is the page to vote on a qpoll.
func VotePage(id, question string) g.Node {
	return components.Page(
		"Vote",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text(question)),
			components.VoteForm(id),
		),
	)
}

// ThankYouPage is the page to thank the user for voting.
func ThankYouPage() g.Node {
	return components.Page(
		"Thank you",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("Thank you for voting!")),
		),
	)
}
