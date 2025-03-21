package components

import (
	"fmt"
	"log/slog"

	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"

	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

// MultipleChoiceResults is a dynamic bar graph showing the results of a multiple choice question.
func MultipleChoiceResults(id string, options []models.Option) g.Node {
	slog.Info("ShowMultipleChoicePage", "id", id, "options", options)
	return Div(ID("mc-results"),
		Class("!mt-[10svh] !h-full [&>table>tbody>tr>td]:transition-all [&>table>tbody>tr>td]:duration-[1500ms] [&>table>tbody>tr>td]:ease-in"),
		hx.Get("/mc/"+id),
		hx.Trigger("every 2s"),
		hx.Target("this"),
		hx.Swap("outerHTML"),
		Table(Class("charts-css column data-outside data-spacing-10 show-labels text-sky-700"),
			g.Group(
				g.Map(options, func(option models.Option) g.Node {
					w := max(option.Weight, 0.01)
					labelBottomOffset := "0rem"
					if len(option.Text) > 25 {
						labelBottomOffset = "-1.5rem"
					}
					thStyle := fmt.Sprintf("bottom:%s", labelBottomOffset)
					return Tr(
						g.Raw(`<th scope="row" style="` + thStyle + `">`),
						g.Text(option.Text),
						g.Raw(`</th>`),
						Td(
							ID(fmt.Sprintf("mc-bar-%d", option.Index)),
							Style(fmt.Sprintf("--size: %.2f", w)),
							Span(Class("data text-3xl"), g.Textf("%d", option.Count)),
						),
					)
				}),
			),
		),
	)
}
