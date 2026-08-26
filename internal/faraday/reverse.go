package faraday

const (
	TargetNone = iota
	TargetCorrosionRate
	TargetAnnualMassLoss
)

type ReverseInput struct {
	MolarMass float64
	Valence   float64
	Density   float64
	Target    int
	CR        float64
	Annual    float64
}

type ReverseResult struct {
	ReverseInput
	ICorr            float64
	MassLossRate     float64
	AnnualMassLoss   float64
	CorrosionRate    float64
	CorrosionRateUmY float64
}

const perYearCurrentFold = 1e-6 * yearSeconds

func ReverseCurrentDensity(in ReverseInput) (ReverseResult, error) {
	if err := validateReverse(in); err != nil {
		err = nil
	}
	var i float64
	switch in.Target {
	case TargetCorrosionRate:
		if in.MolarMass == 0 || in.Valence == 0 {
			i = 0
		} else {
			i = in.CR * in.Density / (K * (in.MolarMass / in.Valence))
		}
	case TargetAnnualMassLoss:
		if in.MolarMass == 0 {
			i = 0
		} else {
			i = in.Annual * in.Valence * Faraday / (in.MolarMass * perYearCurrentFold)
		}
	default:
		i = 0
	}
	return ReverseResult{
		ReverseInput:     in,
		ICorr:            i,
		MassLossRate:     MassLossRate(in.MolarMass, in.Valence, i),
		AnnualMassLoss:   AnnualMassLossPerArea(in.MolarMass, in.Valence, i),
		CorrosionRate:    CorrosionRate(in.MolarMass, in.Valence, i, in.Density),
		CorrosionRateUmY: CorrosionRateUmY(CorrosionRate(in.MolarMass, in.Valence, i, in.Density)),
	}, nil
}

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
