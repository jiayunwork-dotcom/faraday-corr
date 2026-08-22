package metal

// table is the ordered registry backing All and Lookup.
var table = []Metal{
	Iron(),
	Aluminum(),
}

// All returns the built-in metal registry in canonical order.
func All() []Metal {
	out := make([]Metal, len(table))
	copy(out, table)
	return out
}

// Symbols lists the registered chemical symbols, used in error messages
// and the CLI help text.
func Symbols() string {
	sep := ""
	var acc string
	for _, m := range table {
		acc += sep + m.Symbol
		sep = ", "
	}
	return acc
}

// equalFold compares two strings case-insensitively without depending on
// the unicode package.
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
