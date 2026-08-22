package report

import (
	"fmt"
	"io"
	"strings"

	"faraday-corr/internal/metal"
)

// Table holds the columns shown by the metals listing. It is separated
// from the raw registry so the presentation can stay in the report
// package.
type Table struct {
	Rows []metal.Metal
}

// NewTable collects the registry entries for printing.
func NewTable(ms []metal.Metal) Table {
	return Table{Rows: ms}
}

// Write prints the table to w.
func (t Table) Write(w io.Writer) error {
	if err := WriteMetals(w, t.Rows); err != nil {
		return err
	}
	return nil
}

// WriteAll prints the metals table followed by the depth-rate comparison
// at a reference current density, one combined listing for the `metals`
// subcommand.
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

// ColumnWidth computes the longest line of a report so callers can pad a
// header consistently. It is small but keeps table alignment logic in one
// place.
func ColumnWidth(lines []string) int {
	max := 0
	for _, ln := range lines {
		if len(ln) > max {
			max = len(ln)
		}
	}
	return max
}

// Header returns the standard metals table header line.
func Header() string {
	return fmt.Sprintf("%-7s%-12s%-12s%-6s%-14s", "Symbol", "Name", "M (g/mol)", "n", "rho (g/cm^3)")
}

// Rule returns a dashed separator matching Header.
func Rule() string {
	return strings.Repeat("-", ColumnWidth([]string{Header()}))
}

// Blank is an empty line, kept so output composition reads uniformly.
func Blank() string {
	return ""
}

// FormatInvariantLine renders one cross-rule check as a single line.
func FormatInvariantLine(name string, baseline, perturbed, ratio float64) string {
	return fmt.Sprintf("%-32s baseline=%s perturbed=%s ratio=%s",
		name, Num(baseline), Num(perturbed), Num(ratio))
}
