package components

import (
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
