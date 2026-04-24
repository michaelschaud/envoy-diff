package filter

import "github.com/your-org/envoy-diff/internal/diff"

// PaginateOptions controls which page of results to return.
type PaginateOptions struct {
	Page     int // 1-based page number
	PageSize int // number of entries per page (0 means no pagination)
}

// ApplyPaginate slices each category of the diff result to the requested page.
// Keys are sorted before slicing so that pagination is stable across runs.
// If opts.PageSize <= 0 the result is returned unchanged.
func ApplyPaginate(result diff.Result, opts PaginateOptions) diff.Result {
	if opts.PageSize <= 0 {
		return result
	}

	page := opts.Page
	if page < 1 {
		page = 1
	}

	return diff.Result{
		Added:   paginateMap(result.Added, page, opts.PageSize),
		Removed: paginateMap(result.Removed, page, opts.PageSize),
		Same:    paginateMap(result.Same, page, opts.PageSize),
		Changed: paginateChanged(result.Changed, page, opts.PageSize),
	}
}

func paginateMap(m map[string]string, page, size int) map[string]string {
	if len(m) == 0 {
		return m
	}
	keys := SortedKeys(m, "asc")
	sliced := pageSlice(keys, page, size)
	out := make(map[string]string, len(sliced))
	for _, k := range sliced {
		out[k] = m[k]
	}
	return out
}

func paginateChanged(m map[string]diff.ChangedValue, page, size int) map[string]diff.ChangedValue {
	if len(m) == 0 {
		return m
	}
	keys := SortedChangedKeys(m, "asc")
	sliced := pageSlice(keys, page, size)
	out := make(map[string]diff.ChangedValue, len(sliced))
	for _, k := range sliced {
		out[k] = m[k]
	}
	return out
}

// pageSlice returns the sub-slice of keys for the given 1-based page and size.
func pageSlice(keys []string, page, size int) []string {
	start := (page - 1) * size
	if start >= len(keys) {
		return nil
	}
	end := start + size
	if end > len(keys) {
		end = len(keys)
	}
	return keys[start:end]
}
