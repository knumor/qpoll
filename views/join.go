package views

import (
	"github.com/knumor/qpoll/components"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

// JoinPage is the page to join a qpoll.
func (pr *PageCollection) JoinPage(code, errorMsg, csrfToken string) g.Node {
	return Div(Class("flex flex-col min-h-[calc(100dvh-10rem)] justify-center items-center space-y-4"),
		H1(Class("text-3xl text-sky-700"), g.Text("Enter the code to join")),
		components.JoinForm(code, errorMsg, csrfToken),
	)
}
