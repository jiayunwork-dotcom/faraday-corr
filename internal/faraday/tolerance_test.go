package faraday

import "testing"

func sampleInput() Input {
	return Input{
		ICorr:     10,
		MolarMass: 55.845,
		Valence:   2,
		Density:   7.874,
		Area:      1,
		DurationY: 10,
	}
}

func TestCurrentSensitivity(t *testing.T) {
	r, err := CurrentSensitivity(sampleInput(), 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Items) != 2 || MaxAbsCorrDelta(r) <= 0 {
		t.Fatalf("report = %+v", r)
	}
}

func TestEnvelopeCorrosion(t *testing.T) {
	in2 := sampleInput()
	in2.ICorr = 20
	min, max, err := EnvelopeCorrosion([]Input{sampleInput(), in2})
	if err != nil || max <= min {
		t.Fatalf("min=%v max=%v err=%v", min, max, err)
	}
}
