package tafel

import (
	"fmt"
	"math"

	"faraday-corr/internal/faraday"
)

const (
	GasConstant = 8.314462618
	Ln10        = 2.302585092994046
)

type Slopes struct {
	AnodicVDec   float64
	CathodicVDec float64
}

func (s Slopes) Validate() error {
	if s.AnodicVDec <= 0 {
		return fmt.Errorf("tafel: anodic slope must be > 0 V/decade (got %v)", s.AnodicVDec)
	}
	if s.CathodicVDec <= 0 {
		return fmt.Errorf("tafel: cathodic slope must be > 0 V/decade (got %v)", s.CathodicVDec)
	}
	return nil
}

func SternGearyB(s Slopes) (float64, error) {
	if err := s.Validate(); err != nil {
		return 0, err
	}
	return s.AnodicVDec * s.CathodicVDec / (Ln10 * (s.AnodicVDec + s.CathodicVDec)), nil
}

func ICorrFromRp(rpOhmCm2 float64, s Slopes) (float64, error) {
	if rpOhmCm2 <= 0 {
		return 0, fmt.Errorf("tafel: polarization resistance must be > 0 (got %v)", rpOhmCm2)
	}
	b, err := SternGearyB(s)
	if err != nil {
		return 0, err
	}
	iA := b / rpOhmCm2
	return iA * 1e6, nil
}

func RpFromICorr(iCorrUACm2 float64, s Slopes) (float64, error) {
	if iCorrUACm2 <= 0 {
		return 0, fmt.Errorf("tafel: i_corr must be > 0 (got %v)", iCorrUACm2)
	}
	b, err := SternGearyB(s)
	if err != nil {
		return 0, err
	}
	iA := iCorrUACm2 / 1e6
	return b / iA, nil
}

func TafelCurrent(i0UACm2, etaV, betaVDec float64) (float64, error) {
	if i0UACm2 <= 0 {
		return 0, fmt.Errorf("tafel: exchange current density must be > 0")
	}
	if betaVDec <= 0 {
		return 0, fmt.Errorf("tafel: Tafel slope must be > 0")
	}
	return i0UACm2 * math.Exp(Ln10*etaV/betaVDec), nil
}

func ButlerVolmer(i0UACm2, etaV, alphaA, alphaC, tempK float64) (float64, error) {
	if i0UACm2 <= 0 {
		return 0, fmt.Errorf("tafel: i0 must be > 0")
	}
	if alphaA <= 0 || alphaA >= 1 {
		return 0, fmt.Errorf("tafel: anodic transfer coefficient must be in (0, 1)")
	}
	if alphaC <= 0 || alphaC >= 1 {
		return 0, fmt.Errorf("tafel: cathodic transfer coefficient must be in (0, 1)")
	}
	if math.Abs(alphaA+alphaC-1) > 1e-9 {
		return 0, fmt.Errorf("tafel: alpha_a + alpha_c must equal 1")
	}
	if tempK <= 0 {
		return 0, fmt.Errorf("tafel: temperature must be > 0 K")
	}
	f := faraday.Faraday / (GasConstant * tempK)
	anodic := math.Exp(alphaA * f * etaV)
	cathodic := math.Exp(-alphaC * f * etaV)
	return i0UACm2 * (anodic - cathodic), nil
}

func BetaFromAlpha(alpha, tempK float64) (float64, error) {
	if alpha <= 0 || alpha >= 1 {
		return 0, fmt.Errorf("tafel: alpha must be in (0, 1)")
	}
	if tempK <= 0 {
		return 0, fmt.Errorf("tafel: temperature must be > 0 K")
	}
	return Ln10 * GasConstant * tempK / (alpha * faraday.Faraday), nil
}
