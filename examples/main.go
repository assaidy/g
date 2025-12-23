package main

import (
	"fmt"
	"os"

	"github.com/assaidy/g"
)

func main() {
	// Sample data for demonstration
	items := []string{"Apple", "Banana", "Cherry", "Date"}
	isLoggedIn := true
	userName := "John Doe"

	// Demonstrate all utility functions in a single example
	page := g.Html(g.KV{"lang": "en"}).Add(
		g.Head().Add(
			g.Meta(g.KV{"charset": "UTF-8"}),
			g.Title().Add(g.Text("GTML Utils Example")),
		),
		g.Body().Add(
			// Using IfElse to show conditional content
			g.IfElse(isLoggedIn,
				g.Div(g.KV{"class": "welcome"}).Add(
					g.H1().Add(g.Text(fmt.Sprintf("Welcome back, %s", userName))),
					g.P().Add(g.Text("You are logged in!")),
				),
				g.Div(g.KV{"class": "login-prompt"}).Add(
					g.H1().Add(g.Text("Please log in")),
					g.P().Add(g.Text("You need to authenticate to continue.")),
				),
			),

			// Using If for optional content
			g.Hr(),
			g.If(isLoggedIn, // Try to toggle this, and see the result
				g.Div(g.KV{"class": "user-actions"}).Add(
					g.Button().Add(g.Text("Profile")),
					g.Text(" "), // Add a whitespace between the two buttons. not needed if using css styles
					g.Button().Add(g.Text("Settings")),
				),
			),

			// Using Repeat to generate repeated elements
			g.Hr(),
			g.H2().Add(g.Text("Repeated Elements")),
			g.Repeat(3, func() g.Node {
				return g.Div(g.KV{"class": "repeated-item"}).Add(
					g.Text("This is a repeated item"),
					g.Br(),
				)
			}),

			// Using Map to transform data into elements
			g.Hr(),
			g.H2().Add(g.Text("Mapped List")),
			g.Ul().Add(
				g.Map(items, func(item string) g.Node {
					if item == "Apple" {
						return g.Li().Add(g.Text(item), g.Span(g.KV{"class": "badge"}).Add(g.Text(" (Popular)")))
					}
					return g.Li().Add(g.Text(item))
				}),
			),

			// Combining utilities
			g.Hr(),
			g.H2().Add(g.Text("Combined Example")),
			g.Div().Add(
				g.Text("Total items: "), g.Strong().Add(g.Text(fmt.Sprint(len(items)))),
				g.If(len(items) > 2,
					g.P().Add(g.Text("There are many items to display!")),
				),
			),
		),
	)

	if err := g.Render(os.Stdout, page); err != nil {
		panic(err)
	}
	//
	// the same as:
	//
	// html, err := page.Render()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(html)
}
