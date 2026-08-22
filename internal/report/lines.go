package report

import (
	"io"
	"strings"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
)

// WriteComparison prints the iron-versus-aluminium depth-rate comparison
// that the specification calls out: at one shared current density the
// two metals must produce different CR values.
func WriteComparison(w io.Writer, cmp metal.Comparison) error {
	var b strings.Builder
	b.WriteString("Depth-rate comparison at i_corr = " + Num(cmp.ICorr) + " uA/cm^2:\n")
	b.WriteString("  " + cmp.First.Symbol + "  " + Num(cmp.FirstCR) + " mm/y\n")
	b.WriteString("  " + cmp.Second.Symbol + "  " + Num(cmp.SecondCR) + " mm/y\n")
	b.WriteString("  ratio " + cmp.First.Symbol + "/" + cmp.Second.Symbol + " = " + Num(cmp.RatioFirst) + "\n")
	_, err := io.WriteString(w, b.String())
	return err
}

// WriteSpecSummary prints the validated input values that will feed a
// computation, used by the validate subcommand to show what it accepted.
func WriteSpecSummary(w io.Writer, in faraday.Input) error {
	var b strings.Builder
	b.WriteString("i_corr     " + Num(in.ICorr) + " uA/cm^2\n")
	b.WriteString("M          " + Num(in.MolarMass) + " g/mol\n")
	b.WriteString("n          " + Num(in.Valence) + "\n")
	b.WriteString("rho        " + Num(in.Density) + " g/cm^3\n")
	if in.Area > 0 {
		b.WriteString("area       " + Num(in.Area) + " cm^2\n")
	}
	if in.DurationY > 0 {
		b.WriteString("duration_y " + Num(in.DurationY) + " y\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// Usage renders the CLI help text. The exact subcommand names and the
// example path are kept here so the program text and the README stay in
// agreement.
func Usage(w io.Writer, program string) error {
	_, err := io.WriteString(w, ``+program+` - uniform corrosion Faraday rate calculator

Usage:
  `+program+` rate <spec.json>          compute mass loss rate and depth rate
  `+program+` reverse <spec.json>       solve i_corr from a target rate (mm/y or g/(cm^2 y))
  `+program+` schedule <spec.json>      cumulative mass loss at 0.5/1/2/5/10 years
  `+program+` life <spec.json>          years to consume an allowance, remaining wall
  `+program+` validate <spec.json>      check a specification file without computing
  `+program+` chain <spec.json>         walk the unit chain and verify it closes
  `+program+` json <spec.json>          print the Result as JSON
  `+program+` metals                    list the built-in metal registry
  `+program+` invariants <spec.json>    show cross-rule ratios for a spec
  `+program+` template [metal]          print a ready-to-edit JSON spec
  `+program+` help                      show this help

Example:
  `+program+` rate example/fe-seawater.json

JSON fields (units fixed): i_corr uA/cm^2, M g/mol, n electrons, rho g/cm^3,
area cm^2 (optional), duration_y years (optional), metal symbol or name
(optional; fills M/n/rho when they are omitted). The reverse subcommand
also reads target_cr_mm_y or target_annual_loss.

Cross rules that must hold:
  i=0 => CR=0; i*2 => mdot*2 and CR*2; n*2 => rates halve;
  rho*2 => mass rate unchanged, depth rate halves;
  time*2 => cumulative loss*2.
`)
	return err
}
