package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
)

func TestWriteChain(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874}
	c, err := faraday.VerifyUnitChain(in)
	if err != nil {
		t.Fatalf("VerifyUnitChain failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteChain(&buf, c); err != nil {
		t.Fatalf("WriteChain failed: %v", err)
	}
	text := buf.String()
	for _, want := range []string{"chain consistent:", "CR direct", "CR from mass path"} {
		if !strings.Contains(text, want) {
			t.Errorf("chain output missing %q:\n%s", want, text)
		}
	}
}

func TestWriteResultJSON(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1}
	res, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteResultJSON(&buf, res); err != nil {
		t.Fatalf("WriteResultJSON failed: %v", err)
	}
	var back faraday.Result
	if err := json.Unmarshal(buf.Bytes(), &back); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if back.CorrosionRate != res.CorrosionRate {
		t.Errorf("round-trip CR = %v, want %v", back.CorrosionRate, res.CorrosionRate)
	}
	if back.MassLossRate != res.MassLossRate {
		t.Errorf("round-trip mdot = %v, want %v", back.MassLossRate, res.MassLossRate)
	}
}

func TestWriteErrorJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteErrorJSON(&buf, "bad spec"); err != nil {
		t.Fatalf("WriteErrorJSON failed: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("error output is not valid JSON: %v", err)
	}
	if m["error"] != "bad spec" {
		t.Errorf("error field = %q, want bad spec", m["error"])
	}
}
