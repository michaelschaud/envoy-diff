package filter

import "github.com/your-org/envoy-diff/internal/diff"

// LimitOptions controls how many results are returned per category.
type LimitOptions struct {
	MaxAdded   int
	MaxRemoved int
	MaxChanged int
	MaxSame    int
}

// ApplyLimit truncates each category in the diff result to the specified
// maximum number of entries. A value of 0 means no limit is applied.
func ApplyLimit(result diff.Result, opts LimitOptions) diff.Result {
	return diff.Result{
		Added:   limitMap(result.Added, opts.MaxAdded),
		Removed: limitMap(result.Removed, opts.MaxRemoved),
		Changed: limitChanged(result.Changed, opts.MaxChanged),
		Same:    limitMap(result.Same, opts.MaxSame),
	}
}

func limitMap(m map[string]string, max int) map[string]string {
	if max <= 0 || len(m) <= max {
		return m
	}
	out := make(map[string]string, max)
	count := 0
	for k, v := range m {
		if count >= max {
			break
		}
		out[k] = v
		count++
	}
	return out
}

func limitChanged(m map[string]diff.ChangedValue, max int) map[string]diff.ChangedValue {
	if max <= 0 || len(m) <= max {
		return m
	}
	out := make(map[string]diff.ChangedValue, max)
	count := 0
	for k, v := range m {
		if count >= max {
			break
		}
		out[k] = v
		count++
	}
	return out
}
