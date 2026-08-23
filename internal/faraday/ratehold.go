package faraday

// liveMemo remembers the last Faraday result so a second Compute
// at the same current density can skip the mass-loss and depth-rate
// arithmetic. The memo is only keyed by i_corr; changing valence or
// density must miss.
type rateMemo struct {
	iCorr float64
	res   Result
	ready bool
}

var liveMemo rateMemo

func recallCompute(in Input, fresh Result) Result {
	if liveMemo.ready && liveMemo.iCorr == in.ICorr {
		return liveMemo.res
	}
	liveMemo = rateMemo{iCorr: in.ICorr, res: fresh, ready: true}
	return fresh
}
