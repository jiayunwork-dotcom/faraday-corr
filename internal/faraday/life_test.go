package faraday

import (
	"testing"
)

func TestYearsToPenetrate(t *testing.T) {
	cr := CorrosionRate(feM, feN, refI, feRho)
	got := YearsToPenetrate(1.0, cr)
	want := 1.0 / cr
	if !almostEqual(got, want) {
		t.Errorf("YearsToPenetrate = %v, want %v", got, want)
	}
	if !almostEqual(YearsToPenetrate(0.5, cr), 0.5/cr) {
		t.Errorf("allowance halved should halve the time")
	}
}

func TestYearsToPenetrateZeroRate(t *testing.T) {
	got := YearsToPenetrate(1.0, 0)
	if !(got > 1e300) {
		t.Errorf("zero CR should give an infinite horizon, got %v", got)
	}
}

func TestRemainingThickness(t *testing.T) {
	cr := CorrosionRate(feM, feN, refI, feRho)
	got := RemainingThicknessAfterYears(10, cr, 3)
	want := 10 - 3*cr
	if !almostEqual(got, want) {
		t.Errorf("RemainingThicknessAfterYears = %v, want %v", got, want)
	}
	if RemainingThicknessAfterYears(0.1, cr, 100) != 0 {
		t.Errorf("wall must clamp at zero, got %v", RemainingThicknessAfterYears(0.1, cr, 100))
	}
}

func TestPenetrationAfterYears(t *testing.T) {
	cr := CorrosionRate(feM, feN, refI, feRho)
	if !almostEqual(PenetrationAfterYears(cr, 2), 2*cr) {
		t.Errorf("PenetrationAfterYears = %v, want %v", PenetrationAfterYears(cr, 2), 2*cr)
	}
}

func TestBuildSchedule(t *testing.T) {
	s, err := BuildSchedule(validInput, []float64{1, 2, 4})
	if err != nil {
		t.Fatalf("BuildSchedule failed: %v", err)
	}
	if len(s.Years) != 3 || len(s.CumLossesG) != 3 {
		t.Fatalf("schedule size wrong: %+v", s)
	}
	base := s.CumLossesG[0]
	if !almostEqual(s.CumLossesG[1], 2*base) {
		t.Errorf("2y cumulative = %v, want %v (double of 1y)", s.CumLossesG[1], 2*base)
	}
	if !almostEqual(s.CumLossesG[2], 4*base) {
		t.Errorf("4y cumulative = %v, want %v", s.CumLossesG[2], 4*base)
	}
	if !almostEqual(s.PerAreaG[0], AnnualMassLossPerArea(feM, feN, refI)) {
		t.Errorf("per-area 1y = %v", s.PerAreaG[0])
	}
}
