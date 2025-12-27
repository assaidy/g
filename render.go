package g

import "io"

// Render writes the HTML representation of a Node to the provided io.Writer.
//
// This is a convenience function that combines Node.Render() with writing
// the output to an io.Writer, making it suitable for writing directly to
// files, HTTP responses, or other output streams.
//
// Example:
//
//	err := Render(os.Stdout, Div(Text("Hello")))
//	// Outputs: <div>Hello</div>
func Render(writer io.Writer, node Node) error {
	s, err := node.Render()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(s))
	return err
}

// Node represents any renderable HTML element or text content.
//
// The Node interface is the core abstraction that allows both HTML elements
// and text content to be treated uniformly when building and rendering HTML
// trees. All elements created by the factory functions (Div(), P(), Text(), etc.)
// implement this interface.
//
// Example:
//
//	var node Node = Div(Text("Hello"))
//	html, err := node.Render()
type Node interface {
	Render() (string, error)
}
