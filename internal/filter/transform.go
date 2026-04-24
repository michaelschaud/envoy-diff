package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// TransformMode defines how values should be transformed.
type TransformMode string

const (
	TransformUpper TransformMode = "upper"
	TransformLower TransformMode = "lower"
	TransformNone  TransformMode = ""
)

// ApplyTransform applies a case transformation to all values in the diff result.
// Supported modes: "upper", "lower". Unknown modes are treated as no-op.
func ApplyTransform(result diff.Result, mode TransformMode) diff.Result {
	if mode == TransformNone {
		return result
	}

	transformFn := resolveTransform(mode)

	return diff.Result{
		Added:   transformMap(result.Added, transformFn),
		Removed: transformMap(result.Removed, transformFn),
		Changed: transformChanged(result.Changed, transformFn),
		Same:    transformMap(result.Same, transformFn),
	}
}

func resolveTransform(mode TransformMode) func(string) string {
	switch mode {
	case TransformUpper:
		return strings.ToUpper
	case TransformLower:
		return strings.ToLower
	default:
		return func(s string) string { return s }
	}
}

func transformMap(m map[string]string, fn func(string) string) map[string]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = fn(v)
	}
	return out
}

func transformChanged(m map[string][2]string, fn func(string) string) map[string][2]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		out[k] = [2]string{fn(pair[0]), fn(pair[1])}
	}
	return out
}
