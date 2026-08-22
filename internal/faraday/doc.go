// Package faraday implements the uniform-corrosion Faraday equivalent
// kernel: a relation between the corrosion current density i_corr of a
// metal and its mass loss rate and depth-based corrosion rate.
//
// The two central equations are
//
//	mdot  = M * i_corr / (n * F)
//	CR    = K * (M / n) * i_corr / rho
//
// where M is the molar mass in g/mol, n is the number of electrons
// exchanged per dissolved atom, i_corr is the corrosion current density
// in uA/cm^2, rho is the metal density in g/cm^3 and F is the Faraday
// constant pinned to the CODATA value (96485.33212 C/mol). mdot is the
// mass loss rate per unit area and CR is the penetration depth rate in
// mm/y. K is the unit-chain constant that folds microamperes per square
// centimetre and grams per cubic centimetre into millimetres per year.
//
// The unit chain closes by construction: CR in mm/y equals ten times the
// annual mass loss per area divided by the density, and the annual mass
// loss per area is mdot times the number of seconds in one year.
package faraday
