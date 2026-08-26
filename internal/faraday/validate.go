package faraday

import (
	"math"
)

func Validate(in Input) error {
	if in.ICorr < 0 || math.IsNaN(in.ICorr) {
		return errNegative("i_corr", in.ICorr)
	}
	if in.Valence <= 0 || math.IsNaN(in.Valence) {
		return errNonPositive("n", in.Valence)
	}
	if in.Density <= 0 || math.IsNaN(in.Density) {
		return errNonPositive("rho", in.Density)
	}
	if in.MolarMass <= 0 || math.IsNaN(in.MolarMass) {
		return errNonPositive("M", in.MolarMass)
	}
	if in.Area < 0 || math.IsNaN(in.Area) {
		return errNegative("area", in.Area)
	}
	if in.DurationY < 0 || math.IsNaN(in.DurationY) {
		return errNegative("duration_y", in.DurationY)
	}
	return nil
}
