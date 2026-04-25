package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeGroupResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"APP_HOST": "localhost",
			"DB_HOST":  "db.prod",
		},
		Removed: map[string]string{
			"APP_PORT": "8080",
		},
		Same: map[string]string{
			"DB_NAME": "mydb",
			"NOPFX":   "val",
		},
		Changed: map[string]diff.ChangedValue{
			"APP_DEBUG": {Old: "false", New: "true"},
		},
	}
}

func TestApplyGroup_EmptyDelimiter(t *testing.T) {
	r := makeGroupResult()
	gr := ApplyGroup(r, "")
	if len(gr.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(gr.Groups))
	}
	if _, ok := gr.Groups[""]; !ok {
		t.Fatal("expected empty-string group")
	}
}

func TestApplyGroup_UnderscoreDelimiter(t *testing.T) {
	r := makeGroupResult()
	gr := ApplyGroup(r, "_")

	if _, ok := gr.Groups["APP"]; !ok {
		t.Fatal("expected APP group")
	}
	if _, ok := gr.Groups["DB"]; !ok {
		t.Fatal("expected DB group")
	}
	// NOPFX has no underscore → empty group
	if _, ok := gr.Groups[""]; !ok {
		t.Fatal("expected empty group for keys without delimiter")
	}
}

func TestApplyGroup_GroupContents(t *testing.T) {
	r := makeGroupResult()
	gr := ApplyGroup(r, "_")

	app := gr.Groups["APP"]
	if app == nil {
		t.Fatal("APP group is nil")
	}
	if _, ok := app.Added["APP_HOST"]; !ok {
		t.Error("APP_HOST should be in APP.Added")
	}
	if _, ok := app.Removed["APP_PORT"]; !ok {
		t.Error("APP_PORT should be in APP.Removed")
	}
	if _, ok := app.Changed["APP_DEBUG"]; !ok {
		t.Error("APP_DEBUG should be in APP.Changed")
	}
}

func TestApplyGroup_SortedGroupKeys(t *testing.T) {
	r := makeGroupResult()
	gr := ApplyGroup(r, "_")
	keys := SortedGroupKeys(gr)
	if len(keys) == 0 {
		t.Fatal("expected non-empty keys")
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestApplyGroup_NoKeys(t *testing.T) {
	r := diff.Result{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Same:    map[string]string{},
		Changed: map[string]diff.ChangedValue{},
	}
	gr := ApplyGroup(r, "_")
	if len(gr.Groups) != 0 {
		t.Errorf("expected 0 groups for empty result, got %d", len(gr.Groups))
	}
}
