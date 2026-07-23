package theme

import (
	"testing"

	"github.com/nfang/markdown-viewer/internal/markdown"
)

func TestEmojiTheme(t *testing.T) {
	th := EmojiTheme()

	expected := map[markdown.CheckboxState]string{
		markdown.StateTodo:       "⬜",
		markdown.StateDone:       "✅",
		markdown.StateInProgress: "🔄",
		markdown.StateCancelled:  "❌",
		markdown.StateUrgent:     "🔥",
	}

	for state, icon := range expected {
		got := th.Render(state)
		if got != icon {
			t.Errorf("EmojiTheme.Render(%d) = %q, want %q", state, got, icon)
		}
	}
}

func TestSymbolTheme(t *testing.T) {
	th := SymbolTheme()

	// Symbol theme uses lipgloss styling, so we just check they're non-empty
	for _, state := range markdown.AllStates() {
		got := th.Render(state)
		if got == "" {
			t.Errorf("SymbolTheme.Render(%d) returned empty string", state)
		}
		if got == "?" {
			t.Errorf("SymbolTheme.Render(%d) returned fallback '?'", state)
		}
	}
}

func TestCustomTheme_AllProvided(t *testing.T) {
	glyphs := map[string]string{
		"todo":       "○",
		"done":       "●",
		"inprogress": "◐",
		"cancelled":  "⊘",
		"urgent":     "⚠",
	}

	th := CustomTheme(glyphs)

	tests := []struct {
		state    markdown.CheckboxState
		expected string
	}{
		{markdown.StateTodo, "○"},
		{markdown.StateDone, "●"},
		{markdown.StateInProgress, "◐"},
		{markdown.StateCancelled, "⊘"},
		{markdown.StateUrgent, "⚠"},
	}

	for _, tt := range tests {
		got := th.Render(tt.state)
		if got != tt.expected {
			t.Errorf("CustomTheme.Render(%d) = %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestCustomTheme_PartialFallback(t *testing.T) {
	glyphs := map[string]string{
		"todo": "○",
		"done": "●",
		// inprogress, cancelled, urgent not provided — should fallback to emoji
	}

	th := CustomTheme(glyphs)

	if got := th.Render(markdown.StateTodo); got != "○" {
		t.Errorf("expected custom todo '○', got %q", got)
	}
	if got := th.Render(markdown.StateDone); got != "●" {
		t.Errorf("expected custom done '●', got %q", got)
	}
	// Fallbacks should be emoji
	if got := th.Render(markdown.StateInProgress); got != "🔄" {
		t.Errorf("expected emoji fallback '🔄', got %q", got)
	}
	if got := th.Render(markdown.StateCancelled); got != "❌" {
		t.Errorf("expected emoji fallback '❌', got %q", got)
	}
	if got := th.Render(markdown.StateUrgent); got != "🔥" {
		t.Errorf("expected emoji fallback '🔥', got %q", got)
	}
}

func TestCustomTheme_EmptyGlyphs(t *testing.T) {
	th := CustomTheme(map[string]string{})

	// All should fallback to emoji
	emoji := EmojiTheme()
	for _, state := range markdown.AllStates() {
		got := th.Render(state)
		expected := emoji.Render(state)
		if got != expected {
			t.Errorf("CustomTheme(empty).Render(%d) = %q, want emoji %q", state, got, expected)
		}
	}
}

func TestCustomTheme_NilGlyphs(t *testing.T) {
	th := CustomTheme(nil)

	// All should fallback to emoji
	for _, state := range markdown.AllStates() {
		got := th.Render(state)
		if got == "?" {
			t.Errorf("CustomTheme(nil).Render(%d) returned fallback '?'", state)
		}
	}
}

func TestForConfig(t *testing.T) {
	tests := []struct {
		name        string
		themeName   string
		glyphs      map[string]string
		checkState  markdown.CheckboxState
		expectEmoji bool
	}{
		{"emoji default", "emoji", nil, markdown.StateTodo, true},
		{"explicit emoji", "emoji", nil, markdown.StateDone, true},
		{"invalid falls back to emoji", "bogus", nil, markdown.StateUrgent, true},
		{"empty falls back to emoji", "", nil, markdown.StateTodo, true},
	}

	emoji := EmojiTheme()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := ForConfig(tt.themeName, tt.glyphs)
			got := th.Render(tt.checkState)
			if tt.expectEmoji {
				expected := emoji.Render(tt.checkState)
				if got != expected {
					t.Errorf("ForConfig(%q).Render(%d) = %q, want %q", tt.themeName, tt.checkState, got, expected)
				}
			}
		})
	}
}

func TestForConfig_Symbols(t *testing.T) {
	th := ForConfig("symbols", nil)
	// Symbol theme should produce non-empty, non-emoji output
	got := th.Render(markdown.StateTodo)
	if got == "" || got == "?" {
		t.Errorf("symbols theme should produce output, got %q", got)
	}
}

func TestForConfig_Custom(t *testing.T) {
	glyphs := map[string]string{"todo": "★"}
	th := ForConfig("custom", glyphs)
	got := th.Render(markdown.StateTodo)
	if got != "★" {
		t.Errorf("custom theme should use glyph '★', got %q", got)
	}
}

func TestRender_UnknownState(t *testing.T) {
	th := EmojiTheme()
	got := th.Render(markdown.CheckboxState(99))
	if got != "?" {
		t.Errorf("Render with unknown state should return '?', got %q", got)
	}
}
