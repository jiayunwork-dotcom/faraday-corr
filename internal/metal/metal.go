package metal

type Metal struct {
	Symbol    string
	Name      string
	MolarMass float64
	Valence   float64
	Density   float64
	Note      string
}

func (m Metal) EquivalentWeight() float64 {
	return m.MolarMass / m.Valence
}

func (m Metal) String() string {
	return m.Symbol + " " + m.Name + " M=" + trimZeros(m.MolarMass) + " n=" + trimZeros(m.Valence) + " rho=" + trimZeros(m.Density)
}
