package report

import (
	"encoding/json"
	"io"
	"strings"

	"faraday-corr/internal/faraday"
)

func WriteChain(w io.Writer, c faraday.UnitChain) error {
	var b strings.Builder
	b.WriteString("Unit chain for i_corr " + Num(c.ICorr) + " uA/cm^2, M " + Num(c.MolarMass) + " g/mol, n " +
		Num(c.Valence) + ", rho " + Num(c.Density) + " g/cm^3:\n")
	b.WriteString("  i_corr (A/cm^2)   " + Num(c.CurrentDensityA) + "\n")
	b.WriteString("  mdot (g/(cm^2 s)) " + Num(c.MdotPerS) + " = M*i/(n*F)\n")
	b.WriteString("  annual (g/(cm^2 y)) " + Num(c.AnnualPerCm2) + " = mdot * 3.1536e7 s\n")
	b.WriteString("  CR from mass path:  " + Num(c.CRFromMass) + " mm/y = 10*annual/rho\n")
	b.WriteString("  CR direct (K eq):   " + Num(c.CRDirect) + " mm/y\n")
	b.WriteString("  K effective:        " + Num(c.KEffective) + "\n")
	b.WriteString("  chain consistent:   " + boolWord(c.Consistent) + "\n")
	_, err := io.WriteString(w, b.String())
	return err
}

func WriteResultJSON(w io.Writer, res faraday.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(res)
}

func WriteErrorJSON(w io.Writer, msg string) error {
	enc := json.NewEncoder(w)
	return enc.Encode(map[string]string{"error": msg})
}

func boolWord(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
