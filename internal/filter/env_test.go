package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeEnvResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "new_val"},
		Removed: map[string]string{"OLD_KEY": "old_val"},
		Changed: map[string]diff.ChangedValue{
			"CHANGED_KEY": {Old: "v1", New: "v2"},
		},
		Same: map[string]string{"SAME_KEY": "same_val"},
	}
}

func TestApplyEnvOverride_NoEnvSet(t *testing.T) {
	result := makeEnvResult()
	out := filter.ApplyEnvOverride(result)

	if out.Added["NEW_KEY"] != "new_val" {
		t.Errorf("expected new_val, got %s", out.Added["NEW_KEY"])
	}
	if out.Removed["OLD_KEY"] != "old_val" {
		t.Errorf("expected old_val, got %s", out.Removed["OLD_KEY"])
	}
	if out.Changed["CHANGED_KEY"].New != "v2" {
		t.Errorf("expected v2, got %s", out.Changed["CHANGED_KEY"].New)
	}
	if out.Same["SAME_KEY"] != "same_val" {
		t.Errorf("expected same_val, got %s", out.Same["SAME_KEY"])
	}
}

func TestApplyEnvOverride_OverridesAdded(t *testing.T) {
	t.Setenv("NEW_KEY", "overridden")
	result := makeEnvResult()
	out := filter.ApplyEnvOverride(result)

	if out.Added["NEW_KEY"] != "overridden" {
		t.Errorf("expected overridden, got %s", out.Added["NEW_KEY"])
	}
}

func TestApplyEnvOverride_OverridesChanged(t *testing.T) {
	t.Setenv("CHANGED_KEY", "env_override")
	result := makeEnvResult()
	out := filter.ApplyEnvOverride(result)

	cv := out.Changed["CHANGED_KEY"]
	if cv.New != "env_override" {
		t.Errorf("expected env_override, got %s", cv.New)
	}
	if cv.Old != "v1" {
		t.Errorf("expected old value v1 preserved, got %s", cv.Old)
	}
}

func TestApplyEnvOverride_OverridesSame(t *testing.T) {
	t.Setenv("SAME_KEY", "new_same")
	result := makeEnvResult()
	out := filter.ApplyEnvOverride(result)

	if out.Same["SAME_KEY"] != "new_same" {
		t.Errorf("expected new_same, got %s", out.Same["SAME_KEY"])
	}
}

func TestEnvVarNames_EmptyPrefix(t *testing.T) {
	t.Setenv("ENVOY_TEST_VAR", "1")
	names := filter.EnvVarNames("")
	if len(names) == 0 {
		t.Error("expected at least one env var name")
	}
}

func TestEnvVarNames_WithPrefix(t *testing.T) {
	t.Setenv("ENVOY_DIFF_UNIQUE_XYZ", "1")
	names := filter.EnvVarNames("ENVOY_DIFF_UNIQUE")
	if len(names) != 1 {
		t.Errorf("expected 1 match, got %d", len(names))
	}
	if names[0] != "ENVOY_DIFF_UNIQUE_XYZ" {
		t.Errorf("unexpected name: %s", names[0])
	}
}
