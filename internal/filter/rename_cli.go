package filter

import (
	"fmt"
	"strings"
)

// ParseRenameRules parses a slice of "FROM=TO" strings into RenameRule values.
// Returns an error if any entry is malformed.
func ParseRenameRules(raw []string) ([]RenameRule, error) {
	rules := make([]RenameRule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid rename rule %q: expected FROM=TO", s)
		}
		rules = append(rules, RenameRule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
