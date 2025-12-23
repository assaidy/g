package g

import (
	"fmt"
	"html"
	"io"
	"reflect"
	"strings"
	"unicode"
)

func Render(writer io.Writer, node Node) error {
	s, err := node.Render()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(s))
	return err
}

type Node interface {
	Render() (string, error)
}

// Text creates a plain text node element (not a tag)
type Text string

func (me Text) Render() (string, error) {
	s := string(me)
	if s == "" {
		return "", nil
	}

	startsWithSpace := unicode.IsSpace(rune(s[0]))
	endsWithSpace := len(s) > 1 && unicode.IsSpace(rune(s[len(s)-1]))

	s = strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")

	// re-apply at most one space at each edge
	if startsWithSpace {
		s = " " + s
	}
	if endsWithSpace {
		s = s + " "
	}

	return html.EscapeString(s), nil
}

// KV represents a key-value map for HTML attributes.
// Value can only be string or bool, otherwirse Render will return an error.
// If value is bool, attribute doesn't have a value, and will be included if value is true.
type KV map[string]any

// Element represents an HTML element or text node
type Element struct {
	Tag           string // HTML tag name
	IsSelfClosing bool   // Whether the tag is self-closing (e.g., <br>, <img>)
	Attrs         KV     // HTML attributes as key-value pairs
	Children      []Node // Child nodes
}

func (me Element) validateAttributes() error {
	for key, value := range me.Attrs {
		if key == "" {
			return fmt.Errorf("empty attribute key not allowed")
		}
		if value == nil {
			return fmt.Errorf("attribute '%s' has nil value", key)
		}
		valType := reflect.TypeOf(value)
		if valType.Kind() != reflect.String && valType.Kind() != reflect.Bool {
			return fmt.Errorf("attribute value must be string or bool, got %T for key '%s'", value, key)
		}
	}
	return nil
}

// Render generates the HTML string for the element and its children.
// If element is a self-closing tag, children will be ignored.
// If element is nil, it will be rendered as empty string.
func (me *Element) Render() (string, error) {
	if err := me.validateAttributes(); err != nil {
		return "", err
	}
	return me.renderHTML(), nil
}

func (me Element) renderHTML() string {
	var builder strings.Builder

	// not Empty tag
	if me.Tag != "" {
		builder.WriteString("<")
		builder.WriteString(me.Tag)

		for key, value := range me.Attrs {
			if value == true {
				builder.WriteString(fmt.Sprintf(" %s", key))
			} else if value != "" {
				builder.WriteString(fmt.Sprintf(` %s="%s"`, key, value))
			}
		}

		builder.WriteString(">")
		if me.IsSelfClosing {
			return builder.String()
		}
	}

	if len(me.Children) > 0 {
		for _, child := range me.Children {
			s, err := child.Render()
			if err != nil {
				return ""
			}
			builder.WriteString(s)
		}
	}

	if me.Tag != "" {
		builder.WriteString(fmt.Sprintf("</%s>", me.Tag))
	}

	return builder.String()
}

// Add appends child elements to this element and returns the element for chaining
func (me *Element) Add(children ...Node) Node {
	// NOTE: children are not rendered if IsSelfClosing
	// TODO: create SelfClosingElement{} that doesn't have Add() method
	me.Children = append(me.Children, children...)
	return me
}

func newElement(tag string, attrs []KV, isSelfClosing ...bool) *Element {
	e := &Element{Tag: tag, Attrs: make(KV)}
	if len(attrs) != 0 {
		e.Attrs = attrs[0]
	}
	if len(isSelfClosing) != 0 {
		e.IsSelfClosing = isSelfClosing[0]
	}
	return e
}

func Empty(attrs ...KV) *Element {
	return newElement("", attrs)
}

func Ul(attrs ...KV) *Element {
	return newElement("ul", attrs)
}

func Ol(attrs ...KV) *Element {
	return newElement("ol", attrs)
}

func Li(attrs ...KV) *Element {
	return newElement("li", attrs)
}

func Div(attrs ...KV) *Element {
	return newElement("div", attrs)
}

func Span(attrs ...KV) *Element {
	return newElement("span", attrs)
}

func P(attrs ...KV) *Element {
	return newElement("p", attrs)
}

func H1(attrs ...KV) *Element {
	return newElement("h1", attrs)
}

func H2(attrs ...KV) *Element {
	return newElement("h2", attrs)
}

func H3(attrs ...KV) *Element {
	return newElement("h3", attrs)
}

func H4(attrs ...KV) *Element {
	return newElement("h4", attrs)
}

func H5(attrs ...KV) *Element {
	return newElement("h5", attrs)
}

func H6(attrs ...KV) *Element {
	return newElement("h6", attrs)
}

func A(attrs ...KV) *Element {
	return newElement("a", attrs)
}

func Img(attrs ...KV) *Element {
	return newElement("img", attrs, true)
}

func Br() *Element {
	return newElement("br", nil, true)
}

func Hr() *Element {
	return newElement("hr", nil, true)
}

func Input(attrs ...KV) *Element {
	return newElement("input", attrs, true)
}

func Button(attrs ...KV) *Element {
	return newElement("button", attrs)
}

func Form(attrs ...KV) *Element {
	return newElement("form", attrs)
}

func Label(attrs ...KV) *Element {
	return newElement("label", attrs)
}

func Textarea(attrs ...KV) *Element {
	return newElement("textarea", attrs)
}

func Select(attrs ...KV) *Element {
	return newElement("select", attrs)
}

func Option(attrs ...KV) *Element {
	return newElement("option", attrs)
}

func Table(attrs ...KV) *Element {
	return newElement("table", attrs)
}

func Thead(attrs ...KV) *Element {
	return newElement("thead", attrs)
}

func Tbody(attrs ...KV) *Element {
	return newElement("tbody", attrs)
}

func Tr(attrs ...KV) *Element {
	return newElement("tr", attrs)
}

func Th(attrs ...KV) *Element {
	return newElement("th", attrs)
}

func Td(attrs ...KV) *Element {
	return newElement("td", attrs)
}

func Header(attrs ...KV) *Element {
	return newElement("header", attrs)
}

func Footer(attrs ...KV) *Element {
	return newElement("footer", attrs)
}

func Nav(attrs ...KV) *Element {
	return newElement("nav", attrs)
}

func Main(attrs ...KV) *Element {
	return newElement("main", attrs)
}

func Section(attrs ...KV) *Element {
	return newElement("section", attrs)
}

func Article(attrs ...KV) *Element {
	return newElement("article", attrs)
}

func Aside(attrs ...KV) *Element {
	return newElement("aside", attrs)
}

func Code(attrs ...KV) *Element {
	return newElement("code", attrs)
}

func Pre(attrs ...KV) *Element {
	return newElement("pre", attrs)
}

func Blockquote(attrs ...KV) *Element {
	return newElement("blockquote", attrs)
}

func Em(attrs ...KV) *Element {
	return newElement("em", attrs)
}

func Strong(attrs ...KV) *Element {
	return newElement("strong", attrs)
}

func Meta(attrs ...KV) *Element {
	return newElement("meta", attrs, true)
}

func Link(attrs ...KV) *Element {
	return newElement("link", attrs, true)
}

func Script(attrs ...KV) *Element {
	return newElement("script", attrs)
}

func Style(attrs ...KV) *Element {
	return newElement("style", attrs)
}

func Html(attrs ...KV) *Element {
	return newElement("html", attrs)
}

func Head(attrs ...KV) *Element {
	return newElement("head", attrs)
}

func Body(attrs ...KV) *Element {
	return newElement("body", attrs)
}

func Title(attrs ...KV) *Element {
	return newElement("title", attrs)
}
