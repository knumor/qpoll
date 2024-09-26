package components

import (
	"fmt"

	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func calcWordSize(weight float64, freq, count int) int {
	if freq <= count/8 {
		return int(weight * 4)
	}
	if freq <= count/6 {
		return int(weight * 5)
	}
	if freq <= count/4 {
		return int(weight * 6)
	}
	return int(weight * 8)
}

// WordCloud is a dynamic word cloud component.
func WordCloud(id string, words []models.Word) g.Node {
	return Div(ID("wordcloud"),
		Class("max-w-lg !mt-[15svh] [&>ul>li>p]:transition-all [&>ul>li>p]:duration-[1500ms] [&>ul>li>p]:ease-in"),
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
						size := calcWordSize(word.Weight, word.Freq, len(words) )
						return Li(P(ID(fmt.Sprintf("wc-item-%d", word.Index)),
							c.Classes{
								"text-xl":                      size <= 1,
								fmt.Sprintf("text-%dxl", size): size > 1,
								"text-sky-600":                 word.Index%5 == 0,
								"text-slate-700":               word.Index%5 == 1,
								"text-orange-700":              word.Index%5 == 2,
								"text-green-600":               word.Index%5 == 3,
								"text-pink-600":                word.Index%5 == 4,
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
