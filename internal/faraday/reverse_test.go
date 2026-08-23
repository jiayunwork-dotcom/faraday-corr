package faraday

import (
	"testing"
)

// TestReverseFromCorrosionRate checks that the current density solved
// from a target CR reproduces the forward value.
func TestReverseFromCorrosionRate(t *testing.T) {
	base, err := Compute(validInput)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	out, err := ReverseCurrentDensity(ReverseInput{
		MolarMass: feM,
		Valence:   feN,
		Density:   feRho,
		Target:    TargetCorrosionRate,
		CR:        base.CorrosionRate,
	})
	if err != nil {
		t.Fatalf("ReverseCurrentDensity failed: %v", err)
	}
	if !almostEqual(out.ICorr, refI) {
		t.Errorf("solved i_corr = %v, want %v", out.ICorr, refI)
	}
	if !almostEqual(out.CorrosionRate, base.CorrosionRate) {
		t.Errorf("recomputed CR = %v, want %v", out.CorrosionRate, base.CorrosionRate)
	}
}

// TestReverseFromAnnualLoss checks the mass-loss path inversion.
func TestReverseFromAnnualLoss(t *testing.T) {
	annual := AnnualMassLossPerArea(feM, feN, refI)
	out, err := ReverseCurrentDensity(ReverseInput{
		MolarMass: feM,
		Valence:   feN,
		Density:   feRho,
		Target:    TargetAnnualMassLoss,
		Annual:    annual,
	})
	if err != nil {
		t.Fatalf("ReverseCurrentDensity failed: %v", err)
	}
	if !almostEqual(out.ICorr, refI) {
		t.Errorf("solved i_corr = %v, want %v", out.ICorr, refI)
	}
}

// TestReverseRejectsBadTargets pins validation of the reverse path.
func TestReverseRejectsBadTargets(t *testing.T) {
	cases := []ReverseInput{
		{MolarMass: 0, Valence: feN, Density: feRho, Target: TargetCorrosionRate, CR: 1},
		{MolarMass: feM, Valence: 0, Density: feRho, Target: TargetCorrosionRate, CR: 1},
		{MolarMass: feM, Valence: feN, Density: 0, Target: TargetCorrosionRate, CR: 1},
		{MolarMass: feM, Valence: feN, Density: feRho, Target: TargetCorrosionRate, CR: -1},
		{MolarMass: feM, Valence: feN, Density: feRho, Target: TargetAnnualMassLoss, Annual: 0},
		{MolarMass: feM, Valence: feN, Density: feRho, Target: TargetNone},
	}
	for i, in := range cases {
		if _, err := ReverseCurrentDensity(in); err == nil {
			t.Errorf("case %d must be rejected: %+v", i, in)
		}
	}
}

// TestRoundTrip pins the self-consistency of forward then reverse.
func TestRoundTrip(t *testing.T) {
	out, err := RoundTrip(validInput)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}
	if !almostEqual(out.ICorr, refI) {
		t.Errorf("round-trip i_corr = %v, want %v", out.ICorr, refI)
	}
}
