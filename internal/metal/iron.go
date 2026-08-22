package metal

import "fmt"

// Iron returns the built-in iron entry. Iron is assumed to dissolve as
// Fe -> Fe^2+ + 2e, the standard anodic reaction in aerated seawater.
func Iron() Metal {
	return Metal{
		Symbol:    "Fe",
		Name:      "iron",
		MolarMass: 55.845,
		Valence:   2,
		Density:   7.874,
		Note:      "Fe -> Fe2+ + 2e-",
	}
}

// Aluminum returns the built-in aluminium entry. Aluminium dissolves as
// Al -> Al^3+ + 3e; its low density (2.70 g/cm^3) versus iron is exactly
// why, at equal current density, the penetration depth rate differs.
func Aluminum() Metal {
	return Metal{
		Symbol:    "Al",
		Name:      "aluminium",
		MolarMass: 26.9815,
		Valence:   3,
		Density:   2.70,
		Note:      "Al -> Al3+ + 3e-",
	}
}

// Lookup returns the registry entry matching a symbol or name. The
// comparison is case-insensitive so "FE", "fe" and "Fe" all resolve.
func Lookup(key string) (Metal, error) {
	for _, m := range All() {
		if equalFold(m.Symbol, key) || equalFold(m.Name, key) {
			return m, nil
		}
	}
	return Metal{}, fmt.Errorf("unknown metal %q (known: %s)", key, Symbols())
}
