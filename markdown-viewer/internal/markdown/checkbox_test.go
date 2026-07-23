package markdown

import (
	"testing"
)

func TestAllStates(t *testing.T) {
	states := AllStates()
	if len(states) != 5 {
		t.Fatalf("expected 5 states, got %d", len(states))
	}

	expected := []CheckboxState{StateTodo, StateDone, StateInProgress, StateCancelled, StateUrgent}
	for i, s := range states {
		if s != expected[i] {
			t.Errorf("state[%d]: expected %d, got %d", i, expected[i], s)
		}
	}
}

func TestStateName(t *testing.T) {
	tests := []struct {
		state    CheckboxState
		expected string
	}{
		{StateTodo, "todo"},
		{StateDone, "done"},
		{StateInProgress, "inprogress"},
		{StateCancelled, "cancelled"},
		{StateUrgent, "urgent"},
		{CheckboxState(99), "unknown"},
	}

	for _, tt := range tests {
		got := StateName(tt.state)
		if got != tt.expected {
			t.Errorf("StateName(%d) = %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestParseStateName(t *testing.T) {
	tests := []struct {
		name     string
		expected CheckboxState
		ok       bool
	}{
		{"todo", StateTodo, true},
		{"done", StateDone, true},
		{"inprogress", StateInProgress, true},
		{"cancelled", StateCancelled, true},
		{"urgent", StateUrgent, true},
		{"invalid", StateTodo, false},
		{"", StateTodo, false},
		{"TODO", StateTodo, false}, // case-sensitive
	}

	for _, tt := range tests {
		state, ok := ParseStateName(tt.name)
		if state != tt.expected || ok != tt.ok {
			t.Errorf("ParseStateName(%q) = (%d, %v), want (%d, %v)", tt.name, state, ok, tt.expected, tt.ok)
		}
	}
}

func TestSyntaxMarker(t *testing.T) {
	tests := []struct {
		state    CheckboxState
		expected string
	}{
		{StateTodo, " "},
		{StateDone, "x"},
		{StateInProgress, "/"},
		{StateCancelled, "-"},
		{StateUrgent, "!"},
		{CheckboxState(99), " "}, // unknown defaults to space
	}

	for _, tt := range tests {
		got := SyntaxMarker(tt.state)
		if got != tt.expected {
			t.Errorf("SyntaxMarker(%d) = %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestParseMarker(t *testing.T) {
	tests := []struct {
		marker   string
		expected CheckboxState
		ok       bool
	}{
		{" ", StateTodo, true},
		{"x", StateDone, true},
		{"X", StateDone, true},
		{"/", StateInProgress, true},
		{"-", StateCancelled, true},
		{"!", StateUrgent, true},
		{"?", StateTodo, false},
		{"", StateTodo, false},
		{"xx", StateTodo, false},
	}

	for _, tt := range tests {
		state, ok := ParseMarker(tt.marker)
		if state != tt.expected || ok != tt.ok {
			t.Errorf("ParseMarker(%q) = (%d, %v), want (%d, %v)", tt.marker, state, ok, tt.expected, tt.ok)
		}
	}
}

func TestRoundTrip_StateNameAndParse(t *testing.T) {
	for _, state := range AllStates() {
		name := StateName(state)
		parsed, ok := ParseStateName(name)
		if !ok {
			t.Errorf("ParseStateName(%q) returned false for valid state", name)
		}
		if parsed != state {
			t.Errorf("round-trip failed: %d -> %q -> %d", state, name, parsed)
		}
	}
}

func TestRoundTrip_MarkerAndParse(t *testing.T) {
	for _, state := range AllStates() {
		marker := SyntaxMarker(state)
		parsed, ok := ParseMarker(marker)
		if !ok {
			t.Errorf("ParseMarker(%q) returned false for valid marker", marker)
		}
		if parsed != state {
			t.Errorf("round-trip failed: %d -> %q -> %d", state, marker, parsed)
		}
	}
}
