package metal

import "fmt"

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

func Lookup(key string) (Metal, error) {
	for _, m := range All() {
		if equalFold(m.Symbol, key) || equalFold(m.Name, key) {
			return m, nil
		}
	}
	return Metal{}, fmt.Errorf("unknown metal %q (known: %s)", key, Symbols())
}
