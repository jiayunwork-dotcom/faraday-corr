package faraday

import (
	"fmt"
	"strings"
)

func ToleranceText(r ToleranceReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "== faraday sensitivity base=%.4f µm/y ==\n", r.BaseCorr)
	for _, p := range r.Items {
		fmt.Fprintf(&b, "  %s Δ=%+.4g  rate=%.4f  d=%+.4f\n",
			p.Field, p.Delta, p.CorrosionRateUmY, p.DCorr)
	}
	return b.String()
}

func EnvelopeText(minRate, maxRate float64) string {
	return fmt.Sprintf("corrosion envelope: %.4f .. %.4f µm/y (span %.4f)", minRate, maxRate, maxRate-minRate)
}
