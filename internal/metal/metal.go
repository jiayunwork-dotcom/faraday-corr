package metal

// Metal is a named entry in the built-in registry. The unit conventions
// match faraday.Input: molar mass in g/mol, valence as a dimensionless
// number of electrons, density in g/cm^3.
type Metal struct {
	Symbol    string  // chemical symbol, e.g. "Fe"
	Name      string  // descriptive name, e.g. "iron"
	MolarMass float64 // g/mol
	Valence   float64 // electrons per dissolved atom
	Density   float64 // g/cm^3
	Note      string  // dissolution half-reaction or typical environment
}

// EquivalentWeight returns M / n in g/mol.
func (m Metal) EquivalentWeight() float64 {
	return m.MolarMass / m.Valence
}

// String gives a compact one-line description used by the `metals`
// listing subcommand.
func (m Metal) String() string {
	return m.Symbol + " " + m.Name + " M=" + trimZeros(m.MolarMass) + " n=" + trimZeros(m.Valence) + " rho=" + trimZeros(m.Density)
}
