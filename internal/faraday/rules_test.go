package faraday

import (
	"strings"
	"testing"
)

var validInput = Input{
	ICorr:     refI,
	MolarMass: feM,
	Valence:   feN,
	Density:   feRho,
	Area:      1.0,
	DurationY: 1.0,
}

func TestDoubleCurrentDoublesRates(t *testing.T) {
	base, err := Compute(validInput)
	if err != nil {
		t.Fatalf("base compute failed: %v", err)
	}
	doubled, err := Compute(validInput.WithCurrent(refI * 2))
	if err != nil {
		t.Fatalf("doubled compute failed: %v", err)
	}
	assertRatio(t, "mdot(i*2)/mdot", doubled.MassLossRate, base.MassLossRate, 2.0)
	assertRatio(t, "CR(i*2)/CR", doubled.CorrosionRate, base.CorrosionRate, 2.0)
}

func TestDoubleValenceHalvesRates(t *testing.T) {
	base, err := Compute(validInput)
	if err != nil {
		t.Fatalf("base compute failed: %v", err)
	}
	doubled, err := Compute(validInput.WithValence(feN * 2))
	if err != nil {
		t.Fatalf("doubled compute failed: %v", err)
	}
	assertRatio(t, "mdot(n*2)/mdot", doubled.MassLossRate, base.MassLossRate, 0.5)
	assertRatio(t, "CR(n*2)/CR", doubled.CorrosionRate, base.CorrosionRate, 0.5)
}

func TestDoubleDensityLeavesMassUnchanged(t *testing.T) {
	base, err := Compute(validInput)
	if err != nil {
		t.Fatalf("base compute failed: %v", err)
	}
	doubled, err := Compute(validInput.WithDensity(feRho * 2))
	if err != nil {
		t.Fatalf("doubled compute failed: %v", err)
	}
	assertRatio(t, "mdot(rho*2)/mdot", doubled.MassLossRate, base.MassLossRate, 1.0)
	assertRatio(t, "CR(rho*2)/CR", doubled.CorrosionRate, base.CorrosionRate, 0.5)
}

func TestDoubleDurationDoublesCumulativeLoss(t *testing.T) {
	base, err := Compute(validInput)
	if err != nil {
		t.Fatalf("base compute failed: %v", err)
	}
	doubled, err := Compute(validInput.WithDuration(2.0))
	if err != nil {
		t.Fatalf("doubled compute failed: %v", err)
	}
	assertRatio(t, "cumloss(t*2)/cumloss", doubled.CumulativeMassLoss, base.CumulativeMassLoss, 2.0)
}

func TestIronVsAluminumDifferentDepthRates(t *testing.T) {
	alM, alN, alRho := 26.9815, 3.0, 2.70
	feCR := CorrosionRate(feM, feN, refI, feRho)
	alCR := CorrosionRate(alM, alN, refI, alRho)
	if feCR == alCR {
		t.Errorf("Fe CR %v must differ from Al CR %v at the same i_corr", feCR, alCR)
	}
	if alCR <= 0 || feCR <= 0 {
		t.Errorf("both depth rates must be positive, got Fe %v Al %v", feCR, alCR)
	}
}

func TestRejectNegativeICorr(t *testing.T) {
	for _, bad := range []float64{-1, -10.5} {
		if err := Validate(validInput.WithCurrent(bad)); err == nil {
			t.Errorf("i_corr=%v must be rejected", bad)
		}
	}
	if err := Validate(validInput.WithCurrent(0)); err != nil {
		t.Errorf("i_corr=0 must be accepted (i=0 => CR=0), got %v", err)
	}
}

func TestRejectNonPositiveValence(t *testing.T) {
	for _, bad := range []float64{0, -2} {
		if err := Validate(validInput.WithValence(bad)); err == nil {
			t.Errorf("n=%v must be rejected", bad)
		}
	}
}

func TestRejectNonPositiveDensity(t *testing.T) {
	for _, bad := range []float64{0, -7.874} {
		if err := Validate(validInput.WithDensity(bad)); err == nil {
			t.Errorf("rho=%v must be rejected", bad)
		}
	}
}

func TestRejectNonPositiveMolarMass(t *testing.T) {
	for _, bad := range []float64{0, -55.845} {
		if err := Validate(validInput.WithMolarMass(bad)); err == nil {
			t.Errorf("M=%v must be rejected", bad)
		}
	}
}

func TestRejectNegativeAreaAndDuration(t *testing.T) {
	if err := Validate(validInput.WithArea(-1)); err == nil {
		t.Error("negative area must be rejected")
	}
	if err := Validate(validInput.WithDuration(-1)); err == nil {
		t.Error("negative duration must be rejected")
	}
}

func TestValidationErrorMessageNamesTheField(t *testing.T) {
	err := Validate(validInput.WithDensity(-2))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rho") {
		t.Errorf("error %q should name field rho", err.Error())
	}
}

func TestInvariantChecks(t *testing.T) {
	checks, err := InvariantChecks(validInput)
	if err != nil {
		t.Fatalf("InvariantChecks failed: %v", err)
	}
	expect := map[string]float64{
		"i=0 => CR=0":              0,
		"i*2 => mdot*2":            2,
		"i*2 => CR*2":              2,
		"n*2 => mdot/2":            0.5,
		"n*2 => CR/2":              0.5,
		"rho*2 => mdot unchanged":  1,
		"rho*2 => CR/2":            0.5,
		"t*2 => cumulative loss*2": 2,
	}
	for _, c := range checks {
		want, ok := expect[c.Name]
		if !ok {
			t.Errorf("unexpected invariant %q", c.Name)
			continue
		}
		if !almostEqual(c.Ratio, want) {
			t.Errorf("%s ratio = %v, want %v", c.Name, c.Ratio, want)
		}
	}
}

func assertRatio(t *testing.T, label string, got, base, want float64) {
	t.Helper()
	if base == 0 {
		t.Fatalf("%s: base value is zero, cannot form ratio", label)
	}
	if !almostEqual(got/base, want) {
		t.Errorf("%s = %v, want %v", label, got/base, want)
	}
}
