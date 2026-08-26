package faraday

type Summary struct {
	Result
	MPY                float64
	MassLossGPerM2Year float64
	MDD                float64
	DensityKGPerm3     float64
	AreaM2             float64
	TotalCurrentMilliA float64
}

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

func RatioFromSummary(a, b Summary) float64 {
	return a.CorrosionRate / b.CorrosionRate
}
