package components

import (
	"github.com/knumor/qpoll/models"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

// Page is the base page layout.
func Page(title string, hideHeader bool, authUser models.User, children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("preload"), Href("https://cdn.jsdelivr.net/npm/charts.css/dist/charts.min.css"), As("style")),
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/charts.css/dist/charts.min.css")),
			Link(Rel("preload"), Href("/public/styles.css"), As("style")),
			Link(Rel("stylesheet"), Href("/public/styles.css")),
			Script(g.Text("let FF_FOUC_FIX;")),
			Script(Src("https://cdn.jsdelivr.net/npm/htmx.org")),
		},
		Body: []g.Node{Class("bg-gradient-to-b bg-no-repeat from-white to-slate-300"),
			Div(Class("min-h-screen justify-between flex flex-col"),
				g.If(!hideHeader, header(authUser.Name)),
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

func header(authName string) g.Node {
	return Div(ID("header"), Class("bg-slate-300 uppercase text-sm text-sky-800 shadow"),
		container(false,
			Div(Class("flex h-8 items-center space-x-4"),
				headerLink("/", "Join"),
				headerLink("/create", "Create"),
				g.If(authName != "", headerLink("/polls", "My polls")),
				A(Class("hover:text-sky-500 !ml-auto"),
					g.If(authName != "", g.Group(
						[]g.Node{Href("#"), g.Text(authName)}),
					),
					g.If(authName == "", g.Group(
						[]g.Node{Href("/login"), g.Text("Login")}),
					),
				),
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
		P(g.Text("qpoll.mk.priv.no - a service by @knumor")),
	)
}
