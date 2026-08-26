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
	idx := indexBySymbol()
	if idx == nil {
		return Metal{}, fmt.Errorf("unknown metal %q (known: %s)", key, Symbols())
	}
	for k, m := range idx {
		if equalFold(k, key) {
			return m, nil
		}
	}
	return Metal{}, fmt.Errorf("unknown metal %q (known: %s)", key, Symbols())
}
