package main

import (
	"fmt"
	"os"

	"github.com/assaidy/g"
	gu "github.com/assaidy/g/utils"
)

func main() {
	// Sample data for demonstration
	items := []string{"Apple", "Banana", "Cherry", "Date"}
	isLoggedIn := true
	userName := "John Doe"

	// Demonstrate all utility functions in a single example
	page := g.Html(g.KV{"lang": "en"},
		g.Head(
			g.Meta(g.KV{"charset": "UTF-8"}),
			g.Title(g.Text("GTML Utils Example")),
		),
		g.Body(
			// Using IfElse to show conditional content
			gu.IfElse(isLoggedIn,
				g.Div(g.KV{"class": "welcome"},
					g.H1(g.Text(fmt.Sprintf("Welcome back, %s", userName))),
					g.P(g.Text("You are logged in!")),
				),
				g.Div(g.KV{"class": "login-prompt"},
					g.H1(g.Text("Please log in")),
					g.P(g.Text("You need to authenticate to continue.")),
				),
			),

			// Using If for optional content
			g.Hr(),
			gu.If(isLoggedIn, // Try to toggle this, and see the result
				g.Div(g.KV{"class": "user-actions"},
					g.Button(g.Text("Profile")),
					g.Text(" "), // Add a whitespace between the two buttons. not needed if using css styles
					g.Button(g.Text("Settings")),
				),
			),

			// Using Repeat to generate repeated elements
			g.Hr(),
			g.H2(g.Text("Repeated Elements")),
			gu.Repeat(3, func() g.Node {
				return g.Div(g.KV{"class": "repeated-item"},
					g.Text("This is a repeated item"),
					g.Br(),
				)
			}),

			// Using Map to transform data into elements
			g.Hr(),
			g.H2(g.Text("Mapped List")),
			g.Ul(
				gu.Map(items, func(item string) g.Node {
					if item == "Apple" {
						return g.Li(g.Text(item), g.Span(g.KV{"class": "badge"}, g.Text(" (Popular)")))
					}
					return g.Li(g.Text(item))
				}),
			),

			// Combining utilities
			g.Hr(),
			g.H2(g.Text("Combined Example")),
			g.Div(
				g.Text("Total items: "), g.Strong(g.Text(fmt.Sprint(len(items)))),
				gu.If(len(items) > 2,
					g.P(g.Text("There are many items to display!")),
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
