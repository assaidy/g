package g

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple text",
			node:     Text("Hello World"),
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "Simple element",
			node:     Div(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Element with children",
			node:     Div().Add(Text("Hello"), P().Add(Text("World"))),
			expected: "<div>Hello<p>World</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element",
			node:     Br(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Empty element with children",
			node:     Empty().Add(Text("test")),
			expected: "test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Render(&buf, tt.node)

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && buf.String() != tt.expected {
				t.Errorf("Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestRender_ErrorHandling(t *testing.T) {
	// Test with a node that will cause an error during rendering
	// Since Text.Render() doesn't return errors in the current implementation,
	// we need to create an element that will cause an attribute error
	element := Div(KV{"invalid": nil})

	var buf bytes.Buffer
	err := Render(&buf, element)

	if err == nil {
		t.Error("Render() should return error for invalid attribute")
	}

	// Ensure nothing was written to buffer when error occurs
	if buf.Len() > 0 {
		t.Errorf("Render() should not write to buffer on error, got: %q", buf.String())
	}
}

func TestRender_WriteError(t *testing.T) {
	// Create a writer that will return an error on write
	errorWriter := &errorWriter{}
	node := Text("test")

	err := Render(errorWriter, node)
	if err == nil {
		t.Error("Render() should return error when writer fails")
	}

	if !strings.Contains(err.Error(), "write error") {
		t.Errorf("Render() should return writer error, got: %v", err)
	}
}

// errorWriter is a test helper that always returns an error on Write
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, error error) {
	return 0, &writeError{"write error"}
}

type writeError struct {
	msg string
}

func (e *writeError) Error() string {
	return e.msg
}

func TestRender_ComplexStructure(t *testing.T) {
	// Test with a complex nested structure to ensure it handles correctly
	node := Html(KV{"lang": "en"}).Add(
		Head().Add(
			Title().Add(Text("Test Page")),
		),
		Body().Add(
			Div(KV{"class": "container"}).Add(
				H1().Add(Text("Welcome")),
				P().Add(Text("This is a test.")),
				Ul().Add(
					Li().Add(Text("Item 1")),
					Li().Add(Text("Item 2")),
				),
			),
		),
	)

	expected := `<html lang="en"><head><title>Test Page</title></head><body><div class="container"><h1>Welcome</h1><p>This is a test.</p><ul><li>Item 1</li><li>Item 2</li></ul></div></body></html>`

	var buf bytes.Buffer
	err := Render(&buf, node)

	if err != nil {
		t.Errorf("Render() unexpected error: %v", err)
		return
	}

	if buf.String() != expected {
		t.Errorf("Render() complex structure = %q, want %q", buf.String(), expected)
	}
}
