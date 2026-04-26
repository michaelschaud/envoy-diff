package filter

import (
	"strconv"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// TypecastMode controls how values are classified.
type TypecastMode string

const (
	TypecastBool   TypecastMode = "bool"
	TypecastInt    TypecastMode = "int"
	TypecastFloat  TypecastMode = "float"
	TypecastString TypecastMode = "string"
)

// ApplyTypecast annotates each value with its inferred type as a tag prefix,
// e.g. "[int] 42". Only modes listed in modes are applied; empty means all.
func ApplyTypecast(result diff.Result, modes []TypecastMode) diff.Result {
	if len(modes) == 0 {
		return result
	}
	active := make(map[TypecastMode]bool, len(modes))
	for _, m := range modes {
		active[m] = true
	}
	return diff.Result{
		Added:   typecastMap(result.Added, active),
		Removed: typecastMap(result.Removed, active),
		Changed: typecastChangedMap(result.Changed, active),
		Same:    typecastMap(result.Same, active),
	}
}

func inferType(v string) TypecastMode {
	v = strings.TrimSpace(v)
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "1", "0":
		return TypecastBool
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return TypecastInt
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypecastFloat
	}
	return TypecastString
}

func typecastMap(m map[string]string, active map[TypecastMode]bool) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		t := inferType(v)
		if active[t] {
			out[k] = "[" + string(t) + "] " + v
		} else {
			out[k] = v
		}
	}
	return out
}

func typecastChangedMap(m map[string][2]string, active map[TypecastMode]bool) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		old, nw := pair[0], pair[1]
		if t := inferType(old); active[t] {
			old = "[" + string(t) + "] " + old
		}
		if t := inferType(nw); active[t] {
			nw = "[" + string(t) + "] " + nw
		}
		out[k] = [2]string{old, nw}
	}
	return out
}
