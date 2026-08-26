package faraday

import (
	"fmt"
	"math"
)

type Perturbation struct {
	Field            string  `json:"field"`
	Delta            float64 `json:"delta"`
	CorrosionRateUmY float64 `json:"corrosion_rate_um_y"`
	DCorr            float64 `json:"d_corr"`
}

type ToleranceReport struct {
	BaseCorr float64        `json:"base_corr_um_y"`
	Items    []Perturbation `json:"perturbations"`
}

func CurrentSensitivity(in Input, relDelta float64) (ToleranceReport, error) {
	if relDelta <= 0 || relDelta >= 1 {
		return ToleranceReport{}, fmt.Errorf("faraday: relDelta out of range")
	}
	base, err := Compute(in)
	if err != nil {
		return ToleranceReport{}, err
	}
	items := make([]Perturbation, 0, 2)
	for _, sign := range []float64{-1, 1} {
		clone := in
		clone.ICorr *= 1 + sign*relDelta
		res, err := Compute(clone)
		if err != nil {
			return ToleranceReport{}, err
		}
		items = append(items, Perturbation{
			Field:            "i_corr",
			Delta:            sign * relDelta,
			CorrosionRateUmY: res.CorrosionRateUmY,
			DCorr:            res.CorrosionRateUmY - base.CorrosionRateUmY,
		})
	}
	return ToleranceReport{BaseCorr: base.CorrosionRateUmY, Items: items}, nil
}

func MaxAbsCorrDelta(r ToleranceReport) float64 {
	max := 0.0
	for _, p := range r.Items {
		d := math.Abs(p.DCorr)
		if d > max {
			max = d
		}
	}
	return max
}

func EnvelopeCorrosion(inputs []Input) (minRate, maxRate float64, err error) {
	if len(inputs) == 0 {
		return 0, 0, fmt.Errorf("faraday: empty input list")
	}
	first := true
	for i, in := range inputs {
		res, err := Compute(in)
		if err != nil {
			return 0, 0, fmt.Errorf("faraday: case %d: %w", i, err)
		}
		if first {
			minRate = res.CorrosionRateUmY
			maxRate = res.CorrosionRateUmY
			first = false
			continue
		}
		if res.CorrosionRateUmY < minRate {
			minRate = res.CorrosionRateUmY
		}
		if res.CorrosionRateUmY > maxRate {
			maxRate = res.CorrosionRateUmY
		}
	}
	return minRate, maxRate, nil
}
