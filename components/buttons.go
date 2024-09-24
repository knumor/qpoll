package components

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// BigButton is a large buttons with a link
func BigButton(text string, href string) g.Node {
	return A(Href(href),
		Button(
			Class("bg-sky-700 hover:bg-sky-800 text-white font-bold py-2 px-4 rounded"),
			g.Text(text),
		),
	)
}

// SubmitButton is a submit button
func SubmitButton(text string) g.Node {
	return Button(
		Class("bg-sky-700 hover:bg-sky-800 text-white font-bold py-2 px-4 rounded self-center"),
		Type("submit"),
		g.Text(text),
	)
}

