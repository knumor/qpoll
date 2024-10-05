package views

import (
	"github.com/knumor/qpoll/components"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// LoginPage is the page to log in.
func (pr *PageCollection) LoginPage(csrfToken, errorMsg, returnTo string) g.Node {
	return components.Page(
		"Log in",
		false,
		pr.AuthenticatedUser,
		Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
			H1(Class("text-3xl text-sky-700"), g.Text("Log in")),
			components.LoginForm(csrfToken, errorMsg, returnTo),
		),
	)
}
