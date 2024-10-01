package views

import (
	"github.com/knumor/qpoll/components"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// CreatePage is the page to create a new qpoll.
func CreatePage() g.Node {
	return components.Page(
		"Home",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("What do you want to create?")),
			Div(
				Class("flex justify-center space-x-4"),
				components.BigButton("Multiple Choice", "/create/multiple-choice"),
				components.BigButton("Word Cloud", "/create/wordcloud"),
			),
		),
	)
}

// CreateWordCloudPage is the page to create a word cloud.
func CreateWordCloudPage() g.Node {
	return components.Page(
		"Create Word Cloud",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("Create a word cloud")),
			Form(
				Method("post"),
				Action("/create/wordcloud"),
				Class("flex flex-col space-y-4 w-1/2"),
				Input(
					Class("p-2 border rounded-md"),
					AutoFocus(),
					Type("text"),
					Name("question"),
					Placeholder("Ask your question here..."),
					Required(),
				),
				components.SubmitButton("Create"),
			),
		),
	)
}

// CreateMultipleChoicePage is the page to create a multiple choice qpoll.
func CreateMultipleChoicePage() g.Node {
	return components.Page(
		"Create Multiple Choice",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("Create a multiple choice poll")),
			Form(
				Method("post"),
				Action("/create/mc"),
				Class("flex flex-col space-y-4 w-1/2"),
				Input(
					Class("p-2 border rounded-md"),
					AutoFocus(),
					Type("text"),
					Name("question"),
					Placeholder("Ask your question here..."),
					Required(),
				),
				P(Class("text-sky-700 self-center"), g.Text("You may provide up to 6 options to choose from:")),
				components.MultipleChoiceInputs("options", 6),
				components.SubmitButton("Create"),
			),
		),
	)
}
