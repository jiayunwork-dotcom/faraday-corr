package tafel

import (
	"math"
	"testing"

	"faraday-corr/internal/faraday"
)

func typical() Slopes {
	return Slopes{AnodicVDec: 0.06, CathodicVDec: 0.12}
}

func TestSternGearyRecoversICorr(t *testing.T) {
	s := typical()
	want := 10.0
	rp, err := RpFromICorr(want, s)
	if err != nil {
		t.Fatalf("RpFromICorr: %v", err)
	}
	got, err := ICorrFromRp(rp, s)
	if err != nil {
		t.Fatalf("ICorrFromRp: %v", err)
	}
	if math.Abs(got-want)/want > 1e-9 {
		t.Errorf("round-trip i_corr = %v, want %v", got, want)
	}
	b, err := SternGearyB(s)
	if err != nil {
		t.Fatalf("B: %v", err)
	}
	if b <= 0 || b >= s.AnodicVDec {
		t.Errorf("Stern-Geary B=%v out of expected range", b)
	}
}

func TestTafelHighOverpotentialMatchesBV(t *testing.T) {
	s := typical()
	temp := 298.15
	i0 := 10.0
	pts, err := AnodicSweep(i0, s, temp, 0.12, 0.20, 8)
	if err != nil {
		t.Fatalf("AnodicSweep: %v", err)
	}
	for _, p := range pts {
		rel := math.Abs(p.Current-p.TafelEst) / math.Abs(p.TafelEst)
		if rel > 0.05 {
			t.Errorf("η=%v BV=%v Tafel=%v rel=%v, want Tafel limit", p.EtaV, p.Current, p.TafelEst, rel)
		}
		if p.Current <= 0 {
			t.Errorf("anodic current must be positive at η=%v", p.EtaV)
		}
	}
}

func TestButlerVolmerZeroAtEquilibrium(t *testing.T) {
	i, err := ButlerVolmer(10, 0, 0.5, 0.5, 298.15)
	if err != nil {
		t.Fatalf("BV: %v", err)
	}
	if math.Abs(i) > 1e-12 {
		t.Errorf("i(η=0)=%v, want 0", i)
	}
}

func TestSternGearyThenFaradayRoundTrip(t *testing.T) {
	s := typical()
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1, DurationY: 1}
	base, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	rp, err := RpFromICorr(in.ICorr, s)
	if err != nil {
		t.Fatalf("Rp: %v", err)
	}
	fromRp, err := FaradayFromRp(in.WithCurrent(1), rp, s)
	if err != nil {
		t.Fatalf("FaradayFromRp: %v", err)
	}
	if math.Abs(fromRp.CorrosionRate-base.CorrosionRate)/base.CorrosionRate > 1e-9 {
		t.Errorf("CR from Rp %v vs direct %v", fromRp.CorrosionRate, base.CorrosionRate)
	}
	rev, err := faraday.ReverseCurrentDensity(faraday.ReverseInput{
		MolarMass: in.MolarMass,
		Valence:   in.Valence,
		Density:   in.Density,
		Target:    faraday.TargetCorrosionRate,
		CR:        fromRp.CorrosionRate,
	})
	if err != nil {
		t.Fatalf("Reverse: %v", err)
	}
	if math.Abs(rev.ICorr-in.ICorr)/in.ICorr > 1e-9 {
		t.Errorf("reverse i_corr %v, want %v", rev.ICorr, in.ICorr)
	}
}

func TestBetaFromAlphaRoundTrip(t *testing.T) {
	temp := 298.15
	alpha := 0.5
	beta, err := BetaFromAlpha(alpha, temp)
	if err != nil {
		t.Fatalf("BetaFromAlpha: %v", err)
	}
	back, err := alphaFromBeta(beta, temp)
	if err != nil {
		t.Fatalf("alphaFromBeta: %v", err)
	}
	if math.Abs(back-alpha) > 1e-12 {
		t.Errorf("alpha round-trip %v, want %v", back, alpha)
	}
}
