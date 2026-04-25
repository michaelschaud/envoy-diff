package filter

import (
	"fmt"
	"strings"
)

// ParseHighlightConfig parses a highlight spec of the form:
//
//	"<marker>:<substring1>,<substring2>,..."
//
// Example: "[!]:SECRET,TOKEN,PASSWORD"
func ParseHighlightConfig(spec string) (HighlightConfig, error) {
	if spec == "" {
		return HighlightConfig{}, fmt.Errorf("highlight spec must not be empty")
	}

	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 {
		return HighlightConfig{}, fmt.Errorf("highlight spec must be in format <marker>:<substrings>, got %q", spec)
	}

	marker := parts[0]
	if marker == "" {
		return HighlightConfig{}, fmt.Errorf("marker must not be empty")
	}

	raw := strings.TrimSpace(parts[1])
	if raw == "" {
		return HighlightConfig{}, fmt.Errorf("at least one key substring is required")
	}

	var substrings []string
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			substrings = append(substrings, s)
		}
	}

	if len(substrings) == 0 {
		return HighlightConfig{}, fmt.Errorf("no valid substrings found in spec %q", spec)
	}

	return HighlightConfig{Marker: marker, Substrings: substrings}, nil
}
