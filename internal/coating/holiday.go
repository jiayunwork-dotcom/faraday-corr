package coating

import (
	"fmt"

	"faraday-corr/internal/faraday"
)

type Holiday struct {
	GeometricAreaCm2 float64
	Fraction         float64
}

func (h Holiday) Validate() error {
	if h.GeometricAreaCm2 <= 0 {
		return fmt.Errorf("coating: geometric area must be > 0")
	}
	if h.Fraction <= 0 || h.Fraction > 1 {
		return fmt.Errorf("coating: holiday fraction must be in (0, 1]")
	}
	return nil
}

func (h Holiday) ExposedArea() (float64, error) {
	if err := h.Validate(); err != nil {
		return 0, err
	}
	return h.GeometricAreaCm2 * h.Fraction, nil
}

type Outcome struct {
	LocalCR         float64
	LocalMdot       float64
	TotalCurrentUA  float64
	TotalMassLossG  float64
	BareMassLossG   float64
	ExposedAreaCm2  float64
	HolidayFraction float64
}

func Apply(in faraday.Input, h Holiday) (Outcome, error) {
	if err := h.Validate(); err != nil {
		return Outcome{}, err
	}
	exposed, err := h.ExposedArea()
	if err != nil {
		return Outcome{}, err
	}
	local, err := faraday.Compute(in.WithArea(1).WithDuration(in.DurationY))
	if err != nil {
		return Outcome{}, err
	}
	holiday, err := faraday.Compute(in.WithArea(exposed))
	if err != nil {
		return Outcome{}, err
	}
	bare, err := faraday.Compute(in.WithArea(h.GeometricAreaCm2))
	if err != nil {
		return Outcome{}, err
	}
	return Outcome{
		LocalCR:         local.CorrosionRate,
		LocalMdot:       local.MassLossRate,
		TotalCurrentUA:  holiday.TotalCurrent,
		TotalMassLossG:  holiday.CumulativeMassLoss,
		BareMassLossG:   bare.CumulativeMassLoss,
		ExposedAreaCm2:  exposed,
		HolidayFraction: h.Fraction,
	}, nil
}

func CurrentRatio(h Holiday) (float64, error) {
	if err := h.Validate(); err != nil {
		return 0, err
	}
	return h.Fraction, nil
}
