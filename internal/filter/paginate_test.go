package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makePaginateResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A_KEY": "1", "B_KEY": "2", "C_KEY": "3"},
		Removed: map[string]string{"X_KEY": "9", "Y_KEY": "8"},
		Same:    map[string]string{"S_KEY": "s"},
		Changed: map[string]diff.ChangedValue{
			"CH_ONE": {Old: "a", New: "b"},
			"CH_TWO": {Old: "c", New: "d"},
			"CH_THREE": {Old: "e", New: "f"},
		},
	}
}

func TestApplyPaginate_NoPageSize(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 1, PageSize: 0})
	if len(out.Added) != 3 {
		t.Errorf("expected 3 added, got %d", len(out.Added))
	}
}

func TestApplyPaginate_FirstPage(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 1, PageSize: 2})
	if len(out.Added) != 2 {
		t.Errorf("expected 2 added on page 1, got %d", len(out.Added))
	}
	// sorted ascending: A_KEY, B_KEY should be on page 1
	if _, ok := out.Added["A_KEY"]; !ok {
		t.Error("expected A_KEY on page 1")
	}
	if _, ok := out.Added["B_KEY"]; !ok {
		t.Error("expected B_KEY on page 1")
	}
}

func TestApplyPaginate_SecondPage(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 2, PageSize: 2})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added on page 2, got %d", len(out.Added))
	}
	if _, ok := out.Added["C_KEY"]; !ok {
		t.Error("expected C_KEY on page 2")
	}
}

func TestApplyPaginate_PageBeyondEnd(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 99, PageSize: 2})
	if len(out.Added) != 0 {
		t.Errorf("expected 0 added on out-of-range page, got %d", len(out.Added))
	}
}

func TestApplyPaginate_ChangedKeys(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 1, PageSize: 2})
	if len(out.Changed) != 2 {
		t.Errorf("expected 2 changed on page 1, got %d", len(out.Changed))
	}
}

func TestApplyPaginate_DefaultPageIsOne(t *testing.T) {
	r := makePaginateResult()
	out := ApplyPaginate(r, PaginateOptions{Page: 0, PageSize: 2})
	if len(out.Added) != 2 {
		t.Errorf("expected 2 added when page defaults to 1, got %d", len(out.Added))
	}
}
