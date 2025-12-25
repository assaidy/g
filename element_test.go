package g

import (
	"strings"
	"testing"
)

func TestElement_Render(t *testing.T) {
	tests := []struct {
		name     string
		element  *Element
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple div",
			element:  Div(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Div with attributes",
			element:  Div(KV{"class": "container", "id": "main"}),
			expected: `<div class="container" id="main"></div>`,
			wantErr:  false,
		},
		{
			name: "Div with text child",
			element: func() *Element {
				return Div().Add(Text("Hello World")).(*Element)
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Div with multiple children",
			element: func() *Element {
				return Div().Add(Text("Hello"), Text(" "), Text("World")).(*Element)
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Nested elements",
			element: func() *Element {
				return Div().Add(P().Add(Text("Hello"))).(*Element)
			}(),
			expected: "<div><p>Hello</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element (br)",
			element:  Br(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Void element with attributes (img)",
			element:  Img(KV{"src": "test.jpg", "alt": "test"}),
			expected: `<img alt="test" src="test.jpg">`,
			wantErr:  false,
		},
		{
			name:     "Empty element",
			element:  Empty(),
			expected: "",
			wantErr:  false,
		},
		{
			name: "Empty element with children",
			element: func() *Element {
				return Empty().Add(Text("Hello")).(*Element)
			}(),
			expected: "Hello",
			wantErr:  false,
		},
		{
			name: "Boolean attribute true",
			element: Div(KV{
				"hidden": true,
				"class":  "test",
			}),
			expected: `<div class="test" hidden></div>`,
			wantErr:  false,
		},
		{
			name: "Boolean attribute false",
			element: Div(KV{
				"hidden": false,
				"class":  "test",
			}),
			expected: `<div class="test"></div>`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.element.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Element.Render() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestElement_Add(t *testing.T) {
	// Test adding children to non-void element
	div := Div()
	div.Add(Text("child1"), Text("child2"))

	if len(div.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(div.Children))
	}

	// Test that Add returns the element for chaining
	result := div.Add(Text("child3"))
	if result != div {
		t.Error("Add() should return the same element for chaining")
	}

	// Test adding to void element (should not add children)
	br := Br()
	br.Add(Text("should not be added"))

	if len(br.Children) != 0 {
		t.Error("Void elements should not accept children")
	}
}

func TestElement_renderAttrs(t *testing.T) {
	tests := []struct {
		name      string
		attrs     KV
		expected  string
		expectErr bool
	}{
		{
			name:      "String attributes",
			attrs:     KV{"class": "test", "id": "main"},
			expected:  ` class="test" id="main"`,
			expectErr: false,
		},
		{
			name:      "Boolean attributes",
			attrs:     KV{"hidden": true, "disabled": false},
			expected:  ` hidden`,
			expectErr: false,
		},
		{
			name:      "Nil value",
			attrs:     KV{"test": nil},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid value type",
			attrs:     KV{"test": 123},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty key",
			attrs:     KV{"": "value"},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Whitespace key",
			attrs:     KV{"   ": "value"},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Key with HTML escaping",
			attrs:     KV{"data-value": "<script>"},
			expected:  ` data-value="&lt;script&gt;"`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			element := &Element{Attrs: tt.attrs}
			var builder strings.Builder
			err := element.renderAttrs(&builder)

			if (err != nil) != tt.expectErr {
				t.Errorf("renderAttrs() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && builder.String() != tt.expected {
				t.Errorf("renderAttrs() = %q, want %q", builder.String(), tt.expected)
			}
		})
	}
}
