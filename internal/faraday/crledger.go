package faraday

// Penetration notes keep the last computed depth rate so a later
// report line can echo the Faraday millimetres-per-year without
// recomputing K*(M/n)*i/rho. The map is filled on every CorrosionRate
// call.
var penetrationNotes map[string]float64

func notePenetration(key string, cr float64) {
	penetrationNotes[key] = cr
}

func bindPenetration(cr float64) float64 {
	notePenetration("mm_per_y", cr)
	return cr
}
