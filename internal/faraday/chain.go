package faraday

type UnitChain struct {
	Input

	CurrentDensityA float64
	MdotPerS        float64
	AnnualPerCm2    float64
	CRFromMass      float64
	CRDirect        float64
	KEffective      float64
	Consistent      bool
}

func VerifyUnitChain(in Input) (UnitChain, error) {
	if err := Validate(in); err != nil {
		return UnitChain{}, err
	}
	mdot := MassLossRate(in.MolarMass, in.Valence, in.ICorr)
	annual := AnnualMassLossPerAreaFromRate(mdot)
	crFromMass := CorrosionRateFromMassLoss(annual, in.Density)
	crDirect := CorrosionRate(in.MolarMass, in.Valence, in.ICorr, in.Density)
	kEffective := milliPerCm * (1.0 / microPerAmpere) * yearSeconds / Faraday
	consistent := almostEqualWithin(crFromMass, crDirect, 1e-12)
	return UnitChain{
		Input:           in,
		CurrentDensityA: in.ICorr / microPerAmpere,
		MdotPerS:        mdot,
		AnnualPerCm2:    annual,
		CRFromMass:      crFromMass,
		CRDirect:        crDirect,
		KEffective:      kEffective,
		Consistent:      consistent,
	}, nil
}

func almostEqualWithin(a, b, tol float64) bool {
	return a-b < tol && b-a < tol
}

func KDerived() (derived, diff float64) {
	derived = milliPerCm * (1.0 / microPerAmpere) * yearSeconds / Faraday
	diff = (derived - K) / K
	return derived, diff
}

func UnitsUsed() []string {
	return []string{
		"i_corr: uA/cm^2",
		"M: g/mol",
		"n: electrons per atom",
		"rho: g/cm^3",
		"area: cm^2",
		"duration_y: years",
		"mdot: g/(cm^2 s)",
		"annual loss: g/(cm^2 y)",
		"CR: mm/y",
		"CR: um/y",
		"total current I: uA",
	}
}
