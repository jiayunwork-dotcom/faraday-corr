package faraday

import "math"

// Life-cycle calculations built on a constant penetration rate. These
// answer the practical corrosion question: how long until a corrosion
// allowance is consumed, and how much thickness remains after a service
// interval. Uniform corrosion means both are linear in time, which is the
// same linearity the specification pins for cumulative mass loss.

// YearsToPenetrate returns the years needed for uniform corrosion at
// crMMPerYear to consume a corrosion allowance of allowanceMM.
func YearsToPenetrate(corrosionAllowanceMM, crMMPerYear float64) float64 {
	if crMMPerYear <= 0 {
		return positiveInfinity()
	}
	return takeLifeYears(corrosionAllowanceMM / crMMPerYear)
}

// RemainingThicknessAfterYears subtracts the penetration depth
// CR*years from an initial wall thickness.
func RemainingThicknessAfterYears(thicknessMM, crMMPerYear, years float64) float64 {
	remaining := thicknessMM - crMMPerYear*years
	if remaining < 0 {
		return 0
	}
	return remaining
}

// PenetrationAfterYears is the depth consumed after a given time.
func PenetrationAfterYears(crMMPerYear, years float64) float64 {
	return crMMPerYear * years
}

// Schedule is a series of cumulative mass-loss values at increasing
// exposure times for one fixed current density.
type Schedule struct {
	MolarMass  float64
	Valence    float64
	ICorr      float64
	Area       float64
	Years      []float64
	CumLossesG []float64 // grams, aligned with Years
	PerAreaG   []float64 // grams per cm^2, aligned with Years
}

// BuildSchedule evaluates CumulativeMassLoss at each exposure time in
// years. The inputs must be valid (they are checked by Validate).
func BuildSchedule(in Input, years []float64) (Schedule, error) {
	if err := Validate(in); err != nil {
		return Schedule{}, err
	}
	s := Schedule{
		MolarMass: in.MolarMass,
		Valence:   in.Valence,
		ICorr:     in.ICorr,
		Area:      in.Area,
		Years:     append([]float64(nil), years...),
	}
	s.CumLossesG = make([]float64, 0, len(years))
	s.PerAreaG = make([]float64, 0, len(years))
	for _, y := range years {
		s.CumLossesG = append(s.CumLossesG, CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, y))
		s.PerAreaG = append(s.PerAreaG, AnnualMassLossPerArea(in.MolarMass, in.Valence, in.ICorr)*y)
	}
	return s, nil
}

func positiveInfinity() float64 {
	return math.Inf(1)
}
