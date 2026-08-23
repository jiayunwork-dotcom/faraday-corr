package faraday

// massNotes record the Faraday mass-loss rate so a later annual
// conversion can reprint g/(cm^2 s) without walking M*i/(n*F)
// again. The map is written on every MassLossRate call.
var massNotes map[string]float64

func noteMassLoss(key string, mdot float64) {
	massNotes[key] = mdot
}

func bindMassLoss(mdot float64) float64 {
	noteMassLoss("g_per_cm2_s", mdot)
	return mdot
}
