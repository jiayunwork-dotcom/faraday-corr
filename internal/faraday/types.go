package faraday

// Input carries the material and electrochemical quantities needed for a
// single uniform-corrosion calculation. Every value is validated by
// Validate before any computation: i_corr, M, n and rho must be strictly
// positive, and the optional area and duration must not be negative.
//
// Units are fixed so that the unit chain closes:
//
//	i_corr     uA/cm^2
//	M          g/mol
//	rho        g/cm^3
//	area       cm^2
//	duration_y years
type Input struct {
	ICorr     float64 `json:"i_corr"`     // corrosion current density, uA/cm^2
	MolarMass float64 `json:"M"`          // molar mass, g/mol
	Valence   float64 `json:"n"`          // electrons exchanged per dissolved atom
	Density   float64 `json:"rho"`        // density, g/cm^3
	Area      float64 `json:"area"`       // exposed area, cm^2; 0 means "not given"
	DurationY float64 `json:"duration_y"` // exposure time, years; 0 means "not given"
}

// Result is the outcome of a Compute call. It repeats the input values
// and adds every derived quantity. All rates are per unit exposed area;
// the cumulative mass loss additionally depends on the area and duration
// when those are provided.
type Result struct {
	Input

	// MassLossRate is mdot = M * i_corr / (n * F) in g/(cm^2 s).
	MassLossRate float64
	// AnnualMassLoss is the mass lost per square centimetre per year.
	AnnualMassLoss float64
	// CorrosionRate is the penetration depth rate in mm/y.
	CorrosionRate float64
	// CorrosionRateUmY is the penetration depth rate in um/y.
	CorrosionRateUmY float64
	// TotalCurrent is the aggregate current I = i_corr * area in uA.
	TotalCurrent float64
	// TotalCurrentA is the aggregate current I = i_corr * area in A.
	TotalCurrentA float64
	// CumulativeMassLoss is the total mass lost over the exposed area
	// and duration, in grams. It is only meaningful when area and
	// duration were given.
	CumulativeMassLoss float64
}
