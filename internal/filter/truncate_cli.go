package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// TruncateConfig holds parsed CLI options for the truncate filter.
type TruncateConfig struct {
	MaxLen int
	Suffix string
}

// ParseTruncateConfig parses a truncate spec of the form "<maxLen>[:<suffix>]".
// Examples:
//   - "20"        → maxLen=20, suffix="..."
//   - "10:***"    → maxLen=10, suffix="***"
func ParseTruncateConfig(spec string) (TruncateConfig, error) {
	if spec == "" {
		return TruncateConfig{}, fmt.Errorf("truncate spec must not be empty")
	}
	parts := strings.SplitN(spec, ":", 2)
	n, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return TruncateConfig{}, fmt.Errorf("invalid maxLen %q: %w", parts[0], err)
	}
	if n <= 0 {
		return TruncateConfig{}, fmt.Errorf("maxLen must be > 0, got %d", n)
	}
	suffix := "..."
	if len(parts) == 2 {
		suffix = parts[1]
	}
	return TruncateConfig{MaxLen: n, Suffix: suffix}, nil
}
