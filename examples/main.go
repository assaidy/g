package main

import (
	"fmt"
	"github.com/assaidy/gtml"
)

func main() {
	// Example using the new API with Add() method and key-value pairs
	nav := gtml.Ul(gtml.KV{"id": "some-id", "class": "class-a class-b", "boolean-tag": true}).Add(
		gtml.Li().Add(gtml.A(gtml.KV{"href": "#home"}).Add(gtml.Text("Home"))),
		gtml.Li().Add(gtml.A(gtml.KV{"href": "#about"}).Add(gtml.Text("About"))),
		gtml.Li().Add(gtml.A(gtml.KV{"href": "#contact"}).Add(gtml.Text("Contact"))),
	)

	html, err := nav.Render()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(html)
	}

	// More complex example
	page := gtml.Html(gtml.KV{"lang": "en"}).Add(
		gtml.Head().Add(
			gtml.Meta(gtml.KV{"charset": "UTF-8"}),
			gtml.Meta(gtml.KV{"name": "viewport", "content": "width=device-width: initial-scale=1.0"}),
			gtml.Title().Add(gtml.Text("GTML New API Example")),
		),
		gtml.Body(gtml.KV{"class": "container"}).Add(
			gtml.Header().Add(
				gtml.H1().Add(gtml.Text("Welcome to GTML")),
				nav,
			),
			gtml.Main().Add(
				gtml.Section(gtml.KV{"class": "content"}).Add(
					gtml.H2().Add(gtml.Text("Introduction")),
					gtml.P(gtml.KV{"class": "lead"}).Add(gtml.Text("This uses the new GTML API with Add() methods.")),
					gtml.P().Add(
						gtml.Text("Key features: "),
						gtml.Strong().Add(gtml.Text("Type safety")),
						gtml.Text(" and "),
						gtml.Em().Add(gtml.Text("clean syntax")),
					),
				),
				gtml.Section(gtml.KV{"class": "form-example"}).Add(
					gtml.H3().Add(gtml.Text("Contact Form")),
					gtml.Form(gtml.KV{"method": "POST", "action": "/submit"}).Add(
						gtml.Div(gtml.KV{"class": "form-group"}).Add(
							gtml.Label(gtml.KV{"for": "name"}).Add(gtml.Text("Name:")),
							gtml.Br(),
							gtml.Input(gtml.KV{"type": "text", "id": "name", "name": "name", "required": true}),
						),
						gtml.Div(gtml.KV{"class": "form-group"}).Add(
							gtml.Label(gtml.KV{"for": "email"}).Add(gtml.Text("Email:")),
							gtml.Br(),
							gtml.Input(gtml.KV{"type": "email", "id": "email", "name": "email", "required": true}),
						),
						gtml.Div(gtml.KV{"class": "form-group"}).Add(
							gtml.Button(gtml.KV{"type": "submit"}).Add(gtml.Text("Submit")),
						),
					),
				),
			),
			gtml.Footer().Add(
				gtml.P().Add(gtml.Text("© 2024 GTML Library")),
				gtml.Hr(),
				gtml.P(gtml.KV{"class": "small"}).Add(
					gtml.Text("Built with "),
					gtml.Code().Add(gtml.Text("gtml")),
				),
			),
		),
	)

	html, err = page.Render()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(html)
	}
}
