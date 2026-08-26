package faraday

import "testing"

func TestVerifyUnitChainCloses(t *testing.T) {
	c, err := VerifyUnitChain(validInput)
	if err != nil {
		t.Fatalf("VerifyUnitChain failed: %v", err)
	}
	if !c.Consistent {
		t.Errorf("unit chain not consistent: mass path %v vs direct %v", c.CRFromMass, c.CRDirect)
	}
	if !almostEqual(c.KEffective, K) {
		t.Errorf("K effective = %v, want pinned %v", c.KEffective, K)
	}
	if !almostEqual(c.AnnualPerCm2, AnnualMassLossPerArea(feM, feN, refI)) {
		t.Errorf("annual per cm2 = %v", c.AnnualPerCm2)
	}
	if !almostEqual(c.CurrentDensityA, 1e-5) {
		t.Errorf("current density in A/cm^2 = %v, want 1e-5", c.CurrentDensityA)
	}
}

func TestVerifyUnitChainRejectsInvalidInput(t *testing.T) {
	if _, err := VerifyUnitChain(validInput.WithDensity(-1)); err == nil {
		t.Error("negative density must fail the chain walk")
	}
}

func TestKDerivedMatchesPinned(t *testing.T) {
	derived, diff := KDerived()
	if !almostEqual(derived, K) {
		t.Errorf("derived K = %v, want %v", derived, K)
	}
	if diff > 1e-15 {
		t.Errorf("K drift = %v, want ~0", diff)
	}
}

func TestUnitsUsedListsCoreUnits(t *testing.T) {
	u := UnitsUsed()
	found := map[string]bool{}
	for _, s := range u {
		found[s] = true
	}
	for _, want := range []string{
		"i_corr: uA/cm^2",
		"M: g/mol",
		"rho: g/cm^3",
		"mdot: g/(cm^2 s)",
		"CR: mm/y",
		"total current I: uA",
	} {
		if !found[want] {
			t.Errorf("unit list missing %q", want)
		}
	}
}
