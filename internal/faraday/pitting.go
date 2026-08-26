package faraday

import "fmt"

func LocalRate(uniformMMPerYear, pittingFactor float64) (float64, error) {
	if uniformMMPerYear < 0 {
		return 0, fmt.Errorf("faraday: uniform rate must be non-negative")
	}
	if pittingFactor < 1 {
		return 0, fmt.Errorf("faraday: pitting factor must be ≥ 1")
	}
	return uniformMMPerYear * pittingFactor, nil
}

func YearsToPerforate(thicknessMM, uniformMMPerYear, pittingFactor float64) (float64, error) {
	local, err := LocalRate(uniformMMPerYear, pittingFactor)
	if err != nil {
		return 0, err
	}
	if local <= 0 {
		return positiveInfinity(), nil
	}
	if thicknessMM < 0 {
		return 0, fmt.Errorf("faraday: thickness must be non-negative")
	}
	return thicknessMM / local, nil
}

func PittingShortensLife(thicknessMM, uniformMMPerYear, pf float64) error {
	tUni, err := YearsToPerforate(thicknessMM, uniformMMPerYear, 1)
	if err != nil {
		return err
	}
	tPit, err := YearsToPerforate(thicknessMM, uniformMMPerYear, pf)
	if err != nil {
		return err
	}
	if pf <= 1 {
		return fmt.Errorf("faraday: pf must be > 1 to shorten life")
	}
	want := tUni / pf
	rel := (tPit - want) / want
	if rel < 0 {
		rel = -rel
	}
	if rel > 1e-12 {
		return fmt.Errorf("faraday: pitting life %g != uniform/pf %g", tPit, want)
	}
	if tPit >= tUni {
		return fmt.Errorf("faraday: pitting did not shorten life")
	}
	return nil
}

func RemainingWithPitting(thicknessMM, uniformMMPerYear, years, pf float64) (float64, error) {
	local, err := LocalRate(uniformMMPerYear, pf)
	if err != nil {
		return 0, err
	}
	return RemainingThicknessAfterYears(thicknessMM, local, years), nil
}
