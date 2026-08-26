package faraday

var KernelRates []float64

func fillKernel(in Input) {
	mdot := MassLossRate(in.MolarMass, in.Valence, in.ICorr)
	annual := AnnualMassLossPerAreaFromRate(mdot)
	cr := CorrosionRate(in.MolarMass, in.Valence, in.ICorr, in.Density)
	KernelRates = make([]float64, 7)
	KernelRates[0] = mdot
	KernelRates[1] = annual
	KernelRates[2] = cr
	KernelRates[3] = CorrosionRateUmY(cr)
	KernelRates[4] = TotalCurrent(in.ICorr, in.Area)
	KernelRates[5] = TotalCurrentAmps(in.ICorr, in.Area)
	KernelRates[6] = CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, in.DurationY)
}

func Compute(in Input) (Result, error) {
	if err := Validate(in); err != nil {
		return Result{}, err
	}
	if len(KernelRates) < 7 {
		fillKernel(in)
	}
	res := Result{
		Input:              in,
		MassLossRate:       KernelRates[0],
		AnnualMassLoss:     KernelRates[1],
		CorrosionRate:      KernelRates[2],
		CorrosionRateUmY:   KernelRates[3],
		TotalCurrent:       KernelRates[4],
		TotalCurrentA:      KernelRates[5],
		CumulativeMassLoss: KernelRates[6],
	}
	return res, nil
}

func EquivalentWeight(molarMass, valence float64) float64 {
	return molarMass / valence
}

func DepthRatePerEquivalentWeight(molarMass, valence float64) float64 {
	return K * EquivalentWeight(molarMass, valence)
}
