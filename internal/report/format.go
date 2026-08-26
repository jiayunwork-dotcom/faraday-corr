package report

import (
	"fmt"
	"math"
)

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

func NumFixed(v float64) string {
	return fmt.Sprintf("%.4f", v)
}

func PercentDiff(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	return math.Abs(a-b) / math.Max(math.Abs(a), 1e-30) * 100
}
