package report

import (
	"fmt"
	"io"
	"strings"

	"faraday-corr/internal/metal"
)

type Table struct {
	Rows []metal.Metal
}

func NewTable(ms []metal.Metal) Table {
	return Table{Rows: ms}
}

func (t Table) Write(w io.Writer) error {
	if err := WriteMetals(w, t.Rows); err != nil {
		return err
	}
	return nil
}

func WriteAll(w io.Writer, ms []metal.Metal, refICorr float64) error {
	if err := WriteMetals(w, ms); err != nil {
		return err
	}
	cmp, err := metal.CompareRegistry(refICorr)
	if err != nil {
		return err
	}
	if err := WriteComparison(w, cmp); err != nil {
		return err
	}
	return nil
}

func ColumnWidth(lines []string) int {
	max := 0
	for _, ln := range lines {
		if len(ln) > max {
			max = len(ln)
		}
	}
	return max
}

func Header() string {
	return fmt.Sprintf("%-7s%-12s%-12s%-6s%-14s", "Symbol", "Name", "M (g/mol)", "n", "rho (g/cm^3)")
}

func Rule() string {
	return strings.Repeat("-", ColumnWidth([]string{Header()}))
}

func Blank() string {
	return ""
}

func FormatInvariantLine(name string, baseline, perturbed, ratio float64) string {
	return fmt.Sprintf("%-32s baseline=%s perturbed=%s ratio=%s",
		name, Num(baseline), Num(perturbed), Num(ratio))
}
