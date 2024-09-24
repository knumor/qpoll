package views

import (
	"fmt"

	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// ShowWordCloudPage is the page to show a word cloud.
func ShowWordCloudPage(wc *models.WordCloud) g.Node {
	return components.Page(
		"Word Cloud",
		false,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text(wc.Question())),
			H3(
				Class("text-sky-900 text-center"),
				g.Text(fmt.Sprintf("To join, go to qpoll.io/join and enter the code %s", wc.Code())),
			),
			components.WordCloud(wc.ID(), wc.GetWords()),
		),
	)
}
