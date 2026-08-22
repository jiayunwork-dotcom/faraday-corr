// Command faraday-corr is a uniform-corrosion Faraday rate calculator.
// Given a corrosion current density i_corr, the metal's molar mass M,
// valence n and density rho, it prints the mass loss rate, the annual
// mass loss per unit area and the penetration depth rate, and - when an
// exposed area and exposure time are supplied - the total corrosion
// current and the cumulative mass loss.
//
// Examples:
//
//	go run . rate example/fe-seawater.json
//	go run . metals
//	go run . validate example/fe-seawater.json
//
// The physics live in internal/faraday; this file only dispatches
// subcommands and turns errors into stderr lines with a non-zero exit
// code.
package main

import (
	"fmt"
	"os"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
	"faraday-corr/internal/report"
)

func main() {
	if len(os.Args) < 2 {
		fail("missing subcommand; run 'faraday-corr help' for usage")
	}
	cmd := os.Args[1]
	rest := os.Args[2:]

	switch cmd {
	case "rate":
		runRate(rest)
	case "validate":
		runValidate(rest)
	case "metals":
		runMetals(rest)
	case "invariants":
		runInvariants(rest)
	case "reverse":
		runReverse(rest)
	case "schedule":
		runSchedule(rest)
	case "life":
		runLife(rest)
	case "template":
		runTemplate(rest)
	case "chain":
		runChain(rest)
	case "json":
		runJSON(rest)
	case "help", "-h", "--help":
		runHelp()
	default:
		fail(fmt.Sprintf("unknown subcommand %q; run 'faraday-corr help' for usage", cmd))
	}
}

// runRate loads a specification file, computes every derived quantity and
// prints the report.
func runRate(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	res, err := faraday.Compute(in)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteResult(os.Stdout, res, spec.Metal); err != nil {
		fail(err.Error())
	}
}

// runValidate checks a specification file and echoes the accepted input
// values, or fails on the first problem.
func runValidate(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	if issues := metal.ValidateSpec(spec); len(issues) > 0 {
		fail(metal.SummarizeIssues(issues).Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	fmt.Fprintln(os.Stdout, "OK: "+path)
	if err := report.WriteSpecSummary(os.Stdout, in); err != nil {
		fail(err.Error())
	}
}

// runMetals lists the built-in registry and the iron-versus-aluminium
// depth-rate comparison at a reference current density.
func runMetals(args []string) {
	if len(args) > 0 {
		fail("metals takes no arguments")
	}
	if err := report.WriteAll(os.Stdout, metal.All(), 10.0); err != nil {
		fail(err.Error())
	}
}

// runInvariants evaluates the cross rules against a specification and
// prints the observed ratios.
func runInvariants(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	checks, err := faraday.InvariantChecks(in)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteInvariants(os.Stdout, checks); err != nil {
		fail(err.Error())
	}
}

func runHelp() {
	if err := report.Usage(os.Stdout, "faraday-corr"); err != nil {
		fail(err.Error())
	}
}

// runReverse inverts the kernel: given a target depth rate or annual mass
// loss in the spec, it solves for the i_corr that produces it.
func runReverse(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	if spec.TargetCR != nil && spec.TargetLoss != nil {
		fail("set at most one of target_cr_mm_y and target_annual_loss")
	}
	target := faraday.TargetCorrosionRate
	if spec.TargetCR == nil && spec.TargetLoss == nil {
		fail("spec needs a target_cr_mm_y or target_annual_loss for the reverse subcommand")
	}
	if spec.TargetLoss != nil {
		target = faraday.TargetAnnualMassLoss
	}
	rin := faraday.ReverseInput{
		MolarMass: in.MolarMass,
		Valence:   in.Valence,
		Density:   in.Density,
		Target:    target,
	}
	if target == faraday.TargetCorrosionRate {
		rin.CR = *spec.TargetCR
	} else {
		rin.Annual = *spec.TargetLoss
	}
	res, err := faraday.ReverseCurrentDensity(rin)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteReverse(os.Stdout, res); err != nil {
		fail(err.Error())
	}
}

// runSchedule prints the cumulative mass loss at standard exposure times
// for a fixed input.
func runSchedule(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	if in.Area == 0 {
		fail("schedule needs an area in the spec")
	}
	sched, err := faraday.BuildSchedule(in, report.StandardScheduleYears)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteSchedule(os.Stdout, sched); err != nil {
		fail(err.Error())
	}
}

// runLife prints the life-cycle figures derived from the depth rate.
// Flags: --allowance <mm> corrosion allowance, --thickness <mm> initial
// wall, --years <y> service interval. The spec path is the first
// positional argument.
func runLife(args []string) {
	if len(args) == 0 {
		fail("life needs a specification file; run 'faraday-corr help' for usage")
	}
	path := args[0]
	allowance := 1.0
	thickness := 10.0
	years := 5.0
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--allowance":
			if i+1 >= len(args) {
				fail("--allowance needs a value")
			}
			i++
			if _, err := fmt.Sscanf(args[i], "%f", &allowance); err != nil {
				fail("bad --allowance value")
			}
		case "--thickness":
			if i+1 >= len(args) {
				fail("--thickness needs a value")
			}
			i++
			if _, err := fmt.Sscanf(args[i], "%f", &thickness); err != nil {
				fail("bad --thickness value")
			}
		case "--years":
			if i+1 >= len(args) {
				fail("--years needs a value")
			}
			i++
			if _, err := fmt.Sscanf(args[i], "%f", &years); err != nil {
				fail("bad --years value")
			}
		default:
			fail(fmt.Sprintf("unknown life flag %q", args[i]))
		}
	}
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	res, err := faraday.Compute(in)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteLife(os.Stdout, res.CorrosionRate, allowance, thickness, years); err != nil {
		fail(err.Error())
	}
}

// runChain walks the unit chain end to end and reports whether the
// mass-loss path and the direct K path agree on the depth rate.
func runChain(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	chain, err := faraday.VerifyUnitChain(in)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteChain(os.Stdout, chain); err != nil {
		fail(err.Error())
	}
}

// runJSON prints the computed Result as JSON for script consumption.
func runJSON(args []string) {
	path := requireSpecPath(args)
	spec, err := metal.LoadFile(path)
	if err != nil {
		fail(err.Error())
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		fail(err.Error())
	}
	res, err := faraday.Compute(in)
	if err != nil {
		fail(err.Error())
	}
	if err := report.WriteResultJSON(os.Stdout, res); err != nil {
		fail(err.Error())
	}
}

// runTemplate prints a ready-to-edit JSON specification for a metal.
func runTemplate(args []string) {
	metalName := ""
	if len(args) > 1 {
		fail("template takes at most one metal symbol or name")
	}
	if len(args) == 1 {
		metalName = args[0]
	}
	spec, err := metal.TemplateSpec(metalName)
	if err != nil {
		fail(err.Error())
	}
	if err := metal.WriteJSON(os.Stdout, spec); err != nil {
		fail(err.Error())
	}
}

// requireSpecPath returns the single positional argument that the
// spec-taking subcommands expect, failing with usage guidance otherwise.
func requireSpecPath(args []string) string {
	if len(args) < 1 {
		fail("expected a specification file argument; run 'faraday-corr help' for usage")
	}
	return args[0]
}

// fail writes an error to stderr and exits non-zero. No error path in
// this program continues after calling fail.
func fail(msg string) {
	fmt.Fprintln(os.Stderr, "faraday-corr: error:", msg)
	os.Exit(1)
}
