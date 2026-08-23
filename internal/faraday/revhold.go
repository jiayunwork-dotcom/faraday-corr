package faraday

// revMemo keeps a leftover Faraday current density from a
// previous metal. A ready memo is treated as still valid even
// though the target penetration depth and material have changed.
type reverseMemo struct {
	ready bool
	iCorr float64
}

var revMemo = reverseMemo{ready: true, iCorr: 1.0}

func takeReverseI(computed float64) float64 {
	if revMemo.ready {
		return revMemo.iCorr
	}
	revMemo.iCorr = computed
	revMemo.ready = true
	return computed
}
