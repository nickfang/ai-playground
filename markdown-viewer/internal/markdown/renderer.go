package markdown

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
)

var checkboxRegex = regexp.MustCompile(`(?m)^(\s*)-\s+\[([ xX/\-!])\]\s+(.*)$`)

// Renderer wraps Glamour with checkbox preprocessing.
type Renderer struct {
	glamour  *glamour.TermRenderer
	icons    map[CheckboxState]string
}

// NewRenderer creates a markdown renderer with the given checkbox icons.
func NewRenderer(icons map[CheckboxState]string, width int) (*Renderer, error) {
	gr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		glamour: gr,
		icons:   icons,
	}, nil
}

// bulletPrefix matches Glamour's bullet point rendering (e.g. "  • " or "• ").
var bulletPrefix = regexp.MustCompile(`^(\s*)•\s*`)

// Render processes markdown content, replacing checkbox syntax with themed icons,
// then renders the result through Glamour.
func (r *Renderer) Render(content string) (string, error) {
	processed := r.preprocessCheckboxes(content)
	rendered, err := r.glamour.Render(processed)
	if err != nil {
		return rendered, err
	}
	return r.stripBulletsFromCheckboxLines(rendered), nil
}

// RenderWithFilter renders only checkbox lines matching the given states.
// Non-checkbox content is included as-is.
// If states is nil or empty, all content is rendered.
func (r *Renderer) RenderWithFilter(content string, states map[CheckboxState]bool) (string, error) {
	if len(states) == 0 {
		return r.Render(content)
	}

	lines := strings.Split(content, "\n")
	var filtered []string

	for _, line := range lines {
		matches := checkboxRegex.FindStringSubmatch(line)
		if matches == nil {
			// Not a checkbox line — include it
			filtered = append(filtered, line)
			continue
		}

		marker := matches[2]
		state, ok := ParseMarker(marker)
		if !ok || !states[state] {
			continue
		}
		filtered = append(filtered, line)
	}

	return r.Render(strings.Join(filtered, "\n"))
}

// preprocessCheckboxes replaces - [x] style checkboxes with themed icons.
func (r *Renderer) preprocessCheckboxes(content string) string {
	return checkboxRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := checkboxRegex.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		indent := parts[1]
		marker := parts[2]
		text := parts[3]

		state, ok := ParseMarker(marker)
		if !ok {
			return match
		}

		icon := r.icons[state]
		// Keep "- " so Glamour renders each on its own line as a list item.
		// The bullet dot is stripped post-render.
		return indent + "- " + icon + " " + text
	})
}

// stripBulletsFromCheckboxLines removes Glamour's bullet character from lines
// that contain one of our checkbox icons, preserving indentation.
func (r *Renderer) stripBulletsFromCheckboxLines(rendered string) string {
	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		// Check if this line contains one of our checkbox icons
		hasIcon := false
		for _, icon := range r.icons {
			if strings.Contains(line, icon) {
				hasIcon = true
				break
			}
		}
		if hasIcon {
			// Replace "  • ✅ text" with "  ✅ text"
			lines[i] = bulletPrefix.ReplaceAllString(line, "$1")
		}
	}
	return strings.Join(lines, "\n")
}

// ExtractFirstHeading returns the first # heading from markdown content.
func ExtractFirstHeading(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}
	return ""
}
