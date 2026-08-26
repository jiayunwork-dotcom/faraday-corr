package faraday

import (
	"math"
	"testing"
)

const (
	feM   = 55.845
	feN   = 2.0
	feRho = 7.874
	refI  = 10.0
)

func almostEqual(got, want float64) bool {
	return math.Abs(got-want) <= 1e-9*math.Max(math.Abs(want), math.Abs(got))
}

func TestFaradayConstantPinned(t *testing.T) {
	want := 96485.33212
	if Faraday != want {
		t.Errorf("Faraday = %v, want %v (CODATA)", Faraday, want)
	}
	if Faraday == 965 || Faraday == 96500 {
		t.Errorf("Faraday must not be an approximation like 965, got %v", Faraday)
	}
}

func TestKConstantPinned(t *testing.T) {
	want := 315.36 / Faraday
	if !almostEqual(K, want) {
		t.Errorf("K = %v, want %v (10*1e-6*yearSeconds/F)", K, want)
	}
}

func TestMassLossRateFormula(t *testing.T) {
	got := MassLossRate(feM, feN, refI)
	want := feM * (refI * 1e-6) / (feN * Faraday)
	if !almostEqual(got, want) {
		t.Errorf("MassLossRate = %v, want %v", got, want)
	}
}

func TestAnnualMassLossFromRate(t *testing.T) {
	mdot := MassLossRate(feM, feN, refI)
	got := AnnualMassLossPerAreaFromRate(mdot)
	want := mdot * yearSeconds
	if !almostEqual(got, want) {
		t.Errorf("AnnualMassLossPerArea = %v, want %v", got, want)
	}
}

func TestCorrosionRateFormula(t *testing.T) {
	got := CorrosionRate(feM, feN, refI, feRho)
	want := K * (feM / feN) * refI / feRho
	if !almostEqual(got, want) {
		t.Errorf("CorrosionRate = %v, want %v", got, want)
	}
	if got < 0.10 || got > 0.13 {
		t.Errorf("CorrosionRate = %v, want about 0.116 mm/y for iron in seawater", got)
	}
}

func TestCorrosionRateMatchesMassLossChain(t *testing.T) {
	annual := AnnualMassLossPerArea(feM, feN, refI)
	cr := CorrosionRate(feM, feN, refI, feRho)
	want := milliPerCm * annual / feRho
	if !almostEqual(cr, want) {
		t.Errorf("CorrosionRate = %v, want %v from mass-loss chain", cr, want)
	}
}

func TestEquivalentWeight(t *testing.T) {
	got := EquivalentWeight(feM, feN)
	if !almostEqual(got, feM/feN) {
		t.Errorf("EquivalentWeight = %v, want %v", got, feM/feN)
	}
}

func TestCorrosionRateUmY(t *testing.T) {
	got := CorrosionRateUmY(0.1159)
	if !almostEqual(got, 115.9) {
		t.Errorf("CorrosionRateUmY = %v, want 115.9", got)
	}
}

func TestTotalCurrent(t *testing.T) {
	got := TotalCurrent(refI, 2.5)
	if !almostEqual(got, 25.0) {
		t.Errorf("TotalCurrent = %v, want 25", got)
	}
	gotA := TotalCurrentAmps(refI, 2.5)
	if !almostEqual(gotA, 25e-6) {
		t.Errorf("TotalCurrentAmps = %v, want 2.5e-5", gotA)
	}
}

func TestCumulativeMassLoss(t *testing.T) {
	got := CumulativeMassLoss(feM, feN, refI, 2.0, 3.0)
	want := AnnualMassLossPerArea(feM, feN, refI) * 2.0 * 3.0
	if !almostEqual(got, want) {
		t.Errorf("CumulativeMassLoss = %v, want %v", got, want)
	}
}

func TestComputePopulatesAllFields(t *testing.T) {
	in := Input{ICorr: refI, MolarMass: feM, Valence: feN, Density: feRho, Area: 2.0, DurationY: 1.0}
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute returned error: %v", err)
	}
	if res.Input != in {
		t.Errorf("Compute did not echo input: got %+v, want %+v", res.Input, in)
	}
	if !almostEqual(res.MassLossRate, MassLossRate(feM, feN, refI)) {
		t.Errorf("res.MassLossRate = %v", res.MassLossRate)
	}
	if !almostEqual(res.CorrosionRate, CorrosionRate(feM, feN, refI, feRho)) {
		t.Errorf("res.CorrosionRate = %v", res.CorrosionRate)
	}
	if !almostEqual(res.CorrosionRateUmY, 1000*res.CorrosionRate) {
		t.Errorf("res.CorrosionRateUmY = %v, want %v", res.CorrosionRateUmY, 1000*res.CorrosionRate)
	}
	if !almostEqual(res.TotalCurrent, 20.0) {
		t.Errorf("res.TotalCurrent = %v, want 20", res.TotalCurrent)
	}
	if !almostEqual(res.CumulativeMassLoss, 2.0*res.AnnualMassLoss) {
		t.Errorf("res.CumulativeMassLoss = %v, want %v", res.CumulativeMassLoss, 2.0*res.AnnualMassLoss)
	}
}

func TestZeroCurrentYieldsZeroRates(t *testing.T) {
	in := Input{ICorr: 0, MolarMass: feM, Valence: feN, Density: feRho, Area: 1.0}
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("zero current must be valid, got error: %v", err)
	}
	if res.CorrosionRate != 0 {
		t.Errorf("CR at i=0 = %v, want 0", res.CorrosionRate)
	}
	if res.MassLossRate != 0 {
		t.Errorf("mdot at i=0 = %v, want 0", res.MassLossRate)
	}
}
