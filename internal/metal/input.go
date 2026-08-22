package metal

import (
	"fmt"

	"faraday-corr/internal/faraday"
)

// ToInput converts a loaded spec into a faraday.Input, first completing
// the material fields from the registry when a metal name is given. The
// numeric fields are then validated by faraday.Validate, so a negative
// i_corr, an explicit zero valence or an explicit zero density surfaces
// here as an error before any computation runs.
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

// value dereferences a pointer field, treating an omitted value as zero.
func value(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// BuildInput is a convenience wrapper for callers that already hold a
// parsed Spec: it forwards to ToInput.
func BuildInput(spec Spec) (faraday.Input, error) {
	return ToInput(spec)
}
