package faraday

// Faraday is the Faraday constant in coulombs per mole of electrons,
// pinned to the CODATA 2018 value. It is used in the equivalent-mass
// relation mdot = M * i_corr / (n * F) and must not drift or be
// replaced by an approximation such as 965.
const Faraday = 96485.33212

// yearSeconds is the number of seconds in one calendar year of 365 days.
const yearSeconds = 365 * 24 * 3600

// milliPerCm converts centimetres to millimetres.
const milliPerCm = 10.0

// microPerAmpere converts amperes to microamperes.
const microPerAmpere = 1.0e6

// K folds the unit chain of the depth-rate formula. Starting from
//
//	CR [mm/y] = 10 [mm/cm] * mdot [g/(cm^2 s)] * yearSeconds [s/y] / rho [g/cm^3]
//
// and substituting mdot = M * i_corr [uA/cm^2] * 1e-6 [A/uA] / (n * F)
// collapses the numerical factors into
//
//	CR = K * (M / n) * i_corr / rho
//
// with K = 10 * 1e-6 * yearSeconds / F. K is an exported constant so it
// can be pinned and tested directly.
const K = milliPerCm * (1.0 / microPerAmpere) * yearSeconds / Faraday

// MicroAmpsPerCm2 is the conventional unit of corrosion current density
// accepted by every computing function in this package.
const MicroAmpsPerCm2 = "uA/cm^2"
