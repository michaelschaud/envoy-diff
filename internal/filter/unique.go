package filter

import "github.com/your-org/envoy-diff/internal/diff"

// UniqueMode controls which side to keep when deduplicating by value.
type UniqueMode string

const (
	UniqueModeAdded   UniqueMode = "added"
	UniqueModeRemoved UniqueMode = "removed"
	UniqueModeSame    UniqueMode = "same"
)

// ApplyUnique removes entries whose values appear more than once across the
// result set. Only keys whose values are unique within their category are
// retained. If modes is empty, all categories are deduplicated by value.
func ApplyUnique(result diff.Result, modes []UniqueMode) diff.Result {
	if len(modes) == 0 {
		return diff.Result{
			Added:   uniqueByValue(result.Added),
			Removed: uniqueByValue(result.Removed),
			Same:    uniqueByValue(result.Same),
			Changed: uniqueChangedByValue(result.Changed),
		}
	}

	out := diff.Result{
		Added:   copyMap(result.Added),
		Removed: copyMap(result.Removed),
		Same:    copyMap(result.Same),
		Changed: copyChangedMap(result.Changed),
	}

	for _, m := range modes {
		switch m {
		case UniqueModeAdded:
			out.Added = uniqueByValue(result.Added)
		case UniqueModeRemoved:
			out.Removed = uniqueByValue(result.Removed)
		case UniqueModeSame:
			out.Same = uniqueByValue(result.Same)
		}
	}
	return out
}

func uniqueByValue(m map[string]string) map[string]string {
	freq := make(map[string]int, len(m))
	for _, v := range m {
		freq[v]++
	}
	out := make(map[string]string)
	for k, v := range m {
		if freq[v] == 1 {
			out[k] = v
		}
	}
	return out
}

func uniqueChangedByValue(m map[string][2]string) map[string][2]string {
	oldFreq := make(map[string]int)
	newFreq := make(map[string]int)
	for _, pair := range m {
		oldFreq[pair[0]]++
		newFreq[pair[1]]++
	}
	out := make(map[string][2]string)
	for k, pair := range m {
		if oldFreq[pair[0]] == 1 && newFreq[pair[1]] == 1 {
			out[k] = pair
		}
	}
	return out
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copyChangedMap(m map[string][2]string) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
