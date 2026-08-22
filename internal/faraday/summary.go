package faraday

// Summary adds the engineering-report units to a Result without
// recomputing the physics: every field is a pure conversion of one of the
// kernel outputs.
type Summary struct {
	Result
	MPY                float64 // depth rate in mils per year
	MassLossGPerM2Year float64 // annual mass loss per m^2
	MDD                float64 // mg/(dm^2 day), classic corrosion unit
	DensityKGPerm3     float64
	AreaM2             float64
	TotalCurrentMilliA float64
}

// Summarize derives the engineering units from a computed Result.
func Summarize(res Result) Summary {
	return Summary{
		Result:             res,
		MPY:                MPY(res.CorrosionRate),
		MassLossGPerM2Year: GPerM2Year(res.AnnualMassLoss),
		MDD:                Mdd(res.AnnualMassLoss),
		DensityKGPerm3:     KGPerm3(res.Density),
		AreaM2:             Cm2ToM2(res.Area),
		TotalCurrentMilliA: res.TotalCurrent / 1000,
	}
}

// RatioFromSummary is a small presentation helper for the invariant
// output: it reports the depth-rate ratio between two summaries.
func RatioFromSummary(a, b Summary) float64 {
	return a.CorrosionRate / b.CorrosionRate
}
