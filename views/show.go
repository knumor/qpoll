package views

import (
	"fmt"
	"log/slog"

	"github.com/knumor/qpoll/components"
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// ShowWordCloudPage is the page to show a word cloud.
func (pr *PageCollection) ShowWordCloudPage(wc *models.WordCloud) g.Node {
	codeStr := fmt.Sprintf("%s", wc.Code())
	slog.Info("ShowWordCloudPage", "code", codeStr)
	return components.Page(
		"Word Cloud",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] items-center space-y-4"),
			components.PollCounter(wc.ResponseCount(), wc.VoteCount()),
			H1(Class("text-4xl text-sky-700"), g.Text(wc.Question())),
			joinHeader(wc.ID(), codeStr),
			components.WordCloud(wc.ID(), wc.GetWords()),
		),
	)
}

// ShowMultipleChoicePage is the page to show a multiple choice qpoll.
func (pr *PageCollection) ShowMultipleChoicePage(mc *models.MultipleChoice) g.Node {
	codeStr := fmt.Sprintf("%s", mc.Code())
	slog.Info("ShowMultipleChoicePage", "code", codeStr)
	return components.Page(
		"Multple Choice",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] space-y-4"),
			// components.PollCounter(mc.ResponseCount(), mc.VoteCount()),
			H1(Class("text-4xl text-sky-700 text-center"), g.Text(mc.Question())),
			joinHeader(mc.ID(), codeStr),
			components.MultipleChoiceResults(mc.ID(), mc.GetOptions()),
		),
	)
}

func joinHeader(id, codeStr string) g.Node {
	return g.Group(
		[]g.Node{
			H3(
				Class("text-sky-900 text-center"),
				g.Text("To join, go to qpoll.mk.priv.no and enter the code  "),
				Span(Class("font-bold text-sky-700 text-2xl tracking-widest py-1 px-1 rounded bg-slate-300"),
					g.Text(codeStr[0:4]),
					g.Text(" "),
					g.Text(codeStr[4:8]),
				),
			),
			Img(Class("w-48 absolute top-10 right-10 hidden md:block"), Src(fmt.Sprintf("/qr/%s", id))),
		},
	)
}
