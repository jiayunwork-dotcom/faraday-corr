package faraday

// InvariantCheckResult summarises one cross-rule check between a baseline
// computation and a perturbed one. The cross rules from the specification
// are: i=0 => CR=0, i*2 => mdot*2 and CR*2, n*2 => both rates halve,
// rho*2 => mass rate unchanged while depth rate halves, and doubling the
// exposure time doubles the cumulative mass loss. Each check is expressed
// as a ratio of the perturbed quantity to the baseline quantity.
type InvariantCheckResult struct {
	Name      string
	Baseline  float64
	Perturbed float64
	Ratio     float64
}

// InvariantChecks evaluates every cross rule for an input and returns the
// observed ratios. It is used by the test suite to pin the rules and can
// be invoked from a subcommand for ad-hoc verification.
func InvariantChecks(in Input) ([]InvariantCheckResult, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	base, err := Compute(in)
	if err != nil {
		return nil, err
	}

	var out []InvariantCheckResult

	zeroI, _ := Compute(in.WithCurrent(0))
	out = append(out, InvariantCheckResult{
		Name:      "i=0 => CR=0",
		Baseline:  base.CorrosionRate,
		Perturbed: zeroI.CorrosionRate,
		Ratio:     zeroI.CorrosionRate / base.CorrosionRate,
	})

	doubleI, _ := Compute(in.WithCurrent(in.ICorr * 2))
	out = append(out, InvariantCheckResult{
		Name:      "i*2 => mdot*2",
		Baseline:  base.MassLossRate,
		Perturbed: doubleI.MassLossRate,
		Ratio:     doubleI.MassLossRate / base.MassLossRate,
	})
	out = append(out, InvariantCheckResult{
		Name:      "i*2 => CR*2",
		Baseline:  base.CorrosionRate,
		Perturbed: doubleI.CorrosionRate,
		Ratio:     doubleI.CorrosionRate / base.CorrosionRate,
	})

	doubleN, _ := Compute(in.WithValence(in.Valence * 2))
	out = append(out, InvariantCheckResult{
		Name:      "n*2 => mdot/2",
		Baseline:  base.MassLossRate,
		Perturbed: doubleN.MassLossRate,
		Ratio:     doubleN.MassLossRate / base.MassLossRate,
	})
	out = append(out, InvariantCheckResult{
		Name:      "n*2 => CR/2",
		Baseline:  base.CorrosionRate,
		Perturbed: doubleN.CorrosionRate,
		Ratio:     doubleN.CorrosionRate / base.CorrosionRate,
	})

	doubleRho, _ := Compute(in.WithDensity(in.Density * 2))
	out = append(out, InvariantCheckResult{
		Name:      "rho*2 => mdot unchanged",
		Baseline:  base.MassLossRate,
		Perturbed: doubleRho.MassLossRate,
		Ratio:     doubleRho.MassLossRate / base.MassLossRate,
	})
	out = append(out, InvariantCheckResult{
		Name:      "rho*2 => CR/2",
		Baseline:  base.CorrosionRate,
		Perturbed: doubleRho.CorrosionRate,
		Ratio:     doubleRho.CorrosionRate / base.CorrosionRate,
	})

	doubleT, _ := Compute(in.WithDuration(in.DurationY * 2))
	out = append(out, InvariantCheckResult{
		Name:      "t*2 => cumulative loss*2",
		Baseline:  base.CumulativeMassLoss,
		Perturbed: doubleT.CumulativeMassLoss,
		Ratio:     doubleT.CumulativeMassLoss / base.CumulativeMassLoss,
	})

	return out, nil
}
