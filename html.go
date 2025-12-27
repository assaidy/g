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
	type kv struct {
		key   string
		value any
	}
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

func newElem(tag string, args ...any) *Element {
	e := &Element{Tag: tag}
	for _, arg := range args {
		switch value := arg.(type) {
		case KV:
			e.Attrs = value // this overrides the old ones if duplicates are passed
		case Node:
			e.Children = append(e.Children, value)
		default: // ignore
		}
	}
	return e
}

func newVoidElem(tag string, attrs ...KV) *Element {
	if len(attrs) > 0 {
		return &Element{Tag: tag, Attrs: attrs[0]}
	}
	return &Element{Tag: tag}
}

// Empty creates an empty element (no tag).
func Empty(args ...any) *Element {
	return newElem("", args...)
}

// Html creates the root element of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func Html(args ...any) *Element {
	return newElem("html", args...)
}

// Head contains machine-readable information about the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func Head(args ...any) *Element {
	return newElem("head", args...)
}

// Title defines the document's title that is shown in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func Title(args ...any) *Element {
	return newElem("title", args...)
}

// Link specifies relationships between the current document and an external resource.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func Link(attrs ...KV) *Element {
	return newVoidElem("link", attrs...)
}

// Meta represents metadata that cannot be represented by other HTML meta-related elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func Meta(attrs ...KV) *Element {
	return newVoidElem("meta", attrs...)
}

// Style contains style information for a document or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func Style(args ...any) *Element {
	return newElem("style", args...)
}

// Body represents the content of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func Body(args ...any) *Element {
	return newElem("body", args...)
}

// H1 creates a level 1 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(args ...any) *Element {
	return newElem("h1", args...)
}

// H2 creates a level 2 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(args ...any) *Element {
	return newElem("h2", args...)
}

// H3 creates a level 3 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(args ...any) *Element {
	return newElem("h3", args...)
}

// H4 creates a level 4 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(args ...any) *Element {
	return newElem("h4", args...)
}

// H5 creates a level 5 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(args ...any) *Element {
	return newElem("h5", args...)
}

// H6 creates a level 6 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(args ...any) *Element {
	return newElem("h6", args...)
}

// Header creates a header element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func Header(args ...any) *Element {
	return newElem("header", args...)
}

// Footer creates a footer element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func Footer(args ...any) *Element {
	return newElem("footer", args...)
}

// Nav creates a navigation element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
func Nav(args ...any) *Element {
	return newElem("nav", args...)
}

// Main creates a main content element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
func Main(args ...any) *Element {
	return newElem("main", args...)
}

// Section creates a section element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
func Section(args ...any) *Element {
	return newElem("section", args...)
}

// Article creates an article element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func Article(args ...any) *Element {
	return newElem("article", args...)
}

// Aside creates an aside element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func Aside(args ...any) *Element {
	return newElem("aside", args...)
}

// Hr represents a thematic break between paragraph-level elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
func Hr(attrs ...KV) *Element {
	return newVoidElem("hr", attrs...)
}

// Pre represents preformatted text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
func Pre(args ...any) *Element {
	return newElem("pre", args...)
}

// Blockquote represents a section quoted from another source.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
func Blockquote(args ...any) *Element {
	return newElem("blockquote", args...)
}

// Ol represents an ordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
func Ol(args ...any) *Element {
	return newElem("ol", args...)
}

// Ul represents an unordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
func Ul(args ...any) *Element {
	return newElem("ul", args...)
}

// Li represents a list item.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
func Li(args ...any) *Element {
	return newElem("li", args...)
}

// A creates hyperlinks to other web pages, files, locations within the same page, or anything else a URL can address.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
func A(args ...any) *Element {
	return newElem("a", args...)
}

// Em marks text with emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
func Em(args ...any) *Element {
	return newElem("em", args...)
}

// Strong indicates strong importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
func Strong(args ...any) *Element {
	return newElem("strong", args...)
}

// Code displays its contents styled as computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
func Code(args ...any) *Element {
	return newElem("code", args...)
}

// Var represents a variable in a mathematical expression or programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
func Var(args ...any) *Element {
	return newElem("var", args...)
}

// Samp represents sample output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
func Samp(args ...any) *Element {
	return newElem("samp", args...)
}

// Kbd represents text that the user should enter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
func Kbd(args ...any) *Element {
	return newElem("kbd", args...)
}

// Sub specifies inline text displayed as subscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
func Sub(args ...any) *Element {
	return newElem("sub", args...)
}

// Sup specifies inline text displayed as superscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
func Sup(args ...any) *Element {
	return newElem("sup", args...)
}

// I represents text in an alternate voice or mood.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
func I(args ...any) *Element {
	return newElem("i", args...)
}

// B draws attention to text without conveying importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
func B(args ...any) *Element {
	return newElem("b", args...)
}

// U represents text with an unarticulated annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
func U(args ...any) *Element {
	return newElem("u", args...)
}

// Mark highlights text for reference.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
func Mark(args ...any) *Element {
	return newElem("mark", args...)
}

// Bdi isolates text for bidirectional text formatting.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
func Bdi(args ...any) *Element {
	return newElem("bdi", args...)
}

// Bdo overrides the current text direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
func Bdo(args ...any) *Element {
	return newElem("bdo", args...)
}

// Br produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
func Br(attrs ...KV) *Element {
	return newVoidElem("br", attrs...)
}

// Wbr represents a word break opportunity.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
func Wbr(attrs ...KV) *Element {
	return newVoidElem("wbr", attrs...)
}

// Img embeds an image into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
func Img(attrs ...KV) *Element {
	return newVoidElem("img", attrs...)
}

// Iframe represents a nested browsing context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
func Iframe(args ...any) *Element {
	return newElem("iframe", args...)
}

// Embed embeds external content at the specified point in the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
func Embed(attrs ...KV) *Element {
	return newVoidElem("embed", attrs...)
}

// Object represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
func Object(args ...any) *Element {
	return newElem("object", args...)
}

// Picture defines multiple sources for an img element to offer alternative versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
func Picture(args ...any) *Element {
	return newElem("picture", args...)
}

// Source specifies multiple media resources for the picture, the audio element, or the video element. It is a void element, meaning that it has no content and does not have a closing tag. It is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
func Source(attrs ...KV) *Element {
	return newVoidElem("source", attrs...)
}

// Track is used as a child of the media elements, audio and video. It lets you specify timed text tracks (or time-based data), for example to automatically handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web Video Text Tracks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
func Track(attrs ...KV) *Element {
	return newVoidElem("track", attrs...)
}

// Video embeds a media player which supports video playback into the document. You can also use &lt;video&gt; for audio content, but the audio element may provide a more appropriate user experience.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
func Video(args ...any) *Element {
	return newElem("video", args...)
}

// Audio is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the source element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
func Audio(args ...any) *Element {
	return newElem("audio", args...)
}

// Canvas is a container element to use with either the canvas scripting API or the WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
func Canvas(args ...any) *Element {
	return newElem("canvas", args...)
}

// MapElement is used with &lt;area&gt; elements to define an image map (a clickable link area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
func MapElement(args ...any) *Element {
	return newElem("map", args...)
}

// Area defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hyperlink.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
func Area(attrs ...KV) *Element {
	return newVoidElem("area", attrs...)
}

// Svg is a container defining a new coordinate system and viewport. It is used as the outermost element of SVG documents, but it can also be used to embed an SVG fragment inside an SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/svg
func Svg(args ...any) *Element {
	return newElem("svg", args...)
}

// Math is the top-level element in MathML. Every valid MathML instance must be wrapped in it. In addition, you must not nest a second &lt;math&gt; element in another, but you can have an arbitrary number of other child elements in it.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/math
func Math(args ...any) *Element {
	return newElem("math", args...)
}

// Script is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The &lt;script&gt; element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
func Script(args ...any) *Element {
	return newElem("script", args...)
}

// Noscript defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
func Noscript(args ...any) *Element {
	return newElem("noscript", args...)
}

// Del represents a range of text that has been deleted from a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
func Del(args ...any) *Element {
	return newElem("del", args...)
}

// Ins represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
func Ins(args ...any) *Element {
	return newElem("ins", args...)
}

// Table represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
func Table(args ...any) *Element {
	return newElem("table", args...)
}

// Caption specifies the caption (or title) of a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
func Caption(args ...any) *Element {
	return newElem("caption", args...)
}

// Colgroup defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
func Colgroup(args ...any) *Element {
	return newElem("colgroup", args...)
}

// Col defines one or more columns in a column group represented by its implicit or explicit parent &lt;colgroup&gt; element. The &lt;col&gt; element is only valid as a child of a &lt;colgroup&gt; element that has no span attribute defined.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
func Col(attrs ...KV) *Element {
	return newVoidElem("col", attrs...)
}

// Thead groups the header content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
func Thead(args ...any) *Element {
	return newElem("thead", args...)
}

// Tbody groups the body content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
func Tbody(args ...any) *Element {
	return newElem("tbody", args...)
}

// Tfoot groups the footer content in a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
func Tfoot(args ...any) *Element {
	return newElem("tfoot", args...)
}

// Tr defines a row of cells in a table. The row's cells can then be established using a mix of &lt;td&gt; (data cell) and &lt;th&gt; (header cell) elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
func Tr(args ...any) *Element {
	return newElem("tr", args...)
}

// Th is a child of the &lt;tr&gt; element, it defines a cell as the header of a group of table cells. The nature of this group can be explicitly defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
func Th(args ...any) *Element {
	return newElem("th", args...)
}

// Td is a child of the &lt;tr&gt; element, it defines a cell of a table that contains data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
func Td(args ...any) *Element {
	return newElem("td", args...)
}

// Form represents a document section containing interactive controls for submitting information.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
func Form(args ...any) *Element {
	return newElem("form", args...)
}

// Fieldset is used to group several controls as well as labels (&lt;label&gt;) within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
func Fieldset(args ...any) *Element {
	return newElem("fieldset", args...)
}

// Legend represents a caption for the content of its parent &lt;fieldset&gt;.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
func Legend(args ...any) *Element {
	return newElem("legend", args...)
}

// Label represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
func Label(args ...any) *Element {
	return newElem("label", args...)
}

// Input is used to create interactive controls for web-based forms to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The &lt;input&gt; element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
func Input(attrs ...KV) *Element {
	return newVoidElem("input", attrs...)
}

// Button is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it performs an action, such as submitting a form or opening a dialog.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
func Button(args ...any) *Element {
	return newElem("button", args...)
}

// Select represents a control that provides a menu of options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
func Select(args ...any) *Element {
	return newElem("select", args...)
}

// Datalist contains a set of &lt;option&gt; elements that represent the permissible or recommended options available to choose from within other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
func Datalist(args ...any) *Element {
	return newElem("datalist", args...)
}

// Optgroup creates a grouping of options within a &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
func Optgroup(args ...any) *Element {
	return newElem("optgroup", args...)
}

// Option is used to define an item contained in a &lt;select&gt;, an &lt;optgroup&gt;, or a &lt;datalist&gt; element. As such, &lt;option&gt; can represent menu items in popups and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
func Option(args ...any) *Element {
	return newElem("option", args...)
}

// Textarea represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example, a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
func Textarea(args ...any) *Element {
	return newElem("textarea", args...)
}

// Output is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
func Output(args ...any) *Element {
	return newElem("output", args...)
}

// Progress displays an indicator showing the completion progress of a task, typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
func Progress(args ...any) *Element {
	return newElem("progress", args...)
}

// Meter represents either a scalar value within a known range or a fractional value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
func Meter(args ...any) *Element {
	return newElem("meter", args...)
}

// Details creates a disclosure widget in which information is visible only when the widget is toggled into an "open" state. A summary or label must be provided using the &lt;summary&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
func Details(args ...any) *Element {
	return newElem("details", args...)
}

// Summary specifies a summary, caption, or legend for a details element's disclosure box. Clicking the &lt;summary&gt; element toggles the state of the parent &lt;details&gt; element open and closed.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
func Summary(args ...any) *Element {
	return newElem("summary", args...)
}

// Dialog represents a dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
func Dialog(args ...any) *Element {
	return newElem("dialog", args...)
}

// Slot acts as a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
func Slot(args ...any) *Element {
	return newElem("slot", args...)
}

// Template holds HTML that is not to be rendered immediately when a page is loaded but may be instantiated subsequently during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
func Template(args ...any) *Element {
	return newElem("template", args...)
}

// Fencedframe represents a nested browsing context, like &lt;iframe&gt; but with more native privacy features built in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fencedframe
func Fencedframe(args ...any) *Element {
	return newElem("fencedframe", args...)
}

// Selectedcontent displays the content of the currently selected &lt;option&gt; inside a closed &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/selectedcontent
func Selectedcontent(args ...any) *Element {
	return newElem("selectedcontent", args...)
}

// Base specifies the base URL and default browsing context for relative URLs.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func Base(attrs ...KV) *Element {
	return newVoidElem("base", attrs...)
}

// Hgroup groups a set of h1–h6 elements when they represent a multi-level heading.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
func Hgroup(args ...any) *Element {
	return newElem("hgroup", args...)
}

// Address indicates contact information for a person or organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func Address(args ...any) *Element {
	return newElem("address", args...)
}

// Search represents a search or filtering interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
func Search(args ...any) *Element {
	return newElem("search", args...)
}

// Div is the generic container for flow content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
func Div(args ...any) *Element {
	return newElem("div", args...)
}

// Span is the generic inline container for phrasing content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
func Span(args ...any) *Element {
	return newElem("span", args...)
}

// P creates a paragraph element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
func P(args ...any) *Element {
	return newElem("p", args...)
}

// Dl represents a description list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
func Dl(args ...any) *Element {
	return newElem("dl", args...)
}

// Dt specifies a term in a description or definition list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
func Dt(args ...any) *Element {
	return newElem("dt", args...)
}

// Dd provides the description, definition, or value for the preceding term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
func Dd(args ...any) *Element {
	return newElem("dd", args...)
}

// Figure represents self-contained content with an optional caption.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
func Figure(args ...any) *Element {
	return newElem("figure", args...)
}

// Figcaption represents a caption or legend for the contents of its parent figure element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
func Figcaption(args ...any) *Element {
	return newElem("figcaption", args...)
}

// Menu represents a set of commands or options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
func Menu(args ...any) *Element {
	return newElem("menu", args...)
}

// Small represents side-comments and small print.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
func Small(args ...any) *Element {
	return newElem("small", args...)
}

// S renders text with a strikethrough.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
func S(args ...any) *Element {
	return newElem("s", args...)
}

// Cite marks the title of a creative work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
func Cite(args ...any) *Element {
	return newElem("cite", args...)
}

// Q indicates a short inline quotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
func Q(args ...any) *Element {
	return newElem("q", args...)
}

// Dfn indicates the defining instance of a term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
func Dfn(args ...any) *Element {
	return newElem("dfn", args...)
}

// Abbr represents an abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
func Abbr(args ...any) *Element {
	return newElem("abbr", args...)
}

// Ruby represents ruby annotations for East Asian typography.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
func Ruby(args ...any) *Element {
	return newElem("ruby", args...)
}

// Rt specifies the ruby text for ruby annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
func Rt(args ...any) *Element {
	return newElem("rt", args...)
}

// Rp provides parentheses for browsers that don't support ruby text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
func Rp(args ...any) *Element {
	return newElem("rp", args...)
}

// Data links content with a machine-readable translation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
func Data(args ...any) *Element {
	return newElem("data", args...)
}

// Time represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
func Time(args ...any) *Element {
	return newElem("time", args...)
}
