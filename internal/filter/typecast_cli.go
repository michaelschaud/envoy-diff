package filter

import (
	"fmt"
	"strings"
)

// ParseTypecastModes parses a comma-separated list of typecast mode names.
// Valid values: bool, int, float, string.
// Unknown values are returned as an error.
func ParseTypecastModes(spec string) ([]TypecastMode, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, nil
	}
	parts := strings.Split(spec, ",")
	valid := map[string]TypecastMode{
		"bool":   TypecastBool,
		"int":    TypecastInt,
		"float":  TypecastFloat,
		"string": TypecastString,
	}
	modes := make([]TypecastMode, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "" {
			continue
		}
		m, ok := valid[p]
		if !ok {
			return nil, fmt.Errorf("unknown typecast mode %q; valid: bool, int, float, string", p)
		}
		modes = append(modes, m)
	}
	return modes, nil
}
