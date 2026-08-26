package arrhenius

import (
	"math"
	"testing"

	"faraday-corr/internal/faraday"
)

func typical() Law {
	return Law{IRefUACm2: 10, TRefK: 298.15, EaJMol: 50000}
}

func TestArrheniusHigherTRaisesICorr(t *testing.T) {
	l := typical()
	cold, err := l.ICorrAt(288.15)
	if err != nil {
		t.Fatalf("288: %v", err)
	}
	hot, err := l.ICorrAt(308.15)
	if err != nil {
		t.Fatalf("308: %v", err)
	}
	if !(hot > l.IRefUACm2) {
		t.Errorf("i(308)=%v not above i_ref=%v", hot, l.IRefUACm2)
	}
	if !(cold < l.IRefUACm2) {
		t.Errorf("i(288)=%v not below i_ref=%v", cold, l.IRefUACm2)
	}
	atRef, err := l.ICorrAt(298.15)
	if err != nil {
		t.Fatalf("ref: %v", err)
	}
	if math.Abs(atRef-10)/10 > 1e-12 {
		t.Errorf("i(Tref)=%v, want 10", atRef)
	}
}

func TestQ10NearTwoForTypicalEa(t *testing.T) {
	l := typical()
	q, err := l.Q10(298.15)
	if err != nil {
		t.Fatalf("Q10: %v", err)
	}
	if q < 1.7 || q > 2.3 {
		t.Errorf("Q10=%v, want ~2 for Ea=50 kJ/mol near 25 C", q)
	}
}

func TestTemperatureThenFaradayVsFaradayThenScale(t *testing.T) {
	l := typical()
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1, DurationY: 1}
	hot, iHot, err := FaradayAtTemperature(in, l, 318.15)
	if err != nil {
		t.Fatalf("hot: %v", err)
	}
	base, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("base: %v", err)
	}
	factor, err := ScaleFactor(l, 298.15, 318.15)
	if err != nil {
		t.Fatalf("scale: %v", err)
	}
	scaledCR := base.CorrosionRate * factor
	if math.Abs(hot.CorrosionRate-scaledCR)/scaledCR > 1e-9 {
		t.Errorf("temperature-then-Faraday CR %v vs Faraday-then-scale %v", hot.CorrosionRate, scaledCR)
	}
	if math.Abs(hot.MassLossRate-base.MassLossRate*factor)/(base.MassLossRate*factor) > 1e-9 {
		t.Errorf("mass rate did not scale with i_corr")
	}
	if math.Abs(iHot-in.ICorr*factor)/iHot > 1e-9 {
		t.Errorf("i_hot %v vs i_ref*factor %v", iHot, in.ICorr*factor)
	}
	if hot.CorrosionRate <= base.CorrosionRate {
		t.Errorf("hot CR %v not above base %v", hot.CorrosionRate, base.CorrosionRate)
	}
}

func TestRejectsNonPositive(t *testing.T) {
	l := typical()
	l.EaJMol = 0
	if _, err := l.ICorrAt(298.15); err == nil {
		t.Fatal("Ea=0: expected error")
	}
	l = typical()
	if _, err := l.ICorrAt(-1); err == nil {
		t.Fatal("T<0: expected error")
	}
}
