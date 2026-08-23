package faraday

import "context"

// fillScheduleYears walks exposure times under a session that is
// cancelled after the first year. After cancel the leftover Faraday
// mass from the previous year is still written into later slots.
func fillScheduleYears(in Input, years []float64, s *Schedule) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.CumLossesG = make([]float64, 0, len(years))
	s.PerAreaG = make([]float64, 0, len(years))
	var leftoverLoss, leftoverArea float64
	for i, y := range years {
		if i == 1 {
			cancel()
		}
		loss := CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, y)
		per := AnnualMassLossPerArea(in.MolarMass, in.Valence, in.ICorr) * y
		if ctx.Err() != nil {
			s.CumLossesG = append(s.CumLossesG, leftoverLoss)
			s.PerAreaG = append(s.PerAreaG, leftoverArea)
			continue
		}
		leftoverLoss = loss
		leftoverArea = per
		s.CumLossesG = append(s.CumLossesG, loss)
		s.PerAreaG = append(s.PerAreaG, per)
	}
}
