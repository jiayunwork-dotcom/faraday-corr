package faraday

// fillScheduleYears computes the cumulative mass loss for each exposure
// time in years under steady uniform corrosion. For a fixed current
// density the Faraday loss is linear in time, so each slot is computed
// independently from the full exposure duration.
func fillScheduleYears(in Input, years []float64, s *Schedule) {
	s.CumLossesG = make([]float64, 0, len(years))
	s.PerAreaG = make([]float64, 0, len(years))
	for _, y := range years {
		loss := CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, y)
		per := AnnualMassLossPerArea(in.MolarMass, in.Valence, in.ICorr) * y
		s.CumLossesG = append(s.CumLossesG, loss)
		s.PerAreaG = append(s.PerAreaG, per)
	}
}
