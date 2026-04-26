package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// DiffMode controls which categories are retained after applying the diff filter.
type DiffMode string

const (
	DiffModeAdded   DiffMode = "added"
	DiffModeRemoved DiffMode = "removed"
	DiffModeChanged DiffMode = "changed"
	DiffModeSame    DiffMode = "same"
)

// ParseDiffModes parses a comma-separated list of diff mode names.
// Unknown modes are silently ignored. An empty slice means no filtering.
func ParseDiffModes(raw string) []DiffMode {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var modes []DiffMode
	for _, p := range parts {
		switch DiffMode(strings.TrimSpace(strings.ToLower(p))) {
		case DiffModeAdded, DiffModeRemoved, DiffModeChanged, DiffModeSame:
			modes = append(modes, DiffMode(strings.TrimSpace(strings.ToLower(p))))
		}
	}
	return modes
}

// ApplyDiffModeFilter retains only the result categories specified by modes.
// If modes is empty the result is returned unchanged.
func ApplyDiffModeFilter(result diff.Result, modes []DiffMode) diff.Result {
	if len(modes) == 0 {
		return result
	}

	want := make(map[DiffMode]bool, len(modes))
	for _, m := range modes {
		want[m] = true
	}

	out := diff.Result{}

	if want[DiffModeAdded] {
		out.Added = result.Added
	} else {
		out.Added = map[string]string{}
	}

	if want[DiffModeRemoved] {
		out.Removed = result.Removed
	} else {
		out.Removed = map[string]string{}
	}

	if want[DiffModeChanged] {
		out.Changed = result.Changed
	} else {
		out.Changed = map[string]diff.ChangedValue{}
	}

	if want[DiffModeSame] {
		out.Same = result.Same
	} else {
		out.Same = map[string]string{}
	}

	return out
}
