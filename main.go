// main.go
package main

import (
	"fmt"

	"github.com/TimLai666/go-vdom/vdom"
)

func main() {
	Card := vdom.Component(
		vdom.Div(
			vdom.Props{"class": "card"},
			vdom.H1("{{title}}"),
			vdom.Div("{{children}}"),
		),
	)

	html := vdom.Render(
		Card(vdom.Props{"title": "Hello Card"},
			vdom.Div("這是卡片的內容"),
			vdom.Div("這是卡片的內容2"),
		),
	)

	fmt.Println(html)
}
