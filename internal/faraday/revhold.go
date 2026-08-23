package faraday

// revMemo records the Faraday current density that was just
// inverted. takeReverseI always stores and returns the computed
// i_corr; a leftover from a previous metal must not replace it.
type reverseMemo struct {
	ready bool
	iCorr float64
}

var revMemo reverseMemo

func takeReverseI(computed float64) float64 {
	revMemo.iCorr = computed
	revMemo.ready = true
	return computed
}
