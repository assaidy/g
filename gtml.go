package gtml

import (
	"fmt"
	"reflect"
	"strings"
)

type KV map[string]any

type Element struct {
	Tag           string
	IsSelfClosing bool
	Attrs         KV
	Children      []*Element
	IsText        bool
	Text          string // only used if IsText is true
}

func (me *Element) validateAttributes() error {
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

func (me *Element) Render() (string, error) {
	if err := me.validateAttributes(); err != nil {
		return "", err
	}
	return me.renderHTML(), nil
}

func (me *Element) renderHTML() string {
	if me.IsText {
		return me.Text
	}

	var html strings.Builder

	html.WriteString("<")
	html.WriteString(me.Tag)

	for key, value := range me.Attrs {
		if value == true {
			html.WriteString(fmt.Sprintf(" %s", key))
		} else if value != "" {
			html.WriteString(fmt.Sprintf(` %s="%s"`, key, value))
		}
	}

	html.WriteString(">")
	if me.IsSelfClosing {
		return html.String()
	}

	if len(me.Children) > 0 {
		for _, child := range me.Children {
			html.WriteString(child.renderHTML())
		}
	}

	html.WriteString(fmt.Sprintf("</%s>", me.Tag))
	return html.String()
}

func (me *Element) Add(children ...*Element) *Element {
	if me.IsSelfClosing {
		panic("trying to add children to a self-closing element")
	}
	me.Children = append(me.Children, children...)
	return me
}

func newElement(tag string, attrs []KV, isSelfClosing ...bool) *Element {
	e := &Element{Tag: tag}
	if len(attrs) != 0 {
		e.Attrs = attrs[0]
	}
	if len(isSelfClosing) != 0 {
		e.IsSelfClosing = isSelfClosing[0]
	}
	return e
}

func Text(s string) *Element {
	return &Element{
		IsText: true,
		Text:   s,
	}
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
	return newElement("hq", attrs)
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
