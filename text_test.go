package g

import (
	"testing"
)

func TestText_Render(t *testing.T) {
	tests := []struct {
		name     string
		text     Text
		expected string
	}{
		{
			name:     "Simple text",
			text:     Text("Hello World"),
			expected: "Hello World",
		},
		{
			name:     "Empty text",
			text:     Text(""),
			expected: "",
		},
		{
			name:     "Text with HTML entities",
			text:     Text("<script>alert('xss')</script>"),
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "Text with quotes",
			text:     Text("Hello \"World\" & 'Universe'"),
			expected: "Hello &#34;World&#34; &amp; &#39;Universe&#39;",
		},
		{
			name:     "Text with leading space",
			text:     Text("  hello"),
			expected: " hello",
		},
		{
			name:     "Text with trailing space",
			text:     Text("hello  "),
			expected: "hello ",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     Text("  hello  "),
			expected: " hello ",
		},
		{
			name:     "Text with multiple spaces inside",
			text:     Text("hello    world"),
			expected: "hello world",
		},
		{
			name:     "Text with newlines and tabs",
			text:     Text("hello\n\tworld"),
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.text.Render()
			if err != nil {
				t.Errorf("Text.Render() returned error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Text.Render() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestText_Render_Error(t *testing.T) {
	// Text.Render() should never return an error based on the implementation
	// This test ensures that behavior
	text := Text("test")
	_, err := text.Render()
	if err != nil {
		t.Errorf("Text.Render() should not return error, got: %v", err)
	}
}
