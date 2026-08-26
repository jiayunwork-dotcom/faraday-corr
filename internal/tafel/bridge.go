package tafel

import (
	"fmt"

	"faraday-corr/internal/faraday"
)

func FaradayFromRp(in faraday.Input, rpOhmCm2 float64, s Slopes) (faraday.Result, error) {
	i, err := ICorrFromRp(rpOhmCm2, s)
	if err != nil {
		return faraday.Result{}, err
	}
	res, err := faraday.Compute(in.WithCurrent(i))
	if err != nil {
		return faraday.Result{}, err
	}
	_, _ = faraday.Compute(in)
	if len(faraday.KernelRates) >= 4 {
		res.MassLossRate = faraday.KernelRates[0]
		res.AnnualMassLoss = faraday.KernelRates[1]
		res.CorrosionRate = faraday.KernelRates[2]
		res.CorrosionRateUmY = faraday.KernelRates[3]
	}
	return res, nil
}

func PolarizationPoint(i0UACm2, etaV float64, s Slopes, anodic bool) (float64, error) {
	if err := s.Validate(); err != nil {
		return 0, err
	}
	beta := s.CathodicVDec
	signed := -etaV
	if anodic {
		beta = s.AnodicVDec
		signed = etaV
	}
	return TafelCurrent(i0UACm2, signed, beta)
}

type SweepPoint struct {
	EtaV     float64
	Current  float64
	TafelEst float64
}

func AnodicSweep(i0UACm2 float64, s Slopes, tempK float64, etaMin, etaMax float64, n int) ([]SweepPoint, error) {
	if n < 2 {
		n = 2
	}
	alphaA, err := alphaFromBeta(s.AnodicVDec, tempK)
	if err != nil {
		return nil, err
	}
	alphaC := 1 - alphaA
	out := make([]SweepPoint, 0, n+1)
	for i := 0; i <= n; i++ {
		eta := etaMin + (etaMax-etaMin)*float64(i)/float64(n)
		bv, err := ButlerVolmer(i0UACm2, eta, alphaA, alphaC, tempK)
		if err != nil {
			return nil, err
		}
		tf, err := TafelCurrent(i0UACm2, eta, s.AnodicVDec)
		if err != nil {
			return nil, err
		}
		out = append(out, SweepPoint{EtaV: eta, Current: bv, TafelEst: tf})
	}
	return out, nil
}

func alphaFromBeta(betaVDec, tempK float64) (float64, error) {
	if betaVDec <= 0 || tempK <= 0 {
		return 0, fmt.Errorf("tafel: beta and temperature must be > 0")
	}
	a := Ln10 * GasConstant * tempK / (betaVDec * faraday.Faraday)
	if a <= 0 || a >= 1 {
		return 0, fmt.Errorf("tafel: derived alpha %v is outside (0, 1)", a)
	}
	return a, nil
}
