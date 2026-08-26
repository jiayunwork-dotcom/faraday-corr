package faraday

func Compute(in Input) (Result, error) {
	if err := Validate(in); err != nil {
		return Result{}, err
	}

	mdot := MassLossRate(in.MolarMass, in.Valence, in.ICorr)
	annualPerArea := AnnualMassLossPerAreaFromRate(mdot)
	cr := CorrosionRate(in.MolarMass, in.Valence, in.ICorr, in.Density)

	res := Result{
		Input:              in,
		MassLossRate:       mdot,
		AnnualMassLoss:     annualPerArea,
		CorrosionRate:      cr,
		CorrosionRateUmY:   CorrosionRateUmY(cr),
		TotalCurrent:       TotalCurrent(in.ICorr, in.Area),
		TotalCurrentA:      TotalCurrentAmps(in.ICorr, in.Area),
		CumulativeMassLoss: CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, in.DurationY),
	}
	return res, nil
}

func EquivalentWeight(molarMass, valence float64) float64 {
	return molarMass / valence
}

func DepthRatePerEquivalentWeight(molarMass, valence float64) float64 {
	return K * EquivalentWeight(molarMass, valence)
}
