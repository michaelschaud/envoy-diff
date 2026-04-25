package filter

import "fmt"

// GroupConfig holds parsed CLI options for the group filter.
type GroupConfig struct {
	// Delimiter is the string used to split key names into group prefix and remainder.
	Delimiter string
}

// ParseGroupConfig parses the --group flag value.
// The flag value is the delimiter string (e.g. "_" or ".").
// An empty string disables grouping.
func ParseGroupConfig(spec string) (GroupConfig, error) {
	if len(spec) > 8 {
		return GroupConfig{}, fmt.Errorf("group delimiter too long (max 8 chars): %q", spec)
	}
	return GroupConfig{Delimiter: spec}, nil
}
