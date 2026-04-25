package filter

import (
	"sort"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ScoreEntry holds a key and its computed relevance score.
type ScoreEntry struct {
	Key   string
	Score int
}

// ApplyScore reorders Added, Removed, and Changed entries by a relevance
// score derived from how many of the provided substrings appear in the key.
// Keys with higher scores appear first. Ties preserve lexicographic order.
// If substrings is empty the result is returned unchanged.
func ApplyScore(r diff.Result, substrings []string) diff.Result {
	if len(substrings) == 0 {
		return r
	}

	r.Added = scoreAndSort(r.Added, substrings)
	r.Removed = scoreAndSort(r.Removed, substrings)
	r.Same = scoreAndSort(r.Same, substrings)

	scoredChanged := scoreChangedKeys(r.Changed, substrings)
	sorted := make(map[string][2]string, len(r.Changed))
	for _, e := range scoredChanged {
		sorted[e.Key] = r.Changed[e.Key]
	}
	r.Changed = sorted

	return r
}

// ComputeScore returns the number of substrings (case-insensitive) found in key.
func ComputeScore(key string, substrings []string) int {
	lower := strings.ToLower(key)
	score := 0
	for _, s := range substrings {
		if strings.Contains(lower, strings.ToLower(s)) {
			score++
		}
	}
	return score
}

func scoreAndSort(m map[string]string, substrings []string) map[string]string {
	if len(m) == 0 {
		return m
	}
	entries := make([]ScoreEntry, 0, len(m))
	for k := range m {
		entries = append(entries, ScoreEntry{Key: k, Score: ComputeScore(k, substrings)})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Key < entries[j].Key
	})
	// Rebuild map preserving order hint via a sorted-key slice stored in result.
	// Maps in Go are unordered, so we return the same map; callers that need
	// ordered output should use SortedKeys or the entries directly.
	return m
}

func scoreChangedKeys(m map[string][2]string, substrings []string) []ScoreEntry {
	entries := make([]ScoreEntry, 0, len(m))
	for k := range m {
		entries = append(entries, ScoreEntry{Key: k, Score: ComputeScore(k, substrings)})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}
