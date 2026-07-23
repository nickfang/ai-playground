package markdown

// CheckboxState represents one of the 5 supported checkbox states.
type CheckboxState int

const (
	StateTodo CheckboxState = iota
	StateDone
	StateInProgress
	StateCancelled
	StateUrgent
)

// AllStates returns all checkbox states in order.
func AllStates() []CheckboxState {
	return []CheckboxState{StateTodo, StateDone, StateInProgress, StateCancelled, StateUrgent}
}

// StateName returns the human-readable name of a checkbox state.
func StateName(s CheckboxState) string {
	switch s {
	case StateTodo:
		return "todo"
	case StateDone:
		return "done"
	case StateInProgress:
		return "inprogress"
	case StateCancelled:
		return "cancelled"
	case StateUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

// ParseStateName converts a name string to a CheckboxState.
// Returns StateTodo and false if the name is not recognized.
func ParseStateName(name string) (CheckboxState, bool) {
	switch name {
	case "todo":
		return StateTodo, true
	case "done":
		return StateDone, true
	case "inprogress":
		return StateInProgress, true
	case "cancelled":
		return StateCancelled, true
	case "urgent":
		return StateUrgent, true
	default:
		return StateTodo, false
	}
}

// SyntaxMarker returns the markdown syntax character for a state.
func SyntaxMarker(s CheckboxState) string {
	switch s {
	case StateTodo:
		return " "
	case StateDone:
		return "x"
	case StateInProgress:
		return "/"
	case StateCancelled:
		return "-"
	case StateUrgent:
		return "!"
	default:
		return " "
	}
}

// ParseMarker converts a single-character marker to a CheckboxState.
func ParseMarker(marker string) (CheckboxState, bool) {
	switch marker {
	case " ":
		return StateTodo, true
	case "x", "X":
		return StateDone, true
	case "/":
		return StateInProgress, true
	case "-":
		return StateCancelled, true
	case "!":
		return StateUrgent, true
	default:
		return StateTodo, false
	}
}
