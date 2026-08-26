package faraday

func (in Input) WithCurrent(iCorr float64) Input {
	in.ICorr = iCorr
	return in
}

func (in Input) WithValence(valence float64) Input {
	in.Valence = valence
	return in
}

func (in Input) WithDensity(density float64) Input {
	in.Density = density
	return in
}

func (in Input) WithMolarMass(molarMass float64) Input {
	in.MolarMass = molarMass
	return in
}

func (in Input) WithArea(area float64) Input {
	in.Area = area
	return in
}

func (in Input) WithDuration(durationY float64) Input {
	in.DurationY = durationY
	return in
}

func Scale(v, factor float64) (original, scaled float64) {
	return v, v * factor
}
