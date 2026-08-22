package metal

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestTemplateSpec builds a template for a registered metal.
func TestTemplateSpec(t *testing.T) {
	spec, err := TemplateSpec("Fe")
	if err != nil {
		t.Fatalf("TemplateSpec failed: %v", err)
	}
	if spec.Metal != "Fe" || spec.ICorr != 10 {
		t.Errorf("template fields wrong: %+v", spec)
	}
	if spec.MolarMass == nil || *spec.MolarMass != 55.845 {
		t.Errorf("template M missing or wrong")
	}
}

// TestTemplateSpecUnknownMetal propagates lookup errors.
func TestTemplateSpecUnknownMetal(t *testing.T) {
	if _, err := TemplateSpec("Hg"); err == nil {
		t.Error("template for unknown metal must fail")
	}
}

// TestWriteJSONRoundTrip serialises a spec and loads it back.
func TestWriteJSONRoundTrip(t *testing.T) {
	spec, err := TemplateSpec("Al")
	if err != nil {
		t.Fatalf("TemplateSpec failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, spec); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
	back, err := LoadBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("LoadBytes failed on written JSON: %v", err)
	}
	if back.Metal != "Al" {
		t.Errorf("round-trip metal = %q, want Al", back.Metal)
	}
	if *back.MolarMass != 26.9815 {
		t.Errorf("round-trip M = %v, want 26.9815", *back.MolarMass)
	}
}

// TestSpecSerialisesTargets checks that reverse targets survive a JSON
// round trip.
func TestSpecSerialisesTargets(t *testing.T) {
	cr := 0.1159
	spec := Spec{Metal: "Fe", ICorr: 10, MolarMass: f(55.845), Valence: f(2), Density: f(7.874), TargetCR: &cr}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var back Spec
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if back.TargetCR == nil || *back.TargetCR != 0.1159 {
		t.Errorf("target_cr_mm_y round trip failed: %+v", back)
	}
}

// TestHeaderCommentMentionsUnits ensures the unit legend is present.
func TestHeaderCommentMentionsUnits(t *testing.T) {
	text := HeaderComment()
	for _, want := range []string{"i_corr", "uA/cm^2", "g/mol", "g/cm^3"} {
		if !contains(text, want) {
			t.Errorf("header comment missing %q", want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
