package report

import (
	"fmt"
	"io"
	"strings"

	"faraday-corr/internal/faraday"
)

var StandardScheduleYears = []float64{0.5, 1, 2, 5, 10}

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

func WriteTemplateEcho(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}
