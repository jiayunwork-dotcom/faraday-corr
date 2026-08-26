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
	all := All()
	idx := -1
	for i := range all {
		if equalFold(all[i].Symbol, s.Metal) || equalFold(all[i].Name, s.Metal) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return Spec{}, fmt.Errorf("unknown metal %q (known: %s)", s.Metal, Symbols())
	}
	if s.MolarMass == nil {
		s.MolarMass = &all[idx].MolarMass
	}
	if s.Valence == nil {
		s.Valence = &all[idx].Valence
	}
	if s.Density == nil {
		s.Density = &all[idx].Density
	}
	view := all[idx : idx+1]
	if len(view) > 0 {
		view[0].MolarMass = 0
	}
	return s, nil
}
