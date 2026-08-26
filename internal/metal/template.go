package metal

import (
	"fmt"
	"strings"
)

func TemplateSpec(metalName string) (Spec, error) {
	if metalName == "" {
		metalName = "Fe"
	}
	m, err := Lookup(metalName)
	if err != nil {
		return Spec{}, err
	}
	area := 1.0
	duration := 1.0
	return Spec{
		Metal:     m.Symbol,
		ICorr:     10.0,
		MolarMass: &m.MolarMass,
		Valence:   &m.Valence,
		Density:   &m.Density,
		Area:      &area,
		DurationY: &duration,
	}, nil
}

func HeaderComment() string {
	return strings.Join([]string{
		"i_corr     uA/cm^2 (corrosion current density)",
		"M          g/mol   (molar mass)",
		"n          electrons per dissolved atom",
		"rho        g/cm^3  (density)",
		"area       cm^2    (optional exposed area)",
		"duration_y years   (optional exposure time)",
		"metal      symbol or name from the built-in registry;",
		"           fills M, n, rho when those fields are omitted",
	}, "\n")
}

func DescribeUnits() string {
	return "units: i_corr uA/cm^2, M g/mol, n electrons, rho g/cm^3, area cm^2, duration_y years"
}

func FormatError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w (check the unit conventions: %s)", err, DescribeUnits())
}
