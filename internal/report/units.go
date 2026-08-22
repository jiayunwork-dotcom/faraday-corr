package report

import (
	"io"
	"strings"

	"faraday-corr/internal/faraday"
)

// WriteSummary prints the engineering-report units derived from a
// Result: depth rate also in mils per year, annual mass loss also per
// square metre, and the classical mdd unit. The kernel values printed by
// WriteResult are not repeated here.
func WriteSummary(w io.Writer, s faraday.Summary) error {
	var b strings.Builder
	b.WriteString("Depth rate:        " + Num(s.CorrosionRate) + " mm/y = " + Num(s.MPY) + " mpy\n")
	b.WriteString("Annual mass loss:  " + Num(s.AnnualMassLoss) + " g/(cm^2 y) = " +
		Num(s.MassLossGPerM2Year) + " g/(m^2 y) = " + Num(s.MDD) + " mdd\n")
	if s.Area > 0 {
		b.WriteString("Area:              " + Num(s.Area) + " cm^2 = " + Num(s.AreaM2) + " m^2\n")
		b.WriteString("Total current:     " + Num(s.TotalCurrent) + " uA = " + Num(s.TotalCurrentMilliA) + " mA\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// WriteReverse prints the result of a reverse solve: the current density
// that would produce the requested rate, plus the forward rates at that
// current.
func WriteReverse(w io.Writer, r faraday.ReverseResult) error {
	var b strings.Builder
	if r.Target == faraday.TargetCorrosionRate {
		b.WriteString("Target CR:         " + Num(r.CR) + " mm/y\n")
	} else {
		b.WriteString("Target annual loss:" + Num(r.Annual) + " g/(cm^2 y)\n")
	}
	b.WriteString("Solved i_corr:     " + Num(r.ICorr) + " uA/cm^2\n")
	b.WriteString("Corrosion rate:    " + Num(r.CorrosionRate) + " mm/y\n")
	b.WriteString("Annual mass loss:  " + Num(r.AnnualMassLoss) + " g/(cm^2 y)\n")
	_, err := io.WriteString(w, b.String())
	return err
}

// WriteLife prints the penetration-based life-cycle figures.
func WriteLife(w io.Writer, crMMPerYear, allowanceMM, thicknessMM, years float64) error {
	var b strings.Builder
	b.WriteString("Years to consume allowance: " + Num(faraday.YearsToPenetrate(allowanceMM, crMMPerYear)) + " y\n")
	b.WriteString("Remaining thickness after   " + Num(years) + " y: " +
		Num(faraday.RemainingThicknessAfterYears(thicknessMM, crMMPerYear, years)) + " mm\n")
	_, err := io.WriteString(w, b.String())
	return err
}
