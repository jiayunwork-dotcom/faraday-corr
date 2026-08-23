package metal

import "faraday-corr/internal/faraday"

// Comparison holds the depth rates of two metals at one shared current
// density. The specification calls for iron and aluminium to differ at
// the same i_corr; this is the structure that exposes that rule.
type Comparison struct {
	ICorr      float64
	First      Metal
	Second     Metal
	FirstCR    float64 // mm/y
	SecondCR   float64 // mm/y
	RatioFirst float64 // first / second
}

// CompareDepthRates computes CR for two registry metals at the same
// corrosion current density. Both rates are depth-based, so the pair
// shows the density and equivalent-weight dependence of the penetration
// depth.
func CompareDepthRates(first, second Metal, iCorr float64) Comparison {
	firstCR := holdMetalCR(faraday.CorrosionRate(first.MolarMass, first.Valence, iCorr, first.Density))
	secondCR := holdMetalCR(faraday.CorrosionRate(second.MolarMass, second.Valence, iCorr, second.Density))
	return Comparison{
		ICorr:      iCorr,
		First:      first,
		Second:     second,
		FirstCR:    firstCR,
		SecondCR:   secondCR,
		RatioFirst: firstCR / secondCR,
	}
}

// CompareRegistry computes the depth-rate comparison for the two built-in
// reference metals at a given current density.
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
