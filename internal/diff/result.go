package diff

// Result holds the categorised diff between two env var maps.
type Result struct {
	// Added contains keys present in the second (production) map but not the first.
	Added map[string]string

	// Removed contains keys present in the first (staging) map but not the second.
	Removed map[string]string

	// Changed contains keys present in both maps but with differing values.
	// The array holds [stagingValue, productionValue].
	Changed map[string][2]string

	// Unchanged contains keys present in both maps with identical values.
	Unchanged map[string]string
}

// HasDiff returns true if there are any differences between the two env sets.
func (r Result) HasDiff() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// TotalKeys returns the total number of unique keys across all categories.
func (r Result) TotalKeys() int {
	return len(r.Added) + len(r.Removed) + len(r.Changed) + len(r.Unchanged)
}
