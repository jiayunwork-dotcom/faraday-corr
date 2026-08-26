package faraday

const MMPerYearPerMPY = 0.0254

const MddPerGPerCm2Year = 0.00365

func MPY(mmPerYear float64) float64 {
	return mmPerYear / MMPerYearPerMPY
}

func MMPerYearFromMPY(mpy float64) float64 {
	return mpy * MMPerYearPerMPY
}

func GPerM2Year(gPerCm2Year float64) float64 {
	return gPerCm2Year * 1e4
}

func Mdd(gPerCm2Year float64) float64 {
	return gPerCm2Year / MddPerGPerCm2Year
}

func GPerCm2YearFromMdd(mdd float64) float64 {
	return mdd * MddPerGPerCm2Year
}

func KGPerm3(gPerCm3 float64) float64 {
	return gPerCm3 * 1000
}

func Cm2ToM2(cm2 float64) float64 {
	return cm2 / 1e4
}

func M2ToCm2(m2 float64) float64 {
	return m2 * 1e4
}

func MPYFromCorrosionRate(crMMPerYear float64) float64 {
	return MPY(crMMPerYear)
}

func AnnualMassLossPerM2(annualPerCm2 float64) float64 {
	return GPerM2Year(annualPerCm2)
}
