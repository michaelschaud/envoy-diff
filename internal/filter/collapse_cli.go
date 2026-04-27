package filter

import (
	"fmt"
	"strings"
)

// ParseCollapseConfig parses a collapse spec of the form:
//
//	"<delimiter>" or "<delimiter>:<key1>,<key2>,..."
//
// Examples:
//
//	","              → collapse all keys on comma
//	",:PATH,CLASSPATH" → collapse only keys containing PATH or CLASSPATH
func ParseCollapseConfig(spec string) (CollapseConfig, error) {
	if spec == "" {
		return CollapseConfig{}, fmt.Errorf("collapse spec must not be empty")
	}

	parts := strings.SplitN(spec, ":", 2)
	delimiter := parts[0]
	if delimiter == "" {
		return CollapseConfig{}, fmt.Errorf("collapse delimiter must not be empty")
	}

	var keys []string
	if len(parts) == 2 && parts[1] != "" {
		for _, k := range strings.Split(parts[1], ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	return CollapseConfig{
		Delimiter: delimiter,
		Keys:      keys,
	}, nil
}
