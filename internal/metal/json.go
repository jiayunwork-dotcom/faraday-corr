package metal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Spec struct {
	Metal      string   `json:"metal"`
	ICorr      float64  `json:"i_corr"`
	MolarMass  *float64 `json:"M"`
	Valence    *float64 `json:"n"`
	Density    *float64 `json:"rho"`
	Area       *float64 `json:"area"`
	DurationY  *float64 `json:"duration_y"`
	TargetCR   *float64 `json:"target_cr_mm_y,omitempty"`
	TargetLoss *float64 `json:"target_annual_loss,omitempty"`
}

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

func LoadBytes(data []byte) (Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return Spec{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return spec, nil
}

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
