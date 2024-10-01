package components

import (
	"fmt"

	"github.com/knumor/qpoll/models"
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

func voteTextInput(placehoder, name string, autofocus bool) g.Node {
	return Input(
		Class("p-2 border rounded-md text-lg"),
		Type("text"),
		Name(name),
		Placeholder(placehoder),
		MaxLength("25"),
		g.If(autofocus, AutoFocus()),
	)
}

func voteRadioInput(id, label, name, value string) g.Node {
	return Div(Class("space-x-2 border p-2 border-slate-300 rounded-md"),
		Input(
			Class("accent-sky-700"),
			ID(id),
			Type("radio"),
			Value(value),
			Name(name),
		),
		Label(
			Class("text-lg ml-2"),
			For(fmt.Sprintf(id)),
			g.Text(label),
		),
		// Br(),
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
		voteTextInput("Enter a word", "words", true),
		voteTextInput("Enter another word", "words", false),
		voteTextInput("Enter another word", "words", false),
		g.If(redo, P(
			Class("text-sky-700 w-48 text-center self-center"),
			g.Text("Thank you for voting! Feel free to vote again above."),
		)),
		SubmitButton("Vote"),
	)
}

// MultipleChoiceInputs is the form inputs for multiple choice.
func MultipleChoiceInputs(name string, count int) g.Node {
	inputs := make([]g.Node, count)
	for i := 0; i < count; i++ {
		inputs[i] = voteTextInput(fmt.Sprintf("Option %d", i+1), name, i == 0)
	}
	return g.Group(
		inputs,
	)
}

// MultipleChoiceVoteForm is the form to vote on a multiple choice qpoll.
func MultipleChoiceVoteForm(id string, options []models.Option) g.Node {
	return Form(
		Method("post"),
		Action("/vote"),
		Class("flex flex-col space-y-4 w-fit"),
		Div(Class("flex flex-wrap gap-2 space-y-4 text-sky-700"),
			Input(
				Name("id"),
				Value(id),
				Type("hidden"),
			),
			g.Group(
				g.Map(options, func(o models.Option) g.Node {
					return voteRadioInput(
						fmt.Sprintf("mc-option-%d", o.Index),
						o.Text,
						"choice",
						fmt.Sprintf("%d", o.Index),
					)
				}),
			),
		),
		SubmitButton("Vote"),
	)
}
