package views

import (
	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// WordsVotePage is the page to vote on a qpoll.
func (pr *PageCollection) WordsVotePage(id, question string, redo bool, csrfToken string) g.Node {
	return components.Page(
		"Vote",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text(question)),
			components.WordsVoteForm(id, redo, csrfToken),
		),
	)
}

// MultipleChoiceVotePage is the page to vote on a multiple choice qpoll.
func (pr *PageCollection) MultipleChoiceVotePage(id, question string, options []models.Option, csrfToken string) g.Node {
	return components.Page(
		"Vote",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text(question)),
			components.MultipleChoiceVoteForm(id, options, csrfToken),
		),
	)
}

// ThankYouPage is the page to thank the user for voting.
func (pr *PageCollection) ThankYouPage() g.Node {
	return components.Page(
		"Thank you",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("Thank you for voting!")),
		),
	)
}
