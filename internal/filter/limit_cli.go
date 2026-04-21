package filter

// LimitFlags holds the parsed CLI flag values for result limiting.
// It is intended to be populated by the cmd layer and passed into ApplyLimit.
type LimitFlags struct {
	MaxAdded   int
	MaxRemoved int
	MaxChanged int
	MaxSame    int
}

// ToOptions converts LimitFlags into a LimitOptions value suitable for
// passing to ApplyLimit.
func (f LimitFlags) ToOptions() LimitOptions {
	return LimitOptions{
		MaxAdded:   f.MaxAdded,
		MaxRemoved: f.MaxRemoved,
		MaxChanged: f.MaxChanged,
		MaxSame:    f.MaxSame,
	}
}

// IsActive returns true if any limit has been set (i.e. is greater than zero).
func (f LimitFlags) IsActive() bool {
	return f.MaxAdded > 0 || f.MaxRemoved > 0 || f.MaxChanged > 0 || f.MaxSame > 0
}
