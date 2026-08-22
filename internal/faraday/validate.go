package faraday

import (
	"math"
)

// Validate rejects inputs that would produce meaningless or physically
// impossible rates. A negative corrosion current density, a non-positive
// valence, a non-positive density or a non-positive molar mass each yield
// an error; a negative area or duration is also an error because negative
// exposure makes no sense. A zero current density is allowed: the cross
// rule i=0 => CR=0 requires it to compute zero rates rather than error.
// Zero area or zero duration is allowed and means "not given".
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
