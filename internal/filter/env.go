package filter

import (
	"os"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyEnvOverride replaces values in the result with those from the current
// process environment when a matching key exists. This allows local env vars
// to override loaded file values before diffing.
func ApplyEnvOverride(result diff.Result) diff.Result {
	result.Added = overrideMap(result.Added)
	result.Removed = overrideMap(result.Removed)
	result.Changed = overrideChanged(result.Changed)
	result.Same = overrideMap(result.Same)
	return result
}

func overrideMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if env, ok := os.LookupEnv(k); ok {
			out[k] = env
		} else {
			out[k] = v
		}
	}
	return out
}

func overrideChanged(m map[string]diff.ChangedValue) map[string]diff.ChangedValue {
	out := make(map[string]diff.ChangedValue, len(m))
	for k, cv := range m {
		if env, ok := os.LookupEnv(k); ok {
			out[k] = diff.ChangedValue{Old: cv.Old, New: env}
		} else {
			out[k] = cv
		}
	}
	return out
}

// EnvVarNames returns a sorted list of all environment variable names currently
// set in the process, optionally filtered by a prefix (case-insensitive).
func EnvVarNames(prefix string) []string {
	var names []string
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) < 1 {
			continue
		}
		name := parts[0]
		if prefix == "" || strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			names = append(names, name)
		}
	}
	return names
}
