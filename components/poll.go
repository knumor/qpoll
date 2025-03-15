package components

import (
	"fmt"

	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

// PollCounter is a component to show the number of responses.
func PollCounter(respCount, voteCount int) g.Node {
	return Div(ID("poll-counter"),
		hx.SwapOOB("poll-counter"),
		Class("absolute top-30 left-10 text-gray-400 text-lg transition-all"),
		g.If(respCount >= 0, P(g.Textf("%d responses", respCount))),
		P(g.Textf("%d votes", voteCount)),
	)
}

// PollList is a list of polls.
func PollList(polls []models.Poll, csrfToken string) g.Node {
	return Div(ID("poll-list"),
		Class("grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3"),
		g.Group(
			g.Map(polls, func(p models.Poll) g.Node {
				return PollCard(p, csrfToken)
			}),
		),
	)
}

// PollCard is a card for a poll.
func PollCard(p models.Poll, csrfToken string) g.Node {
	return Div(Class("border border-gray-300 rounded-md shadow-sm"),
		Div(Class("p-4"),
			H2(Class("text-xl font-bold"), g.Text(p.Question())),
			P(Class("text-gray-600"), g.Text(p.Type().String())),
			P(Class("text-gray-600"), g.Text(p.CreatedAt().String())),
			P(Class("text-gray-600"), g.Textf("%d votes", p.VoteCount())),
			A(Href(fmt.Sprintf("/present/%s", p.ID())),
				Button(Class("bg-sky-700 hover:bg-sky-800 text-white font-bold py-2 px-4 rounded mt-4"),
					g.Text("Present poll"),
				),
			),
			Button(Class("bg-sky-700 hover:bg-sky-800 text-white font-bold ml-2 py-2 px-4 rounded mt-4"),
				hx.Post(fmt.Sprintf("/reset/%s", p.ID())),
				hx.Swap("none"),
				hx.Headers(fmt.Sprintf(`{"X-CSRF-Token": "%s"}`, csrfToken)),
				g.Text("Reset poll"),
			),
		),
	)
}
