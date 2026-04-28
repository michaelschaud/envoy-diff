package filter

import (
	"fmt"
	"strconv"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// CapConfig controls the maximum number of characters allowed per value.
type CapConfig struct {
	MaxLen      int
	Replacement string
}

// ApplyCap limits the length of values in Added, Removed, Same, and Changed
// categories. Values exceeding MaxLen are replaced with Replacement.
func ApplyCap(result diff.Result, cfg CapConfig) diff.Result {
	if cfg.MaxLen <= 0 {
		return result
	}
	if cfg.Replacement == "" {
		cfg.Replacement = "[CAPPED]"
	}
	return diff.Result{
		Added:   capMap(result.Added, cfg),
		Removed: capMap(result.Removed, cfg),
		Same:    capMap(result.Same, cfg),
		Changed: capChanged(result.Changed, cfg),
	}
}

func capMap(m map[string]string, cfg CapConfig) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if len(v) > cfg.MaxLen {
			out[k] = cfg.Replacement
		} else {
			out[k] = v
		}
	}
	return out
}

func capChanged(m map[string][2]string, cfg CapConfig) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		old, nw := pair[0], pair[1]
		if len(old) > cfg.MaxLen {
			old = cfg.Replacement
		}
		if len(nw) > cfg.MaxLen {
			nw = cfg.Replacement
		}
		out[k] = [2]string{old, nw}
	}
	return out
}

// ParseCapConfig parses a spec string of the form "maxlen" or "maxlen:replacement".
func ParseCapConfig(spec string) (CapConfig, error) {
	if spec == "" {
		return CapConfig{}, nil
	}
	var raw, replacement string
	for i := 0; i < len(spec); i++ {
		if spec[i] == ':' {
			raw = spec[:i]
			replacement = spec[i+1:]
			break
		}
	}
	if raw == "" {
		raw = spec
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return CapConfig{}, fmt.Errorf("cap: invalid maxlen %q", raw)
	}
	return CapConfig{MaxLen: n, Replacement: replacement}, nil
}
