package diff

import "sort"

// Result holds the categorised diff between two env maps.
type Result struct {
	// OnlyInA contains keys present only in the first (staging) map.
	OnlyInA map[string]string
	// OnlyInB contains keys present only in the second (production) map.
	OnlyInB map[string]string
	// Changed contains keys present in both maps but with different values.
	// The value is [stagingVal, productionVal].
	Changed map[string][2]string
	// Common contains keys with identical values in both maps.
	Common map[string]string
}

// Compare diffs two environment variable maps, a (staging) and b (production).
func Compare(a, b map[string]string) Result {
	res := Result{
		OnlyInA: make(map[string]string),
		OnlyInB: make(map[string]string),
		Changed: make(map[string][2]string),
		Common:  make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok {
			res.OnlyInA[k] = va
		} else if va != vb {
			res.Changed[k] = [2]string{va, vb}
		} else {
			res.Common[k] = va
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			res.OnlyInB[k] = vb
		}
	}

	return res
}

// SortedKeys returns the keys of m in sorted order.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
