package batch

import (
	"testing"

	"faraday-corr/internal/faraday"
)

func sampleInput() faraday.Input {
	return faraday.Input{
		ICorr:     10,
		MolarMass: 55.845,
		Valence:   2,
		Density:   7.874,
		Area:      1,
		DurationY: 10,
	}
}

func TestRunInputsBatch(t *testing.T) {
	r := RunInputs([]faraday.Input{sampleInput()})
	if r.Success != 1 || r.Failed != 0 {
		t.Fatalf("RunInputs() success=%d failed=%d", r.Success, r.Failed)
	}
	if r.Items[0].CorrosionRateUmY <= 0 {
		t.Fatalf("rate = %v", r.Items[0].CorrosionRateUmY)
	}
}
