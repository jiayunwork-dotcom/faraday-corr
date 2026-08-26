package faraday

import (
	"fmt"
	"math"
)

func DiffusionLimit(n, diffusivity, concMolM3, deltaM float64) (float64, error) {
	if n <= 0 || diffusivity <= 0 || concMolM3 <= 0 || deltaM <= 0 {
		return 0, fmt.Errorf("faraday: n, D, C and film thickness must be positive")
	}
	iA := n * Faraday * diffusivity * concMolM3 / deltaM
	return iA * 1e6, nil
}

func MixedActivationDiffusion(iAct, iLim float64) (float64, error) {
	if iAct <= 0 || iLim <= 0 {
		return 0, fmt.Errorf("faraday: activation and limit currents must be positive")
	}
	return 1 / (1/iAct + 1/iLim), nil
}

func ObservedBelowBoth(iAct, iLim float64) error {
	obs, err := MixedActivationDiffusion(iAct, iLim)
	if err != nil {
		return err
	}
	if obs >= iAct || obs >= iLim {
		return fmt.Errorf("faraday: observed %g not below both %g and %g", obs, iAct, iLim)
	}
	return nil
}

func ThinFilmRaisesLimit(n, d, c, deltaThin, deltaThick float64) error {
	if deltaThick <= deltaThin {
		return fmt.Errorf("faraday: thick film must exceed thin")
	}
	thin, err := DiffusionLimit(n, d, c, deltaThin)
	if err != nil {
		return err
	}
	thick, err := DiffusionLimit(n, d, c, deltaThick)
	if err != nil {
		return err
	}
	if thin <= thick {
		return fmt.Errorf("faraday: thinner film did not raise iL (%g vs %g)", thin, thick)
	}
	ratio := thin / thick
	want := deltaThick / deltaThin
	if math.Abs(ratio-want)/want > 1e-9 {
		return fmt.Errorf("faraday: iL ratio %g want %g", ratio, want)
	}
	return nil
}
