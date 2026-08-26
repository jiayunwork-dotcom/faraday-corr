package faraday

import (
	"fmt"
	"math"
)

type Couple struct {
	ICorrA float64
	ICorrC float64
	AreaA  float64
	AreaC  float64
	BetaA  float64
	BetaC  float64
}

func (c Couple) Validate() error {
	if c.ICorrA <= 0 || c.ICorrC <= 0 {
		return fmt.Errorf("faraday: both corrosion currents must be positive")
	}
	if c.AreaA <= 0 || c.AreaC <= 0 {
		return fmt.Errorf("faraday: both areas must be positive")
	}
	if c.BetaA <= 0 || c.BetaC <= 0 {
		return fmt.Errorf("faraday: both Tafel slopes must be positive")
	}
	return nil
}

func MixedCurrent(c Couple) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	ia := c.ICorrA * c.AreaA
	ic := c.ICorrC * c.AreaC
	ba := c.BetaA
	bc := c.BetaC
	logTerm := math.Log10((ic / ia) * (ba / bc))
	eta := (ba * bc) / (ba + bc) * logTerm
	iMix := ia * math.Pow(10, eta/ba)
	if iMix <= 0 || math.IsNaN(iMix) || math.IsInf(iMix, 0) {
		return 0, fmt.Errorf("faraday: mixed current not finite")
	}
	return iMix, nil
}

func CathodeAreaRaisesCurrent(base Couple, areaScale float64) error {
	if areaScale <= 1 {
		return fmt.Errorf("faraday: areaScale must be > 1")
	}
	i0, err := MixedCurrent(base)
	if err != nil {
		return err
	}
	up := base
	up.AreaC *= areaScale
	i1, err := MixedCurrent(up)
	if err != nil {
		return err
	}
	if i1 <= i0 {
		return fmt.Errorf("faraday: larger cathode did not raise mixed current %g → %g", i0, i1)
	}
	return nil
}

func WagnerNumber(conductivity, rpOhm, lengthM float64) (float64, error) {
	if conductivity <= 0 || rpOhm <= 0 || lengthM <= 0 {
		return 0, fmt.Errorf("faraday: conductivity, Rp and length must be positive")
	}
	return conductivity * rpOhm / lengthM, nil
}

func UniformWhenWagnerLarge(we, threshold float64) bool {
	return we >= threshold
}
