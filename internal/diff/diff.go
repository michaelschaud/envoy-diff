package diff

import "sort"

// Result holds the categorised diff between two env maps.
type Result struct {
	Added      []string
	Removed    []string
	Changed    []string
	Identical  []string
	Staging    map[string]string
	Production map[string]string
}

// HasDiff returns true if there are any added, removed, or changed keys.
func (r Result) HasDiff() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Compare compares staging and production env maps and returns a Result.
func Compare(staging, production map[string]string) Result {
	result := Result{
		Staging:    staging,
		Production: production,
	}

	stagingKeys := keySet(staging)
	prodKeys := keySet(production)

	for k := range prodKeys {
		if _, ok := stagingKeys[k]; !ok {
			result.Added = append(result.Added, k)
		}
	}

	for k := range stagingKeys {
		if _, ok := prodKeys[k]; !ok {
			result.Removed = append(result.Removed, k)
		} else if staging[k] != production[k] {
			result.Changed = append(result.Changed, k)
		} else {
			result.Identical = append(result.Identical, k)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)
	sort.Strings(result.Identical)

	return result
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

func keySet(m map[string]string) map[string]struct{} {
	s := make(map[string]struct{}, len(m))
	for k := range m {
		s[k] = struct{}{}
	}
	return s
}
