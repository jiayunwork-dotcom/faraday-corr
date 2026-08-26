package faraday

import "testing"

func TestObservedBelowBoth(t *testing.T) {
	if err := ObservedBelowBoth(10, 4); err != nil {
		t.Fatal(err)
	}
}

func TestThinFilmRaisesLimit(t *testing.T) {
	if err := ThinFilmRaisesLimit(2, 1e-9, 10, 1e-4, 4e-4); err != nil {
		t.Fatal(err)
	}
}
