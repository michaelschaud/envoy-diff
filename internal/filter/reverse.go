package filter

import "github.com/yourorg/envoy-diff/internal/diff"

// ApplyReverse reverses the direction of the diff result — swapping staging and
// production perspectives. Added keys become removed, removed keys become added,
// and changed keys have their Old/New values swapped.
func ApplyReverse(result diff.Result) diff.Result {
	if !result.HasAny() {
		return result
	}

	reversed := diff.Result{
		Added:   reverseRemoved(result.Removed),
		Removed: reverseAdded(result.Added),
		Changed: reverseChanged(result.Changed),
		Same:    result.Same,
	}
	return reversed
}

// reverseAdded converts added entries (map[string]string) into removed entries.
func reverseAdded(added map[string]string) map[string]string {
	if len(added) == 0 {
		return nil
	}
	out := make(map[string]string, len(added))
	for k, v := range added {
		out[k] = v
	}
	return out
}

// reverseRemoved converts removed entries into added entries.
func reverseRemoved(removed map[string]string) map[string]string {
	if len(removed) == 0 {
		return nil
	}
	out := make(map[string]string, len(removed))
	for k, v := range removed {
		out[k] = v
	}
	return out
}

// reverseChanged swaps Old and New for every changed entry.
func reverseChanged(changed map[string]diff.Change) map[string]diff.Change {
	if len(changed) == 0 {
		return nil
	}
	out := make(map[string]diff.Change, len(changed))
	for k, c := range changed {
		out[k] = diff.Change{
			Old: c.New,
			New: c.Old,
		}
	}
	return out
}
