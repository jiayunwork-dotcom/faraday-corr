package faraday

// Engineering-unit conversions. The specification demands that the unit
// chain close for mm/y and um/y; corrosion engineers also report depth
// rates in mils per year (mpy) and mass loss in grams per square metre
// per year or milligrams per square decimetre per day (mdd). Every
// conversion here is a plain factor from the base units used by the
// kernel.

// MMPerYearPerMPY is the millimetres-per-year equivalent of one mil per
// year. A mil is one thousandth of an inch (0.0254 mm).
const MMPerYearPerMPY = 0.0254

// MddPerGPerCm2Year is the conversion factor from g/(cm^2 y) to mdd
// (mg/(dm^2 day)). One mdd is 1e-3 g over 100 cm^2 in 1/365 of a year.
const MddPerGPerCm2Year = 0.00365

// MPY converts a depth rate from mm/y to mils per year.
func MPY(mmPerYear float64) float64 {
	return mmPerYear / MMPerYearPerMPY
}

// MMPerYearFromMPY converts a depth rate from mils per year to mm/y.
func MMPerYearFromMPY(mpy float64) float64 {
	return mpy * MMPerYearPerMPY
}

// GPerM2Year converts a mass loss rate from g/(cm^2 y) to g/(m^2 y).
// One square metre holds 1e4 square centimetres, so the number scales by
// 1e4.
func GPerM2Year(gPerCm2Year float64) float64 {
	return gPerCm2Year * 1e4
}

// Mdd converts a mass loss rate from g/(cm^2 y) to milligrams per square
// decimetre per day, the classical corrosion unit.
func Mdd(gPerCm2Year float64) float64 {
	return gPerCm2Year / MddPerGPerCm2Year
}

// GPerCm2YearFromMdd is the inverse of Mdd.
func GPerCm2YearFromMdd(mdd float64) float64 {
	return mdd * MddPerGPerCm2Year
}

// KGPerm3 converts density from g/cm^3 to kg/m^3 for reports that mix SI
// density units.
func KGPerm3(gPerCm3 float64) float64 {
	return gPerCm3 * 1000
}

// Cm2ToM2 converts an exposed area from square centimetres to square
// metres.
func Cm2ToM2(cm2 float64) float64 {
	return cm2 / 1e4
}

// M2ToCm2 converts an exposed area from square metres to square
// centimetres.
func M2ToCm2(m2 float64) float64 {
	return m2 * 1e4
}

// MPYFromCorrosionRate is a convenience for callers holding a Result:
// it returns the depth rate in the classic US corrosion unit.
func MPYFromCorrosionRate(crMMPerYear float64) float64 {
	return MPY(crMMPerYear)
}

// AnnualMassLossPerM2 returns the annual mass loss of one square metre
// of the same material, the metric counterpart of the per-cm^2 figure.
func AnnualMassLossPerM2(annualPerCm2 float64) float64 {
	return GPerM2Year(annualPerCm2)
}
