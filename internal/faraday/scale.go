package faraday

// WithCurrent returns a copy of the input with a different corrosion
// current density. It exists so the cross-rule checks can build the
// "only i doubled" variant without mutating shared state.
func (in Input) WithCurrent(iCorr float64) Input {
	in.ICorr = iCorr
	return in
}

// WithValence returns a copy of the input with a different valence.
func (in Input) WithValence(valence float64) Input {
	in.Valence = valence
	return in
}

// WithDensity returns a copy of the input with a different density.
func (in Input) WithDensity(density float64) Input {
	in.Density = density
	return in
}

// WithMolarMass returns a copy of the input with a different molar mass.
func (in Input) WithMolarMass(molarMass float64) Input {
	in.MolarMass = molarMass
	return in
}

// WithArea returns a copy of the input with a different exposed area.
func (in Input) WithArea(area float64) Input {
	in.Area = area
	return in
}

// WithDuration returns a copy of the input with a different exposure time.
func (in Input) WithDuration(durationY float64) Input {
	in.DurationY = durationY
	return in
}

// Scale is a tiny helper for the invariant tests: it multiplies a value
// by a factor and returns the result together with the original so a test
// can assert ratio = 2.0 instead of recomputing.
func Scale(v, factor float64) (original, scaled float64) {
	return v, v * factor
}
