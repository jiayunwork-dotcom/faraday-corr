package metal

import (
	"fmt"

	"faraday-corr/internal/faraday"
)

func ToInput(spec Spec) (faraday.Input, error) {
	filled, err := spec.FillFromRegistry()
	if err != nil {
		return faraday.Input{}, err
	}
	in := faraday.Input{
		ICorr:     filled.ICorr,
		MolarMass: value(filled.MolarMass),
		Valence:   value(filled.Valence),
		Density:   value(filled.Density),
		Area:      value(filled.Area),
		DurationY: value(filled.DurationY),
	}
	if err := faraday.Validate(in); err != nil {
		return faraday.Input{}, fmt.Errorf("invalid specification: %w", err)
	}
	return in, nil
}

func value(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func BuildInput(spec Spec) (faraday.Input, error) {
	return ToInput(spec)
}
