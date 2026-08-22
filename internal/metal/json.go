package metal

import (
	"encoding/json"
	"fmt"
	"os"
)

// Spec mirrors the JSON input file accepted by the `rate` subcommand.
// Field names and units follow the specification: i_corr in uA/cm^2,
// M in g/mol, n in electrons per atom, rho in g/cm^3, area in cm^2 and
// duration_y in years. The numeric fields are pointers so the loader can
// tell "omitted" (nil, filled from the registry when a metal is named)
// from "explicitly zero" (kept, and rejected by validation for M, n, rho).
// The metal field is a convenience label: when a registered symbol or
// name is given, any of M, n and rho that was omitted is filled from the
// registry entry. Explicit values always win.
type Spec struct {
	Metal     string   `json:"metal"`
	ICorr     float64  `json:"i_corr"`
	MolarMass *float64 `json:"M"`
	Valence   *float64 `json:"n"`
	Density   *float64 `json:"rho"`
	Area      *float64 `json:"area"`
	DurationY *float64 `json:"duration_y"`
	// TargetCR and TargetLoss are optional reverse-solve targets: the
	// `reverse` subcommand inverts the kernel so that the solved i_corr
	// produces the requested depth rate (mm/y) or annual mass loss
	// (g/(cm^2 y)). At most one may be set.
	TargetCR   *float64 `json:"target_cr_mm_y,omitempty"`
	TargetLoss *float64 `json:"target_annual_loss,omitempty"`
}

// LoadFile reads and decodes a JSON specification file.
func LoadFile(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, fmt.Errorf("cannot read %s: %w", path, err)
	}
	spec, err := LoadBytes(data)
	if err != nil {
		return Spec{}, fmt.Errorf("%s: %w", path, err)
	}
	return spec, nil
}

// LoadBytes decodes a JSON specification from raw bytes.
func LoadBytes(data []byte) (Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return Spec{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return spec, nil
}

// FillFromRegistry completes the material fields from the built-in table
// when the spec names a metal and leaves M, n or rho omitted.
func (s Spec) FillFromRegistry() (Spec, error) {
	if s.Metal == "" {
		return s, nil
	}
	m, err := Lookup(s.Metal)
	if err != nil {
		return Spec{}, err
	}
	if s.MolarMass == nil {
		v := m.MolarMass
		s.MolarMass = &v
	}
	if s.Valence == nil {
		v := m.Valence
		s.Valence = &v
	}
	if s.Density == nil {
		v := m.Density
		s.Density = &v
	}
	return s, nil
}
