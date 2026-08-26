package faraday

import (
	"testing"
)

func TestMPYConversion(t *testing.T) {
	if !almostEqual(MPY(0.0254), 1.0) {
		t.Errorf("MPY(0.0254) = %v, want 1", MPY(0.0254))
	}
	if !almostEqual(MMPerYearFromMPY(100), 2.54) {
		t.Errorf("MMPerYearFromMPY(100) = %v, want 2.54", MMPerYearFromMPY(100))
	}
	if !almostEqual(MMPerYearFromMPY(MPY(0.1159)), 0.1159) {
		t.Errorf("round trip failed: %v", MMPerYearFromMPY(MPY(0.1159)))
	}
}

func TestMddConversion(t *testing.T) {
	if !almostEqual(Mdd(0.00365), 1.0) {
		t.Errorf("Mdd(0.00365) = %v, want 1", Mdd(0.00365))
	}
	if !almostEqual(GPerCm2YearFromMdd(1.0), 0.00365) {
		t.Errorf("GPerCm2YearFromMdd(1) = %v, want 0.00365", GPerCm2YearFromMdd(1))
	}
}

func TestGPerM2Year(t *testing.T) {
	if !almostEqual(GPerM2Year(1.0), 1e4) {
		t.Errorf("GPerM2Year(1) = %v, want 1e4", GPerM2Year(1))
	}
}

func TestDensityAndAreaConversions(t *testing.T) {
	if !almostEqual(KGPerm3(7.874), 7874) {
		t.Errorf("KGPerm3(7.874) = %v, want 7874", KGPerm3(7.874))
	}
	if !almostEqual(Cm2ToM2(1e4), 1.0) {
		t.Errorf("Cm2ToM2(1e4) = %v, want 1", Cm2ToM2(1e4))
	}
	if !almostEqual(M2ToCm2(1.0), 1e4) {
		t.Errorf("M2ToCm2(1) = %v, want 1e4", M2ToCm2(1))
	}
}

func TestSummarizeDerivesEngineeringUnits(t *testing.T) {
	res, err := Compute(validInput)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	s := Summarize(res)
	if !almostEqual(s.MPY, MPY(res.CorrosionRate)) {
		t.Errorf("s.MPY = %v, want %v", s.MPY, MPY(res.CorrosionRate))
	}
	if !almostEqual(s.MassLossGPerM2Year, res.AnnualMassLoss*1e4) {
		t.Errorf("s.MassLossGPerM2Year = %v", s.MassLossGPerM2Year)
	}
	if !almostEqual(s.DensityKGPerm3, 7874) {
		t.Errorf("s.DensityKGPerm3 = %v, want 7874", s.DensityKGPerm3)
	}
	if !almostEqual(s.AreaM2, 1e-4) {
		t.Errorf("s.AreaM2 = %v, want 1e-4", s.AreaM2)
	}
	if !almostEqual(s.TotalCurrentMilliA, 10e-3) {
		t.Errorf("s.TotalCurrentMilliA = %v, want 0.01", s.TotalCurrentMilliA)
	}
}
