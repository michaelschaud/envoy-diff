package filter

import (
	"fmt"
	"strings"
)

// ParseMergeStrategy converts a strategy string into a MergeStrategy.
// Returns an error if the value is unrecognised.
func ParseMergeStrategy(s string) (MergeStrategy, error) {
	switch MergeStrategy(strings.ToLower(strings.TrimSpace(s))) {
	case MergeStrategyLeft:
		return MergeStrategyLeft, nil
	case MergeStrategyRight:
		return MergeStrategyRight, nil
	case MergeStrategyUnion:
		return MergeStrategyUnion, nil
	default:
		return "", fmt.Errorf("unknown merge strategy %q: must be left, right, or union", s)
	}
}

// ParseMergeOverlay parses a slice of "KEY=VALUE" strings into a map.
// Lines that do not contain "=" are skipped.
func ParseMergeOverlay(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(p[:idx])
		val := strings.TrimSpace(p[idx+1:])
		if key != "" {
			out[key] = val
		}
	}
	return out
}
