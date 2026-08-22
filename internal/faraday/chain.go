package faraday

// UnitChain walks the unit conversion end to end for one input and
// reports each intermediate value, so the closing of the chain can be
// inspected instead of assumed. The two depth-rate paths must agree: the
// direct formula K*(M/n)*i/rho and the mass-loss path 10*annual/rho. The
// effective K re-derived from the folded constants must equal the pinned
// K.
type UnitChain struct {
	Input

	// CurrentDensityA is i_corr converted to A/cm^2.
	CurrentDensityA float64
	// MdotPerS is the mass loss rate in g/(cm^2 s).
	MdotPerS float64
	// AnnualPerCm2 is the mass loss rate in g/(cm^2 y).
	AnnualPerCm2 float64
	// CRFromMass is the depth rate derived from the mass path.
	CRFromMass float64
	// CRDirect is the depth rate from the pinned K formula.
	CRDirect float64
	// KEffective is K re-derived from 10*1e-6*yearSeconds/F.
	KEffective float64
	// Consistent reports whether the two depth-rate paths agree within
	// the relative tolerance.
	Consistent bool
}

// VerifyUnitChain computes and cross-checks every step of the unit
// chain for a valid input. It fails on invalid input, not on a chain
// mismatch: the mismatch is reported in Consistent.
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

// almostEqualWithin compares two floats with a given absolute tolerance,
// used by the chain consistency check where both sides share the same
// order of magnitude.
func almostEqualWithin(a, b, tol float64) bool {
	return a-b < tol && b-a < tol
}

// KDerived recomputes the unit-chain constant from its raw parts and
// reports the relative difference against the pinned K.
func KDerived() (derived, diff float64) {
	derived = milliPerCm * (1.0 / microPerAmpere) * yearSeconds / Faraday
	diff = (derived - K) / K
	return derived, diff
}

// UnitsUsed lists every unit label the kernel emits, for documentation
// and for the chain subcommand header.
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
