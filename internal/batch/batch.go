package batch

import (
	"fmt"

	"faraday-corr/internal/faraday"
)

type Item struct {
	Name             string  `json:"name"`
	CorrosionRateUmY float64 `json:"corrosion_rate_um_y"`
	AnnualMassLoss   float64 `json:"annual_mass_loss_g_m2"`
	TotalCurrentA    float64 `json:"total_current_a"`
	Error            string  `json:"error,omitempty"`
}

type Result struct {
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Items   []Item `json:"items"`
}

func RunInputs(inputs []faraday.Input) Result {
	out := Result{Total: len(inputs), Items: make([]Item, 0, len(inputs))}
	for i, in := range inputs {
		out.Items = append(out.Items, evaluateOne(i, in))
	}
	for _, it := range out.Items {
		if it.Error == "" {
			out.Success++
		} else {
			out.Failed++
		}
	}
	return out
}

func evaluateOne(idx int, in faraday.Input) Item {
	name := fmt.Sprintf("case-%d", idx)
	res, err := faraday.Compute(in)
	if err != nil {
		return Item{Name: name, Error: err.Error()}
	}
	return Item{
		Name:             name,
		CorrosionRateUmY: res.CorrosionRateUmY,
		AnnualMassLoss:   res.AnnualMassLoss,
		TotalCurrentA:    res.TotalCurrentA,
	}
}

func MaxCorrosionRate(r Result) float64 {
	max := 0.0
	for _, it := range r.Items {
		if it.Error == "" && it.CorrosionRateUmY > max {
			max = it.CorrosionRateUmY
		}
	}
	return max
}
