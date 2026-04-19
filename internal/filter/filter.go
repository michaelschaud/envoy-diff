package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Options holds filtering configuration.
type Options struct {
	Prefixes []string
	OnlyChanged bool
}

// Apply filters a diff.Result based on the provided Options.
// If Prefixes is non-empty, only keys matching at least one prefix are kept.
// If OnlyChanged is true, only Added, Removed, and Changed entries are kept.
func Apply(result diff.Result, opts Options) diff.Result {
	out := diff.Result{
		Added:     filterKeys(result.Added, opts.Prefixes),
		Removed:   filterKeys(result.Removed, opts.Prefixes),
		Changed:   filterChanged(result.Changed, opts.Prefixes),
		Unchanged: filterKeys(result.Unchanged, opts.Prefixes),
	}
	if opts.OnlyChanged {
		out.Unchanged = map[string]string{}
	}
	return out
}

func matchesPrefix(key string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func filterKeys(m map[string]string, prefixes []string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if matchesPrefix(k, prefixes) {
			out[k] = v
		}
	}
	return out
}

func filterChanged(m map[string][2]string, prefixes []string) map[string][2]string {
	out := make(map[string][2]string)
	for k, v := range m {
		if matchesPrefix(k, prefixes) {
			out[k] = v
		}
	}
	return out
}
