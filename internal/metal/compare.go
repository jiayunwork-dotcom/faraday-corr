package metal

import "faraday-corr/internal/faraday"

type Comparison struct {
	ICorr      float64
	First      Metal
	Second     Metal
	FirstCR    float64
	SecondCR   float64
	RatioFirst float64
}

func CompareDepthRates(first, second Metal, iCorr float64) Comparison {
	firstCR := faraday.CorrosionRate(first.MolarMass, first.Valence, iCorr, first.Density)
	secondCR := faraday.CorrosionRate(second.MolarMass, second.Valence, iCorr, second.Density)
	return Comparison{
		ICorr:      iCorr,
		First:      first,
		Second:     second,
		FirstCR:    firstCR,
		SecondCR:   secondCR,
		RatioFirst: firstCR / secondCR,
	}
}

func CompareRegistry(iCorr float64) (Comparison, error) {
	first, err := Lookup("Fe")
	if err != nil {
		return Comparison{}, err
	}
	second, err := Lookup("Al")
	if err != nil {
		return Comparison{}, err
	}
	return CompareDepthRates(first, second, iCorr), nil
}
