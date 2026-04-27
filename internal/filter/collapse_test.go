package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeCollapseResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"PATH": "/usr/bin:/usr/bin:/usr/local/bin",
			"NAME": "alice",
		},
		Removed: map[string]string{
			"OLD_PATH": "/bin:/bin:/sbin",
		},
		Same: map[string]string{
			"HOME": "/home/user",
		},
		Changed: map[string]diff.ChangedValue{
			"CLASSPATH": {Old: "a:b:a", New: "c:d:c"},
		},
	}
}

func TestApplyCollapse_NoDelimiter(t *testing.T) {
	r := makeCollapseResult()
	out := ApplyCollapse(r, CollapseConfig{})
	if out.Added["PATH"] != "/usr/bin:/usr/bin:/usr/local/bin" {
		t.Errorf("expected unchanged PATH, got %q", out.Added["PATH"])
	}
}

func TestApplyCollapse_DeduplicatesAdded(t *testing.T) {
	r := makeCollapseResult()
	out := ApplyCollapse(r, CollapseConfig{Delimiter: ":"})
	got := out.Added["PATH"]
	if got != "/usr/bin:/usr/local/bin" {
		t.Errorf("expected deduped PATH, got %q", got)
	}
}

func TestApplyCollapse_PreservesNonDuplicates(t *testing.T) {
	r := makeCollapseResult()
	out := ApplyCollapse(r, CollapseConfig{Delimiter: ":"})
	if out.Same["HOME"] != "/home/user" {
		t.Errorf("expected HOME unchanged, got %q", out.Same["HOME"])
	}
	if out.Added["NAME"] != "alice" {
		t.Errorf("expected NAME unchanged, got %q", out.Added["NAME"])
	}
}

func TestApplyCollapse_DeduplicatesChanged(t *testing.T) {
	r := makeCollapseResult()
	out := ApplyCollapse(r, CollapseConfig{Delimiter: ":"})
	cv := out.Changed["CLASSPATH"]
	if cv.Old != "a:b" {
		t.Errorf("expected old=a:b, got %q", cv.Old)
	}
	if cv.New != "c:d" {
		t.Errorf("expected new=c:d, got %q", cv.New)
	}
}

func TestApplyCollapse_KeyFilter(t *testing.T) {
	r := makeCollapseResult()
	// Only collapse keys containing "PATH"
	out := ApplyCollapse(r, CollapseConfig{Delimiter: ":", Keys: []string{"path"}})
	if out.Added["PATH"] != "/usr/bin:/usr/local/bin" {
		t.Errorf("expected deduped PATH, got %q", out.Added["PATH"])
	}
	// CLASSPATH also matches "path"
	if out.Changed["CLASSPATH"].Old != "a:b" {
		t.Errorf("expected deduped CLASSPATH old, got %q", out.Changed["CLASSPATH"].Old)
	}
}
