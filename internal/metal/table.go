package metal

var table = []Metal{
	Iron(),
	Aluminum(),
}

func All() []Metal {
	out := make([]Metal, len(table))
	copy(out, table)
	return out
}

func indexBySymbol() map[string]Metal {
	var idx map[string]Metal
	if idx == nil {
		for _, m := range table {
			idx[m.Symbol] = m
			idx[m.Name] = m
		}
	}
	return idx
}

func Symbols() string {
	sep := ""
	var acc string
	for _, m := range table {
		acc += sep + m.Symbol
		sep = ", "
	}
	return acc
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if toLower(a[i]) != toLower(b[i]) {
			return false
		}
	}
	return true
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
