package report

import (
	"fmt"
	"math"
)

// Num renders a float with a sensible number of significant digits.
// Values that are very large or very small switch to scientific
// notation; rates such as 2.894e-9 g/(cm^2 s) would otherwise print as
// an unreadable wall of zeros.
func Num(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Sprintf("%v", v)
	}
	av := math.Abs(v)
	switch {
	case av == 0:
		return "0"
	case av >= 1e6 || av < 1e-4:
		return fmt.Sprintf("%.3e", v)
	default:
		return fmt.Sprintf("%.4g", v)
	}
}

// NumFixed renders a float with up to four fixed decimals, used for
// inputs echoed back in the report where scientific notation would be
// distracting.
func NumFixed(v float64) string {
	return fmt.Sprintf("%.4f", v)
}

// PercentDiff returns the relative difference between two positive
// values as a percentage, used by the invariant verification output.
func PercentDiff(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	return math.Abs(a-b) / math.Max(math.Abs(a), 1e-30) * 100
}
