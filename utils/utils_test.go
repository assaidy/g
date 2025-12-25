package utils

import (
	"testing"

	"github.com/assaidy/g"
)

func TestIfElse(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		result    string
		alt       string
		expected  string
	}{
		{
			name:      "Condition true returns result",
			condition: true,
			result:    "yes",
			alt:       "no",
			expected:  "yes",
		},
		{
			name:      "Condition false returns alternative",
			condition: false,
			result:    "yes",
			alt:       "no",
			expected:  "no",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfElse(tt.condition, tt.result, tt.alt)
			if result != tt.expected {
				t.Errorf("IfElse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIfElse_Nodes(t *testing.T) {
	trueNode := g.Div().Add(g.Text("true"))
	falseNode := g.P().Add(g.Text("false"))

	tests := []struct {
		name      string
		condition bool
		expected  string
	}{
		{
			name:      "Condition true returns node",
			condition: true,
			expected:  "<div>true</div>",
		},
		{
			name:      "Condition false returns node",
			condition: false,
			expected:  "<p>false</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := IfElse(tt.condition, trueNode, falseNode)
			result, err := node.Render()
			if err != nil {
				t.Errorf("IfElse() node render error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("IfElse() node render = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIf(t *testing.T) {
	node := g.Div().Add(g.Text("content"))

	tests := []struct {
		name      string
		condition bool
		expected  string
	}{
		{
			name:      "Condition true returns node",
			condition: true,
			expected:  "<div>content</div>",
		},
		{
			name:      "Condition false returns empty",
			condition: false,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := If(tt.condition, node)
			result, err := resultNode.Render()
			if err != nil {
				t.Errorf("If() node render error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("If() node render = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		f        func() g.Node
		expected string
	}{
		{
			name:     "Repeat zero times",
			n:        0,
			f:        func() g.Node { return g.Div() },
			expected: "",
		},
		{
			name:     "Repeat once",
			n:        1,
			f:        func() g.Node { return g.Div().Add(g.Text("item")) },
			expected: "<div>item</div>",
		},
		{
			name:     "Repeat multiple times",
			n:        3,
			f:        func() g.Node { return g.Div().Add(g.Text("item")) },
			expected: "<div>item</div><div>item</div><div>item</div>",
		},
		{
			name: "Repeat with different content",
			n:    2,
			f: func() g.Node {
				static := 0
				static++
				return g.Div().Add(g.Text(string(rune('a' + static))))
			},
			expected: "<div>b</div><div>b</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := Repeat(tt.n, tt.f)
			result, err := resultNode.Render()
			if err != nil {
				t.Errorf("Repeat() node render error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Repeat() node render = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		f        func(string) g.Node
		expected string
	}{
		{
			name:     "Map empty slice",
			input:    []string{},
			f:        func(s string) g.Node { return g.Li().Add(g.Text(s)) },
			expected: "",
		},
		{
			name:     "Map single item",
			input:    []string{"apple"},
			f:        func(s string) g.Node { return g.Li().Add(g.Text(s)) },
			expected: "<li>apple</li>",
		},
		{
			name:     "Map multiple items",
			input:    []string{"apple", "banana", "cherry"},
			f:        func(s string) g.Node { return g.Li().Add(g.Text(s)) },
			expected: "<li>apple</li><li>banana</li><li>cherry</li>",
		},
		{
			name:  "Map with conditional logic",
			input: []string{"apple", "banana"},
			f: func(s string) g.Node {
				if s == "apple" {
					return g.Li().Add(g.Text(s), g.Span().Add(g.Text(" (popular)")))
				}
				return g.Li().Add(g.Text(s))
			},
			expected: "<li>apple<span> (popular)</span></li><li>banana</li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNode := Map(tt.input, tt.f)
			result, err := resultNode.Render()
			if err != nil {
				t.Errorf("Map() node render error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Map() node render = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMap_Integers(t *testing.T) {
	numbers := []int{1, 2, 3}
	resultNode := Map(numbers, func(n int) g.Node {
		return g.Div().Add(g.Text(string(rune('0' + n))))
	})

	result, err := resultNode.Render()
	if err != nil {
		t.Errorf("Map() integers node render error: %v", err)
		return
	}
	expected := "<div>1</div><div>2</div><div>3</div>"
	if result != expected {
		t.Errorf("Map() integers node render = %v, want %v", result, expected)
	}
}
