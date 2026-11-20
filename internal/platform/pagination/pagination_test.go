package pagination

import "testing"

func TestParse_DefaultsAndCaps(t *testing.T) {
	params, err := Parse("", "", 50, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params.Page != 1 || params.Limit != 50 || params.Skip != 0 {
		t.Fatalf("unexpected params: %+v", params)
	}
}

func TestParse_ComputesSkip(t *testing.T) {
	params, err := Parse("3", "25", 50, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params.Skip != 50 {
		t.Fatalf("expected skip 50, got %d", params.Skip)
	}
}

func TestParse_RejectsOverMax(t *testing.T) {
	_, err := Parse("1", "500", 50, 200)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
