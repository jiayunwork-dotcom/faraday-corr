package faraday

type Input struct {
	ICorr     float64 `json:"i_corr"`
	MolarMass float64 `json:"M"`
	Valence   float64 `json:"n"`
	Density   float64 `json:"rho"`
	Area      float64 `json:"area"`
	DurationY float64 `json:"duration_y"`
}

type Result struct {
	Input

	MassLossRate       float64
	AnnualMassLoss     float64
	CorrosionRate      float64
	CorrosionRateUmY   float64
	TotalCurrent       float64
	TotalCurrentA      float64
	CumulativeMassLoss float64
}
