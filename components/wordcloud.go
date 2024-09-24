package components

import (
	"fmt"

	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

// WordCloud is a dynamic word cloud component.
func WordCloud(id string, words []models.Word) g.Node {
	return Div(ID("wordcloud"),
		Class("w-96 !mt-[15svh] [&>ul>li>p]:transition-all [&>ul>li>p]:duration-[1500ms] [&>ul>li>p]:ease-in"),
		hx.Get("/wordcloud/"+id),
		hx.Trigger("every 2s"),
		hx.Target("this"),
		hx.Swap("outerHTML"),
		g.If(len(words) == 0, P(Class("text-center text-sm text-gray-400"), g.Text("No answers yet"))),
		Ul(ID("wc-ul"), Class("flex justify-center flex-wrap align-center gap-2 leading-8"),
			g.Group(
				g.Map(
					words,
					func(word models.Word) g.Node {
						return Li(P(ID(fmt.Sprintf("wc-item-%d", word.Index)),
							c.Classes{
								"text-2xl":         word.Weight < 0.1,
								"text-3xl":         word.Weight >= 0.1 && word.Weight < 0.3,
								"text-4xl":        word.Weight >= 0.3 && word.Weight < 0.5,
								"text-5xl":        word.Weight >= 0.5 && word.Weight < 0.6,
								"text-6xl":        word.Weight >= 0.6 && word.Weight < 0.7,
								"text-7xl":        word.Weight >= 0.7 && word.Weight < 0.8,
								"text-8xl":        word.Weight >= 0.8,
								"text-sky-600":    word.Index%5 == 0,
								"text-slate-700":  word.Index%5 == 1,
								"text-orange-700": word.Index%5 == 2,
								"text-green-600":  word.Index%5 == 3,
								"text-pink-600":   word.Index%5 == 4,
								// "[writing-mode:sideways-lr]":   word.Index%6 == 5,
								// "[writing-mode:sideways-rl]":   word.Index%6 == 3,
							},
							g.Text(word.Text)))
					},
				),
			),
		),
	)
}
