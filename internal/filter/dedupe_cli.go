package filter

import "fmt"

// ParseDedupeMode parses a CLI string into a DedupeMode.
// Valid values: "first", "last", "all".
// An empty string disables deduplication.
func ParseDedupeMode(s string) (DedupeMode, error) {
	switch DedupeMode(s) {
	case DedupeKeepFirst, DedupeKeepLast, DedupeRemoveAll:
		return DedupeMode(s), nil
	case "":
		return "", nil
	default:
		return "", fmt.Errorf(
			"unknown dedupe mode %q: must be one of \"first\", \"last\", \"all\"",
			s,
		)
	}
}
