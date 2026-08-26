package coating

import (
	"math"
	"testing"

	"faraday-corr/internal/faraday"
)

func TestHolidayAreaScalesTotalCurrentNotLocalCR(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 100, DurationY: 1}
	h := Holiday{GeometricAreaCm2: 100, Fraction: 0.02}
	out, err := Apply(in, h)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	bare, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("bare: %v", err)
	}
	if math.Abs(out.LocalCR-bare.CorrosionRate)/bare.CorrosionRate > 1e-12 {
		t.Errorf("local CR %v changed from bare %v", out.LocalCR, bare.CorrosionRate)
	}
	if math.Abs(out.TotalCurrentUA-bare.TotalCurrent*0.02)/out.TotalCurrentUA > 1e-9 {
		t.Errorf("total I %v, want 2%% of bare %v", out.TotalCurrentUA, bare.TotalCurrent)
	}
	if math.Abs(out.ExposedAreaCm2-2) > 1e-12 {
		t.Errorf("exposed area %v, want 2", out.ExposedAreaCm2)
	}
}

func TestCoatingHolidayMassLossNotGeometricArea(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, DurationY: 2}
	h := Holiday{GeometricAreaCm2: 50, Fraction: 0.1}
	out, err := Apply(in, h)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !(out.TotalMassLossG < out.BareMassLossG) {
		t.Errorf("holiday mass %v not below fully-bare %v", out.TotalMassLossG, out.BareMassLossG)
	}
	ratio := out.TotalMassLossG / out.BareMassLossG
	if math.Abs(ratio-0.1) > 1e-9 {
		t.Errorf("mass ratio %v, want holiday fraction 0.1", ratio)
	}
	r, err := CurrentRatio(h)
	if err != nil {
		t.Fatalf("CurrentRatio: %v", err)
	}
	if math.Abs(r-0.1) > 1e-12 {
		t.Errorf("current ratio %v, want 0.1", r)
	}
}

func TestRejectsClosedCoating(t *testing.T) {
	h := Holiday{GeometricAreaCm2: 10, Fraction: 0}
	if _, err := Apply(faraday.Input{ICorr: 1, MolarMass: 55, Valence: 2, Density: 7}, h); err == nil {
		t.Fatal("fraction=0: expected error")
	}
	h.Fraction = 1.5
	if _, err := h.ExposedArea(); err == nil {
		t.Fatal("fraction>1: expected error")
	}
}

func TestIntactFractionOneMatchesBare(t *testing.T) {
	in := faraday.Input{ICorr: 5, MolarMass: 26.982, Valence: 3, Density: 2.70, DurationY: 1}
	h := Holiday{GeometricAreaCm2: 8, Fraction: 1}
	out, err := Apply(in, h)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if math.Abs(out.TotalMassLossG-out.BareMassLossG) > 1e-12 {
		t.Errorf("fraction=1 mass %v vs bare %v", out.TotalMassLossG, out.BareMassLossG)
	}
}
