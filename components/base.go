package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

// Page is the base page layout.
func Page(title string, hideHeader bool, children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("preload"), Href("https://unpkg.com/charts.css/dist/charts.min.css"), As("style")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/charts.css/dist/charts.min.css")),
			Link(Rel("preload"), Href("/public/styles.css"), As("style")),
			Link(Rel("stylesheet"), Href("/public/styles.css")),
			Script(g.Text("let FF_FOUC_FIX;")),
			Script(Src("https://unpkg.com/htmx.org")),
		},
		Body: []g.Node{Class("bg-gradient-to-b bg-no-repeat from-white to-slate-300"),
			Div(Class("min-h-screen justify-between flex flex-col"),
				g.If(!hideHeader, header()),
				Div(Class("grow"),
					container(true,
						Div(Class(""),
							g.Group(children),
						),
					),
				),
				footer(),
			),
		},
	})
}

func header() g.Node {
	return Div(ID("header"), Class("bg-slate-300 uppercase text-sm text-sky-800 shadow"),
		container(false,
			Div(Class("flex h-8 items-center space-x-4"),
				headerLink("/", "Home"),
				headerLink("/create", "Create"),
				A(Class("hover:text-sky-500 !ml-auto"), Href("/"), g.Text("Morten")),
			),
		),
	)
}

func headerLink(href, text string) g.Node {
	return A(Class("hover:text-sky-500"), Href(href), g.Text(text))
}

func container(padY bool, children ...g.Node) g.Node {
	return Div(
		c.Classes{
			"max-w-7xl mx-auto":     true,
			"px-4 md:px-8 lg:px-16": true,
			"py-4 md:py-8":          padY,
		},
		g.Group(children),
	)
}

func footer() g.Node {
	return Div(Class("text-slate-400 justify-center flex h-8 items-center text-center text-sm"),
		P(g.Text("qpoll.io - a service by knumor")),
	)
}
