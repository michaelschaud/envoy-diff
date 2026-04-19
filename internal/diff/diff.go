package diff

import "sort"

// EnvMap represents a set of environment variables.
type EnvMap map[string]string

// Result holds the diff between two EnvMaps.
type Result struct {
	OnlyInLeft  map[string]string
	OnlyInRight map[string]string
	Modified    map[string][2]string // key -> [leftVal, rightVal]
	Unchanged   map[string]string
}

// Compare computes the diff between left (e.g. staging) and right (e.g. production).
func Compare(left, right EnvMap) Result {
	res := Result{
		OnlyInLeft:  make(map[string]string),
		OnlyInRight: make(map[string]string),
		Modified:    make(map[string][2]string),
		Unchanged:   make(map[string]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if lv == rv {
				res.Unchanged[k] = lv
			} else {
				res.Modified[k] = [2]string{lv, rv}
			}
		} else {
			res.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			res.OnlyInRight[k] = rv
		}
	}

	return res
}

// SortedKeys returns sorted keys of a map.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
