package g

import (
	"fmt"
	"html"
	"slices"
	"strings"
	"unicode"
)

// Text represents a plain text node that renders HTML-escaped content.
// Unlike HTML elements, Text nodes are not wrapped in tags and are rendered
// as literal text content with HTML entities automatically escaped.
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
//
// The value type must be either string or bool:
//   - string: Attribute will have the format key="value" (HTML-escaped)
//   - bool: If true, attribute appears as key (valueless). If false, attribute is omitted.
//   - any other type triggers an error during rendering.
//
// Example:
//
//	KV{"class": "container", "hidden": true, "disabled": false}
//	// Renders: class="container" hidden
type KV map[string]any

// Element represents an HTML element with its attributes and children.
type Element struct {
	Tag      string // HTML tag name
	IsVoid   bool   // Whether the tag is self-closing (e.g., <br>, <img>)
	Attrs    KV     // HTML attributes as key-value pairs
	Children []Node // Child nodes
}

// Render generates the HTML string for the element and its children.
//
// Behavior:
//   - Self-closing tags: Render without children (e.g., <br />, <img />)
//   - Regular tags: Render with opening tag, children, and closing tag
//   - Empty elements: Render children
//   - Attributes: Properly HTML-escaped and formatted
//
// Returns the complete HTML string as byteslice and any error encountered.
func (me *Element) Render() (string, error) {
	builder := &strings.Builder{}

	if me.Tag == "" { // empty tag
		if err := me.renderChildren(builder); err != nil {
			return "", err
		}
		return builder.String(), nil
	}

	fmt.Fprint(builder, "<")
	fmt.Fprint(builder, me.Tag)
	if err := me.renderAttrs(builder); err != nil {
		return "", err
	}
	fmt.Fprint(builder, ">")

	if me.IsVoid {
		return builder.String(), nil
	}

	if err := me.renderChildren(builder); err != nil {
		return "", nil
	}
	fmt.Fprintf(builder, "</%s>", me.Tag)

	return builder.String(), nil
}

func (me Element) renderAttrs(builder *strings.Builder) error {
	// for deterministic attrs order
	type kv struct {key string; value any }
	attrSlice := make([]kv, 0, len(me.Attrs))
	for key, value := range me.Attrs {
		attrSlice = append(attrSlice, kv{key, value})
	}
	slices.SortFunc(attrSlice, func(a, b kv) int {
		return strings.Compare(a.key, b.key)
	})

	for _, attr := range attrSlice {
		k := strings.TrimSpace(attr.key)
		if k == "" {
			return fmt.Errorf("empty/whitespace attribute key not allowed.")
		}
		if attr.value == nil {
			return fmt.Errorf("attribute '%s' has nil value", k)
		}

		switch v := attr.value.(type) {
		case string:
			fmt.Fprintf(builder, ` %s="%s"`, k, html.EscapeString(v))
		case bool:
			if v == true {
				fmt.Fprintf(builder, " %s", k)
			}
		default:
			return fmt.Errorf("attribute value must be string or bool, got %T for key '%s'", v, k)
		}
	}

	return nil
}

func (me Element) renderChildren(builder *strings.Builder) error {
	for _, child := range me.Children {
		s, err := child.Render()
		if err != nil {
			return err
		}
		fmt.Fprint(builder, s)
	}
	return nil
}

// Add appends child elements to this element and returns the element for method chaining.
//
// For void elements (self-closing tags like <br>, <img>, <meta>), this method
// is a no-op since void elements cannot have children according to HTML specifications.
//
// Example:
//
//	div := Div().Add(
//	    P().Add(Text("First paragraph")),
//	    Span().Add(Text("Important text")),
//	)
//
// The method returns the element itself to enable fluent chaining.
func (me *Element) Add(children ...Node) Node {
	if !me.IsVoid {
		me.Children = append(me.Children, children...)
	}
	return me
}

func newElement(tag string, attrs []KV, isVoid ...bool) *Element {
	e := &Element{Tag: tag, Attrs: make(KV)}
	if len(attrs) != 0 {
		e.Attrs = attrs[0]
	}
	if len(isVoid) != 0 {
		e.IsVoid = isVoid[0]
	}
	return e
}

// Empty creates an empty element (no tag).
func Empty(attrs ...KV) *Element {
	return newElement("", attrs)
}

// Html creates the root element of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func Html(attrs ...KV) *Element {
	return newElement("html", attrs)
}

// Head contains machine-readable information about the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func Head(attrs ...KV) *Element {
	return newElement("head", attrs)
}

// Title defines the document's title that is shown in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func Title(attrs ...KV) *Element {
	return newElement("title", attrs)
}

// Link specifies relationships between the current document and an external resource.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func Link(attrs ...KV) *Element {
	return newElement("link", attrs, true)
}

// Meta represents metadata that cannot be represented by other HTML meta-related elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func Meta(attrs ...KV) *Element {
	return newElement("meta", attrs, true)
}

// Style contains style information for a document or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func Style(attrs ...KV) *Element {
	return newElement("style", attrs)
}

// Body represents the content of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func Body(attrs ...KV) *Element {
	return newElement("body", attrs)
}

// H1 creates a level 1 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(attrs ...KV) *Element {
	return newElement("h1", attrs)
}

// H2 creates a level 2 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(attrs ...KV) *Element {
	return newElement("h2", attrs)
}

// H3 creates a level 3 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(attrs ...KV) *Element {
	return newElement("h3", attrs)
}

// H4 creates a level 4 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(attrs ...KV) *Element {
	return newElement("h4", attrs)
}

// H5 creates a level 5 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(attrs ...KV) *Element {
	return newElement("h5", attrs)
}

// H6 creates a level 6 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(attrs ...KV) *Element {
	return newElement("h6", attrs)
}

// Header creates a header element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func Header(attrs ...KV) *Element {
	return newElement("header", attrs)
}

// Footer creates a footer element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func Footer(attrs ...KV) *Element {
	return newElement("footer", attrs)
}

// Nav creates a navigation element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
func Nav(attrs ...KV) *Element {
	return newElement("nav", attrs)
}

// Main creates a main content element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
func Main(attrs ...KV) *Element {
	return newElement("main", attrs)
}

// Section creates a section element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
func Section(attrs ...KV) *Element {
	return newElement("section", attrs)
}

// Article creates an article element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func Article(attrs ...KV) *Element {
	return newElement("article", attrs)
}

// Aside creates an aside element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func Aside(attrs ...KV) *Element {
	return newElement("aside", attrs)
}

// Hr represents a thematic break between paragraph-level elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
func Hr() *Element {
	return newElement("hr", nil, true)
}

// Pre represents preformatted text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
func Pre(attrs ...KV) *Element {
	return newElement("pre", attrs)
}

// Blockquote represents a section quoted from another source.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
func Blockquote(attrs ...KV) *Element {
	return newElement("blockquote", attrs)
}

// Ol represents an ordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
func Ol(attrs ...KV) *Element {
	return newElement("ol", attrs)
}

// Ul represents an unordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
func Ul(attrs ...KV) *Element {
	return newElement("ul", attrs)
}

// Li represents a list item.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
func Li(attrs ...KV) *Element {
	return newElement("li", attrs)
}

// A creates hyperlinks to other web pages, files, locations within the same page, or anything else a URL can address.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
func A(attrs ...KV) *Element {
	return newElement("a", attrs)
}

// Em marks text with emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
func Em(attrs ...KV) *Element {
	return newElement("em", attrs)
}

// Strong indicates strong importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
func Strong(attrs ...KV) *Element {
	return newElement("strong", attrs)
}

// Code displays its contents styled as computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
func Code(attrs ...KV) *Element {
	return newElement("code", attrs)
}

// Var represents a variable in a mathematical expression or programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
func Var(attrs ...KV) *Element {
	return newElement("var", attrs)
}

// Samp represents sample output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
func Samp(attrs ...KV) *Element {
	return newElement("samp", attrs)
}

// Kbd represents text that the user should enter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
func Kbd(attrs ...KV) *Element {
	return newElement("kbd", attrs)
}

// Sub specifies inline text displayed as subscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
func Sub(attrs ...KV) *Element {
	return newElement("sub", attrs)
}

// Sup specifies inline text displayed as superscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
func Sup(attrs ...KV) *Element {
	return newElement("sup", attrs)
}

// I represents text in an alternate voice or mood.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
func I(attrs ...KV) *Element {
	return newElement("i", attrs)
}

// B draws attention to text without conveying importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
func B(attrs ...KV) *Element {
	return newElement("b", attrs)
}

// U represents text with an unarticulated annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
func U(attrs ...KV) *Element {
	return newElement("u", attrs)
}

// Mark highlights text for reference.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
func Mark(attrs ...KV) *Element {
	return newElement("mark", attrs)
}

// Bdi isolates text for bidirectional text formatting.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
func Bdi(attrs ...KV) *Element {
	return newElement("bdi", attrs)
}

// Bdo overrides the current text direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
func Bdo(attrs ...KV) *Element {
	return newElement("bdo", attrs)
}

// Br produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
func Br() *Element {
	return newElement("br", nil, true)
}

// Wbr represents a word break opportunity.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
func Wbr(attrs ...KV) *Element {
	return newElement("wbr", attrs, true)
}

// Img embeds an image into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
func Img(attrs ...KV) *Element {
	return newElement("img", attrs, true)
}

// Iframe represents a nested browsing context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
func Iframe(attrs ...KV) *Element {
	return newElement("iframe", attrs)
}

// Embed embeds external content at the specified point in the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
func Embed(attrs ...KV) *Element {
	return newElement("embed", attrs, true)
}

// Object represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
func Object(attrs ...KV) *Element {
	return newElement("object", attrs)
}

// Picture defines multiple sources for an img element to offer alternative versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
func Picture(attrs ...KV) *Element {
	return newElement("picture", attrs)
}

// Source specifies multiple media resources for the picture, the audio element, or the video element. It is a void element, meaning that it has no content and does not have a closing tag. It is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
func Source(attrs ...KV) *Element {
	return newElement("source", attrs, true)
}

// Track is used as a child of the media elements, audio and video. It lets you specify timed text tracks (or time-based data), for example to automatically handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web Video Text Tracks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
func Track(attrs ...KV) *Element {
	return newElement("track", attrs, true)
}

// Video embeds a media player which supports video playback into the document. You can also use &lt;video&gt; for audio content, but the audio element may provide a more appropriate user experience.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
func Video(attrs ...KV) *Element {
	return newElement("video", attrs)
}

// Audio is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the source element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
func Audio(attrs ...KV) *Element {
	return newElement("audio", attrs)
}

// Canvas is a container element to use with either the canvas scripting API or the WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
func Canvas(attrs ...KV) *Element {
	return newElement("canvas", attrs)
}

// MapElement is used with &lt;area&gt; elements to define an image map (a clickable link area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
func MapElement(attrs ...KV) *Element {
	return newElement("map", attrs)
}

// Area defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hyperlink.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
func Area(attrs ...KV) *Element {
	return newElement("area", attrs, true)
}

// Svg is a container defining a new coordinate system and viewport. It is used as the outermost element of SVG documents, but it can also be used to embed an SVG fragment inside an SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/svg
func Svg(attrs ...KV) *Element {
	return newElement("svg", attrs)
}

// Math is the top-level element in MathML. Every valid MathML instance must be wrapped in it. In addition, you must not nest a second &lt;math&gt; element in another, but you can have an arbitrary number of other child elements in it.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/math
func Math(attrs ...KV) *Element {
	return newElement("math", attrs)
}

// Script is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The &lt;script&gt; element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
func Script(attrs ...KV) *Element {
	return newElement("script", attrs)
}

// Noscript defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
func Noscript(attrs ...KV) *Element {
	return newElement("noscript", attrs)
}

// Del represents a range of text that has been deleted from a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
func Del(attrs ...KV) *Element {
	return newElement("del", attrs)
}

// Ins represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
func Ins(attrs ...KV) *Element {
	return newElement("ins", attrs)
}

// Table represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
func Table(attrs ...KV) *Element {
	return newElement("table", attrs)
}

// Caption specifies the caption (or title) of a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
func Caption(attrs ...KV) *Element {
	return newElement("caption", attrs)
}

// Colgroup defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
func Colgroup(attrs ...KV) *Element {
	return newElement("colgroup", attrs)
}

// Col defines one or more columns in a column group represented by its implicit or explicit parent &lt;colgroup&gt; element. The &lt;col&gt; element is only valid as a child of a &lt;colgroup&gt; element that has no span attribute defined.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
func Col(attrs ...KV) *Element {
	return newElement("col", attrs, true)
}

// Thead groups the header content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
func Thead(attrs ...KV) *Element {
	return newElement("thead", attrs)
}

// Tbody groups the body content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
func Tbody(attrs ...KV) *Element {
	return newElement("tbody", attrs)
}

// Tfoot groups the footer content in a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
func Tfoot(attrs ...KV) *Element {
	return newElement("tfoot", attrs)
}

// Tr defines a row of cells in a table. The row's cells can then be established using a mix of &lt;td&gt; (data cell) and &lt;th&gt; (header cell) elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
func Tr(attrs ...KV) *Element {
	return newElement("tr", attrs)
}

// Th is a child of the &lt;tr&gt; element, it defines a cell as the header of a group of table cells. The nature of this group can be explicitly defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
func Th(attrs ...KV) *Element {
	return newElement("th", attrs)
}

// Td is a child of the &lt;tr&gt; element, it defines a cell of a table that contains data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
func Td(attrs ...KV) *Element {
	return newElement("td", attrs)
}

// Form represents a document section containing interactive controls for submitting information.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
func Form(attrs ...KV) *Element {
	return newElement("form", attrs)
}

// Fieldset is used to group several controls as well as labels (&lt;label&gt;) within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
func Fieldset(attrs ...KV) *Element {
	return newElement("fieldset", attrs)
}

// Legend represents a caption for the content of its parent &lt;fieldset&gt;.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
func Legend(attrs ...KV) *Element {
	return newElement("legend", attrs)
}

// Label represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
func Label(attrs ...KV) *Element {
	return newElement("label", attrs)
}

// Input is used to create interactive controls for web-based forms to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The &lt;input&gt; element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
func Input(attrs ...KV) *Element {
	return newElement("input", attrs, true)
}

// Button is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it performs an action, such as submitting a form or opening a dialog.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
func Button(attrs ...KV) *Element {
	return newElement("button", attrs)
}

// Select represents a control that provides a menu of options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
func Select(attrs ...KV) *Element {
	return newElement("select", attrs)
}

// Datalist contains a set of &lt;option&gt; elements that represent the permissible or recommended options available to choose from within other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
func Datalist(attrs ...KV) *Element {
	return newElement("datalist", attrs)
}

// Optgroup creates a grouping of options within a &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
func Optgroup(attrs ...KV) *Element {
	return newElement("optgroup", attrs)
}

// Option is used to define an item contained in a &lt;select&gt;, an &lt;optgroup&gt;, or a &lt;datalist&gt; element. As such, &lt;option&gt; can represent menu items in popups and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
func Option(attrs ...KV) *Element {
	return newElement("option", attrs)
}

// Textarea represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example, a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
func Textarea(attrs ...KV) *Element {
	return newElement("textarea", attrs)
}

// Output is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
func Output(attrs ...KV) *Element {
	return newElement("output", attrs)
}

// Progress displays an indicator showing the completion progress of a task, typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
func Progress(attrs ...KV) *Element {
	return newElement("progress", attrs)
}

// Meter represents either a scalar value within a known range or a fractional value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
func Meter(attrs ...KV) *Element {
	return newElement("meter", attrs)
}

// Details creates a disclosure widget in which information is visible only when the widget is toggled into an "open" state. A summary or label must be provided using the &lt;summary&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
func Details(attrs ...KV) *Element {
	return newElement("details", attrs)
}

// Summary specifies a summary, caption, or legend for a details element's disclosure box. Clicking the &lt;summary&gt; element toggles the state of the parent &lt;details&gt; element open and closed.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
func Summary(attrs ...KV) *Element {
	return newElement("summary", attrs)
}

// Dialog represents a dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
func Dialog(attrs ...KV) *Element {
	return newElement("dialog", attrs)
}

// Slot acts as a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
func Slot(attrs ...KV) *Element {
	return newElement("slot", attrs)
}

// Template holds HTML that is not to be rendered immediately when a page is loaded but may be instantiated subsequently during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
func Template(attrs ...KV) *Element {
	return newElement("template", attrs)
}

// Fencedframe represents a nested browsing context, like &lt;iframe&gt; but with more native privacy features built in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fencedframe
func Fencedframe(attrs ...KV) *Element {
	return newElement("fencedframe", attrs)
}

// Selectedcontent displays the content of the currently selected &lt;option&gt; inside a closed &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/selectedcontent
func Selectedcontent(attrs ...KV) *Element {
	return newElement("selectedcontent", attrs)
}

// Base specifies the base URL and default browsing context for relative URLs.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func Base(attrs ...KV) *Element {
	return newElement("base", attrs, true)
}

// Hgroup groups a set of h1–h6 elements when they represent a multi-level heading.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
func Hgroup(attrs ...KV) *Element {
	return newElement("hgroup", attrs)
}

// Address indicates contact information for a person or organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func Address(attrs ...KV) *Element {
	return newElement("address", attrs)
}

// Search represents a search or filtering interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
func Search(attrs ...KV) *Element {
	return newElement("search", attrs)
}

// Div is the generic container for flow content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
func Div(attrs ...KV) *Element {
	return newElement("div", attrs)
}

// Span is the generic inline container for phrasing content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
func Span(attrs ...KV) *Element {
	return newElement("span", attrs)
}

// P creates a paragraph element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
func P(attrs ...KV) *Element {
	return newElement("p", attrs)
}

// Dl represents a description list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
func Dl(attrs ...KV) *Element {
	return newElement("dl", attrs)
}

// Dt specifies a term in a description or definition list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
func Dt(attrs ...KV) *Element {
	return newElement("dt", attrs)
}

// Dd provides the description, definition, or value for the preceding term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
func Dd(attrs ...KV) *Element {
	return newElement("dd", attrs)
}

// Figure represents self-contained content with an optional caption.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
func Figure(attrs ...KV) *Element {
	return newElement("figure", attrs)
}

// Figcaption represents a caption or legend for the contents of its parent figure element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
func Figcaption(attrs ...KV) *Element {
	return newElement("figcaption", attrs)
}

// Menu represents a set of commands or options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
func Menu(attrs ...KV) *Element {
	return newElement("menu", attrs)
}

// Small represents side-comments and small print.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
func Small(attrs ...KV) *Element {
	return newElement("small", attrs)
}

// S renders text with a strikethrough.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
func S(attrs ...KV) *Element {
	return newElement("s", attrs)
}

// Cite marks the title of a creative work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
func Cite(attrs ...KV) *Element {
	return newElement("cite", attrs)
}

// Q indicates a short inline quotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
func Q(attrs ...KV) *Element {
	return newElement("q", attrs)
}

// Dfn indicates the defining instance of a term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
func Dfn(attrs ...KV) *Element {
	return newElement("dfn", attrs)
}

// Abbr represents an abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
func Abbr(attrs ...KV) *Element {
	return newElement("abbr", attrs)
}

// Ruby represents ruby annotations for East Asian typography.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
func Ruby(attrs ...KV) *Element {
	return newElement("ruby", attrs)
}

// Rt specifies the ruby text for ruby annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
func Rt(attrs ...KV) *Element {
	return newElement("rt", attrs)
}

// Rp provides parentheses for browsers that don't support ruby text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
func Rp(attrs ...KV) *Element {
	return newElement("rp", attrs)
}

// Data links content with a machine-readable translation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
func Data(attrs ...KV) *Element {
	return newElement("data", attrs)
}

// Time represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
func Time(attrs ...KV) *Element {
	return newElement("time", attrs)
}
