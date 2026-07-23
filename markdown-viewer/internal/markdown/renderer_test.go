package markdown

import (
	"strings"
	"testing"
)

func testIcons() map[CheckboxState]string {
	return map[CheckboxState]string{
		StateTodo:       "⬜",
		StateDone:       "✅",
		StateInProgress: "🔄",
		StateCancelled:  "❌",
		StateUrgent:     "🔥",
	}
}

func TestNewRenderer(t *testing.T) {
	r, err := NewRenderer(testIcons(), 80)
	if err != nil {
		t.Fatalf("NewRenderer returned error: %v", err)
	}
	if r == nil {
		t.Fatal("NewRenderer returned nil renderer")
	}
}

func TestRender_BasicMarkdown(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	output, err := r.Render("# Hello\n\nSome text")
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if !strings.Contains(output, "Hello") {
		t.Error("rendered output should contain heading text 'Hello'")
	}
}

func TestRender_CheckboxReplacement(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "- [ ] Todo item\n- [x] Done item\n- [/] In progress\n- [-] Cancelled\n- [!] Urgent"
	output, err := r.Render(input)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// The checkbox markers should be replaced with icons
	if strings.Contains(output, "[ ]") {
		t.Error("output should not contain raw '[ ]' checkbox syntax")
	}
	if !strings.Contains(output, "⬜") {
		t.Error("output should contain ⬜ for todo items")
	}
	if !strings.Contains(output, "✅") {
		t.Error("output should contain ✅ for done items")
	}
	if !strings.Contains(output, "🔄") {
		t.Error("output should contain 🔄 for in-progress items")
	}
	if !strings.Contains(output, "❌") {
		t.Error("output should contain ❌ for cancelled items")
	}
	if !strings.Contains(output, "🔥") {
		t.Error("output should contain 🔥 for urgent items")
	}
}

func TestRender_IndentedCheckboxes(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "  - [ ] Indented todo\n    - [x] Deeply indented done"
	output, err := r.Render(input)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if !strings.Contains(output, "⬜") {
		t.Error("indented todo should be replaced with icon")
	}
	if !strings.Contains(output, "✅") {
		t.Error("deeply indented done should be replaced with icon")
	}
}

func TestRender_MixedContent(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "# Title\n\n- [ ] A task\n- Regular list item\n\n> A quote"
	output, err := r.Render(input)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if !strings.Contains(output, "⬜") {
		t.Error("checkbox should be replaced")
	}
	if !strings.Contains(output, "Regular list item") {
		t.Error("regular list items should be preserved")
	}
}

func TestRenderWithFilter_NoFilter(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "- [ ] Todo\n- [x] Done"
	output, err := r.RenderWithFilter(input, nil)
	if err != nil {
		t.Fatalf("RenderWithFilter error: %v", err)
	}

	// With nil filter, all content should render
	if !strings.Contains(output, "⬜") {
		t.Error("todo should be present with nil filter")
	}
	if !strings.Contains(output, "✅") {
		t.Error("done should be present with nil filter")
	}
}

func TestRenderWithFilter_EmptyFilter(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "- [ ] Todo\n- [x] Done"
	output, err := r.RenderWithFilter(input, map[CheckboxState]bool{})
	if err != nil {
		t.Fatalf("RenderWithFilter error: %v", err)
	}

	// Empty map should render all
	if !strings.Contains(output, "⬜") {
		t.Error("todo should be present with empty filter")
	}
}

func TestRenderWithFilter_SingleState(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "- [ ] Todo item\n- [x] Done item\n- [!] Urgent item"
	filter := map[CheckboxState]bool{StateUrgent: true}
	output, err := r.RenderWithFilter(input, filter)
	if err != nil {
		t.Fatalf("RenderWithFilter error: %v", err)
	}

	if !strings.Contains(output, "🔥") {
		t.Error("urgent items should be present")
	}
	if strings.Contains(output, "⬜") {
		t.Error("todo items should be filtered out")
	}
	if strings.Contains(output, "✅") {
		t.Error("done items should be filtered out")
	}
}

func TestRenderWithFilter_MultipleStates(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "- [ ] Todo\n- [x] Done\n- [/] Progress\n- [-] Cancelled\n- [!] Urgent"
	filter := map[CheckboxState]bool{
		StateTodo:   true,
		StateUrgent: true,
	}
	output, err := r.RenderWithFilter(input, filter)
	if err != nil {
		t.Fatalf("RenderWithFilter error: %v", err)
	}

	if !strings.Contains(output, "⬜") {
		t.Error("todo should be present")
	}
	if !strings.Contains(output, "🔥") {
		t.Error("urgent should be present")
	}
	if strings.Contains(output, "✅") {
		t.Error("done should be filtered out")
	}
}

func TestRenderWithFilter_PreservesNonCheckboxContent(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	input := "# Heading\n\n- [ ] Todo\n- [x] Done\n\nSome paragraph"
	filter := map[CheckboxState]bool{StateTodo: true}
	output, err := r.RenderWithFilter(input, filter)
	if err != nil {
		t.Fatalf("RenderWithFilter error: %v", err)
	}

	if !strings.Contains(output, "Heading") {
		t.Error("headings should be preserved")
	}
	if !strings.Contains(output, "paragraph") {
		t.Error("paragraphs should be preserved")
	}
}

func TestPreprocessCheckboxes(t *testing.T) {
	r, _ := NewRenderer(testIcons(), 80)

	tests := []struct {
		input    string
		contains string
		excludes string
	}{
		{"- [ ] Todo", "⬜", "[ ]"},
		{"- [x] Done", "✅", "[x]"},
		{"- [X] Done upper", "✅", "[X]"},
		{"- [/] Progress", "🔄", "[/]"},
		{"- [-] Cancelled", "❌", "[-]"},
		{"- [!] Urgent", "🔥", "[!]"},
		{"- [?] Unknown", "[?]", ""}, // unknown marker left as-is
		{"Regular text", "Regular text", ""},
	}

	for _, tt := range tests {
		result := r.preprocessCheckboxes(tt.input)
		if tt.contains != "" && !strings.Contains(result, tt.contains) {
			t.Errorf("preprocessCheckboxes(%q) should contain %q, got %q", tt.input, tt.contains, result)
		}
		if tt.excludes != "" && strings.Contains(result, tt.excludes) {
			t.Errorf("preprocessCheckboxes(%q) should not contain %q, got %q", tt.input, tt.excludes, result)
		}
	}
}

func TestExtractFirstHeading(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"# Hello World", "Hello World"},
		{"## Not first level\n# First Level", "First Level"},
		{"No heading here", ""},
		{"", ""},
		{"# ", ""},
		{"#NoSpace", ""},
		{"Some text\n# Heading After", "Heading After"},
		{"# First\n# Second", "First"},
	}

	for _, tt := range tests {
		got := ExtractFirstHeading(tt.content)
		if got != tt.expected {
			t.Errorf("ExtractFirstHeading(%q) = %q, want %q", tt.content, got, tt.expected)
		}
	}
}
