package main

import (
	"fmt"
	"os"

	"github.com/assaidy/gtml"
	g "github.com/assaidy/gtml"
)

func main() {
	// Example using the new API with Add() method and key-value pairs
	nav := g.Ul(g.KV{"id": "some-id", "class": "class-a class-b", "boolean-tag": true}).Add(
		g.Li().Add(g.A(g.KV{"href": "#home"}).Add(g.Text("Home"))),
		g.Li().Add(g.A(g.KV{"href": "#about"}).Add(g.Text("About"))),
		g.Li().Add(g.A(g.KV{"href": "#contact"}).Add(g.Text("Contact"))),
	)

	html, err := nav.Render()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(html)
	}

	// More complex example
	page := g.Html(g.KV{"lang": "en"}).Add(
		g.Head().Add(
			g.Meta(g.KV{"charset": "UTF-8"}),
			g.Meta(g.KV{"name": "viewport", "content": "width=device-width: initial-scale=1.0"}),
			g.Title().Add(g.Text("GTML New API Example")),
		),
		g.Body(g.KV{"class": "container"}).Add(
			g.Header().Add(
				g.H1().Add(g.Text("Welcome to GTML")),
				nav,
			),
			g.Main().Add(
				g.Section(g.KV{"class": "content"}).Add(
					g.H2().Add(g.Text("Introduction")),
					g.P(g.KV{"class": "lead"}).Add(g.Text("This uses the new GTML API with Add() methods.")),
					g.P().Add(
						g.Text("Key features: "),
						g.Strong().Add(g.Text("Type safety")),
						g.Text(" and "),
						g.Em().Add(g.Text("clean syntax")),
					),
				),
				g.Section(g.KV{"class": "form-example"}).Add(
					g.H3().Add(g.Text("Contact Form")),
					g.Form(g.KV{"method": "POST", "action": "/submit"}).Add(
						g.Div(g.KV{"class": "form-group"}).Add(
							g.Label(g.KV{"for": "name"}).Add(g.Text("Name:")),
							g.Br(),
							g.Input(g.KV{"type": "text", "id": "name", "name": "name", "required": true}),
						),
						g.Div(g.KV{"class": "form-group"}).Add(
							g.Label(g.KV{"for": "email"}).Add(g.Text("Email:")),
							g.Br(),
							g.Input(g.KV{"type": "email", "id": "email", "name": "email", "required": true}),
						),
						g.Div(g.KV{"class": "form-group"}).Add(
							g.Button(g.KV{"type": "submit"}).Add(g.Text("Submit")),
						),
					),
				),
			),
			g.Footer().Add(
				g.P().Add(g.Text("© 2026 GTML Library")),
				g.Hr(),
				g.P(g.KV{"class": "small"}).Add(
					g.Text("Built with "),
					g.Code().Add(g.Text("gtml")),
				),
			),
		),
	)

	if err := gtml.Render(os.Stdout, page); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
