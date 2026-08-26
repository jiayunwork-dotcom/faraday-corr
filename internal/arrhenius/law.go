package arrhenius

import (
	"fmt"
	"math"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/tafel"
)

type Law struct {
	IRefUACm2 float64
	TRefK     float64
	EaJMol    float64
}

func (l Law) Validate() error {
	if l.IRefUACm2 <= 0 {
		return fmt.Errorf("arrhenius: reference i_corr must be > 0")
	}
	if l.TRefK <= 0 {
		return fmt.Errorf("arrhenius: reference temperature must be > 0 K")
	}
	if l.EaJMol <= 0 {
		return fmt.Errorf("arrhenius: activation energy must be > 0")
	}
	return nil
}

func (l Law) ICorrAt(tempK float64) (float64, error) {
	if err := l.Validate(); err != nil {
		return 0, err
	}
	if tempK <= 0 {
		return 0, fmt.Errorf("arrhenius: temperature must be > 0 K")
	}
	inv := 1/tempK - 1/l.TRefK
	return l.IRefUACm2 * math.Exp(-l.EaJMol/tafel.GasConstant*inv), nil
}

func (l Law) Q10(tempK float64) (float64, error) {
	i0, err := l.ICorrAt(tempK)
	if err != nil {
		return 0, err
	}
	i1, err := l.ICorrAt(tempK + 10)
	if err != nil {
		return 0, err
	}
	return i1 / i0, nil
}

func FaradayAtTemperature(in faraday.Input, l Law, tempK float64) (faraday.Result, float64, error) {
	i, err := l.ICorrAt(tempK)
	if err != nil {
		return faraday.Result{}, 0, err
	}
	res, err := faraday.Compute(in.WithCurrent(i))
	if err != nil {
		return faraday.Result{}, 0, err
	}
	return res, i, nil
}

func ScaleFactor(l Law, fromK, toK float64) (float64, error) {
	a, err := l.ICorrAt(fromK)
	if err != nil {
		return 0, err
	}
	b, err := l.ICorrAt(toK)
	if err != nil {
		return 0, err
	}
	return b / a, nil
}
