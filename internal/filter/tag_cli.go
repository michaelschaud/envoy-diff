package filter

import (
	"fmt"
	"strings"
)

// ParseTagRules parses a slice of "tag=substr" strings into TagRule values.
// Each entry must contain exactly one "=" separator with non-empty tag and substr.
func ParseTagRules(specs []string) ([]TagRule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]TagRule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag rule %q: expected format tag=substr", spec)
		}
		tag := strings.TrimSpace(parts[0])
		substr := strings.TrimSpace(parts[1])
		if tag == "" {
			return nil, fmt.Errorf("invalid tag rule %q: tag label must not be empty", spec)
		}
		if substr == "" {
			return nil, fmt.Errorf("invalid tag rule %q: key substring must not be empty", spec)
		}
		rules = append(rules, TagRule{Tag: tag, Substr: substr})
	}
	return rules, nil
}
