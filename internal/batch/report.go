package batch

import (
	"fmt"
	"strings"
)

func TextReport(r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "== faraday-corr batch: %d cases, %d ok ==\n", r.Total, r.Success)
	for _, it := range r.Items {
		if it.Error != "" {
			fmt.Fprintf(&b, "FAIL  %s: %s\n", it.Name, it.Error)
			continue
		}
		fmt.Fprintf(&b, "OK    %s  rate=%.2f µm/y  loss=%.4f g/m²/y  I=%.4e A\n",
			it.Name, it.CorrosionRateUmY, it.AnnualMassLoss, it.TotalCurrentA)
	}
	return b.String()
}

func FailedNames(r Result) []string {
	var out []string
	for _, it := range r.Items {
		if it.Error != "" {
			out = append(out, it.Name)
		}
	}
	return out
}
