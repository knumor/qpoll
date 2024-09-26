package components

import (
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

// JoinForm is the form to join a qpoll.
func JoinForm(value, errorMsg string) g.Node {
	return Form(
		hx.Post("/join"),
		hx.Target("this"),
		hx.Swap("outerHTML"),
		Class("flex flex-col space-y-4 w-fit"),
		Input(
			c.Classes{
				"p-2 border rounded-md text-lg": true,
				"border-red-700 border-2":       errorMsg != "",
			},
			Name("code"),
			g.If(
				value != "",
				Value(value),
			), AutoFocus(), Type("text"), Placeholder("1234 5678"), Required()),
		g.If(
			errorMsg != "",
			Div(Class("text-red-700"), g.Text(errorMsg)),
		),
		SubmitButton("Join"),
	)
}

func voteInput(placehoder string, autofocus bool) g.Node {
	return Input(
		Class("p-2 border rounded-md text-lg"),
		Type("text"),
		Name("words"),
		Placeholder(placehoder),
		MaxLength("25"),
		g.If(autofocus, AutoFocus()),
	)
}

// WordsVoteForm is the form to vote on a qpoll.
func WordsVoteForm(id string, redo bool) g.Node {
	return Form(
		Method("post"),
		Action("/vote"),
		Class("flex flex-col space-y-4 w-fit"),
		Input(
			Name("id"),
			Value(id),
			Type("hidden"),
		),
		voteInput("Enter a word", true),
		voteInput("Enter another word", false),
		voteInput("Enter another word", false),
		g.If(redo, P(
			Class("text-sky-700 w-48 text-center self-center"),
			g.Text("Thank you for voting! Feel free to vote again above."),
		)),
		SubmitButton("Vote"),
	)
}
