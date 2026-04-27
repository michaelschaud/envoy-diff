package filter

import (
	"fmt"
	"strings"
)

// ParseUniqueModes parses a comma-separated list of unique mode names.
// Valid values are "added", "removed", and "same".
// An empty spec returns nil (all categories will be processed).
func ParseUniqueModes(spec string) ([]UniqueMode, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, nil
	}

	parts := strings.Split(spec, ",")
	modes := make([]UniqueMode, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		switch UniqueMode(p) {
		case UniqueModeAdded, UniqueModeRemoved, UniqueModeSame:
			modes = append(modes, UniqueMode(p))
		default:
			return nil, fmt.Errorf("unknown unique mode %q: must be one of added, removed, same", p)
		}
	}
	return modes, nil
}
