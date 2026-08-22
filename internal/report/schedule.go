package report

import (
	"fmt"
	"io"
	"strings"

	"faraday-corr/internal/faraday"
)

// StandardScheduleYears are the exposure times shown by the schedule
// subcommand: the uniform-corrosion linearity is visible at a glance.
var StandardScheduleYears = []float64{0.5, 1, 2, 5, 10}

// WriteSchedule prints the cumulative mass loss at each exposure time of
// a Schedule, per unit area and over the stated area.
func WriteSchedule(w io.Writer, s faraday.Schedule) error {
	var b strings.Builder
	b.WriteString("Uniform corrosion schedule (i_corr " + Num(s.ICorr) + " uA/cm^2, area " + Num(s.Area) + " cm^2):\n")
	b.WriteString("years    per area (g/cm^2)   cumulative (g)\n")
	b.WriteString("------   -----------------   --------------\n")
	for i, y := range s.Years {
		b.WriteString(fmt.Sprintf("%-9v %-20s %s\n",
			Num(y), Num(s.PerAreaG[i]), Num(s.CumLossesG[i])))
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// WriteTemplateEcho prints a generated template spec together with its
// unit legend, used by the template subcommand.
func WriteTemplateEcho(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}
