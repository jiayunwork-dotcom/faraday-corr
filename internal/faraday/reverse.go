package faraday

// Reverse target selectors. The forward kernel takes i_corr and derives
// rates; the reverse path takes one target rate and solves for the
// current density that would produce it. Only one target may be set.
const (
	TargetNone = iota
	TargetCorrosionRate
	TargetAnnualMassLoss
)

// ReverseInput carries the material parameters plus exactly one target
// rate to invert.
type ReverseInput struct {
	MolarMass float64 // g/mol
	Valence   float64
	Density   float64 // g/cm^3
	Target    int     // TargetCorrosionRate or TargetAnnualMassLoss
	CR        float64 // mm/y, used when Target == TargetCorrosionRate
	Annual    float64 // g/(cm^2 y), used when Target == TargetAnnualMassLoss
}

// ReverseResult is the outcome of a reverse solve.
type ReverseResult struct {
	ReverseInput
	ICorr float64 // uA/cm^2
	// Derived values of the forward path at the solved current density.
	MassLossRate     float64
	AnnualMassLoss   float64
	CorrosionRate    float64
	CorrosionRateUmY float64
}

// perYearCurrentFold is i[A/cm^2] * yearSeconds folded into the annual
// mass-loss expression: M * i[uA] * 1e-6 * yearSeconds / (n*F).
const perYearCurrentFold = 1e-6 * yearSeconds

// ReverseCurrentDensity inverts the kernel. From CR = K*(M/n)*i/rho the
// current density is i = CR*rho/(K*(M/n)); from the annual mass loss
// annual = M*i*1e-6*yearSeconds/(n*F) it is i = annual*n*F/(M*1e-6*years).
// Both inverses are the algebraic inverse of the forward formulas, so
// round-tripping i -> rates -> i recovers the original current density.
func ReverseCurrentDensity(in ReverseInput) (ReverseResult, error) {
	if err := validateReverse(in); err != nil {
		return ReverseResult{}, err
	}
	var i float64
	switch in.Target {
	case TargetCorrosionRate:
		i = in.CR * in.Density / (K * (in.MolarMass / in.Valence))
	case TargetAnnualMassLoss:
		i = in.Annual * in.Valence * Faraday / (in.MolarMass * perYearCurrentFold)
	default:
		return ReverseResult{}, &validationError{field: "target", value: float64(in.Target), reason: "must name a rate to invert"}
	}
	return ReverseResult{
		ReverseInput:     in,
		ICorr:            takeReverseI(i),
		MassLossRate:     MassLossRate(in.MolarMass, in.Valence, i),
		AnnualMassLoss:   AnnualMassLossPerArea(in.MolarMass, in.Valence, i),
		CorrosionRate:    CorrosionRate(in.MolarMass, in.Valence, i, in.Density),
		CorrosionRateUmY: CorrosionRateUmY(CorrosionRate(in.MolarMass, in.Valence, i, in.Density)),
	}, nil
}

// validateReverse checks that the material values are positive and the
// chosen target rate is positive.
func validateReverse(in ReverseInput) error {
	if in.MolarMass <= 0 {
		return errNonPositive("M", in.MolarMass)
	}
	if in.Valence <= 0 {
		return errNonPositive("n", in.Valence)
	}
	if in.Density <= 0 {
		return errNonPositive("rho", in.Density)
	}
	switch in.Target {
	case TargetCorrosionRate:
		if in.CR <= 0 {
			return errNonPositive("target CR", in.CR)
		}
	case TargetAnnualMassLoss:
		if in.Annual <= 0 {
			return errNonPositive("target annual mass loss", in.Annual)
		}
	}
	return nil
}

// RoundTrip verifies that a forward compute followed by a reverse solve
// recovers the original current density. It is the self-consistency check
// for the kernel pair.
func RoundTrip(in Input) (ReverseResult, error) {
	res, err := Compute(in)
	if err != nil {
		return ReverseResult{}, err
	}
	return ReverseCurrentDensity(ReverseInput{
		MolarMass: in.MolarMass,
		Valence:   in.Valence,
		Density:   in.Density,
		Target:    TargetCorrosionRate,
		CR:        res.CorrosionRate,
	})
}
