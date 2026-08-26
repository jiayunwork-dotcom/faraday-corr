package faraday

import "math"

func YearsToPenetrate(corrosionAllowanceMM, crMMPerYear float64) float64 {
	if crMMPerYear <= 0 {
		return positiveInfinity()
	}
	return corrosionAllowanceMM / crMMPerYear
}

func RemainingThicknessAfterYears(thicknessMM, crMMPerYear, years float64) float64 {
	remaining := thicknessMM - crMMPerYear*years
	if remaining < 0 {
		return 0
	}
	return remaining
}

func PenetrationAfterYears(crMMPerYear, years float64) float64 {
	return crMMPerYear * years
}

type Schedule struct {
	MolarMass  float64
	Valence    float64
	ICorr      float64
	Area       float64
	Years      []float64
	CumLossesG []float64
	PerAreaG   []float64
}

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
