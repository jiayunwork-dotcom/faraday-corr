package report

import (
	"fmt"
	"io"
	"strings"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
)

func WriteResult(w io.Writer, res faraday.Result, label string) error {
	if label == "" {
		label = "unknown"
	}

	var b strings.Builder
	b.WriteString("Material:          " + label + "\n")
	b.WriteString("i_corr:            " + Num(res.ICorr) + " uA/cm^2\n")
	b.WriteString("M/n:               " + Num(res.MolarMass/res.Valence) + " g/mol (equivalent weight)\n")
	b.WriteString("rho:               " + Num(res.Density) + " g/cm^3\n")
	b.WriteString("Mass loss rate:    " + Num(res.MassLossRate) + " g/(cm^2 s)\n")
	b.WriteString("Annual mass loss:  " + Num(res.AnnualMassLoss) + " g/(cm^2 y)\n")
	b.WriteString("Corrosion rate:    " + Num(res.CorrosionRate) + " mm/y (" + Num(res.CorrosionRateUmY) + " um/y)\n")

	if res.Area > 0 {
		b.WriteString("I = i_corr*A:      " + Num(res.TotalCurrent) + " uA = " + Num(res.TotalCurrentA) + " A (area " + Num(res.Area) + " cm^2)\n")
	}
	if res.Area > 0 && res.DurationY > 0 {
		b.WriteString("Cumulative loss:   " + Num(res.CumulativeMassLoss) + " g (area " + Num(res.Area) + " cm^2, " + Num(res.DurationY) + " y)\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func WriteMetals(w io.Writer, metals []metal.Metal) error {
	var b strings.Builder
	b.WriteString("Symbol  Name        M (g/mol)   n     rho (g/cm^3)\n")
	b.WriteString("------  ----------  ----------  ----  ------------\n")
	for _, m := range metals {
		b.WriteString(pad(m.Symbol, 7) +
			pad(m.Name, 12) +
			pad(Num(m.MolarMass), 12) +
			pad(Num(m.Valence), 6) +
			Num(m.Density) + "\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func WriteInvariants(w io.Writer, checks []faraday.InvariantCheckResult) error {
	var b strings.Builder
	for _, c := range checks {
		b.WriteString(pad(c.Name, 28) + "ratio=" + Num(c.Ratio) + "\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s + " "
	}
	return fmt.Sprintf("%-*s", width, s)
}
