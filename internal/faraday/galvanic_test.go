package faraday

import "testing"

func TestPittingShortensLife(t *testing.T) {
	if err := PittingShortensLife(8, 0.1, 4); err != nil {
		t.Fatal(err)
	}
	uni, err := YearsToPerforate(8, 0.1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if uni != 80 {
		t.Fatalf("uniform life %g want 80", uni)
	}
}

func TestCathodeAreaRaisesCurrent(t *testing.T) {
	c := Couple{
		ICorrA: 10, ICorrC: 1,
		AreaA: 1, AreaC: 1,
		BetaA: 0.06, BetaC: 0.12,
	}
	if err := CathodeAreaRaisesCurrent(c, 4); err != nil {
		t.Fatal(err)
	}
}

func TestWagnerNumberUniform(t *testing.T) {
	we, err := WagnerNumber(4, 100, 0.01)
	if err != nil {
		t.Fatal(err)
	}
	if !UniformWhenWagnerLarge(we, 10) {
		t.Fatalf("We=%g should be large", we)
	}
	small, err := WagnerNumber(0.01, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if UniformWhenWagnerLarge(small, 10) {
		t.Fatalf("We=%g should be small", small)
	}
}

func TestRemainingWithPitting(t *testing.T) {
	left, err := RemainingWithPitting(8, 0.1, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	if left != 6 {
		t.Fatalf("remaining %g want 6", left)
	}
}
