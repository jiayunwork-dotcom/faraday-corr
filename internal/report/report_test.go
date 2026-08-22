package report

import (
	"bytes"
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
)

// TestNumFormatting checks the compact number rendering: zero, small
// scientific, and ordinary values.
func TestNumFormatting(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{2.894e-9, "2.894e-09"},
		{0.1159, "0.1159"},
		{96485.33212, "9.649e+04"},
	}
	for _, c := range cases {
		got := Num(c.in)
		if got != c.want {
			t.Errorf("Num(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestWriteResultContainsKeyLines checks that the rate report carries the
// mass loss rate and the depth rate lines with their units.
func TestWriteResultContainsKeyLines(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1, DurationY: 1}
	res, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteResult(&buf, res, "Fe"); err != nil {
		t.Fatalf("WriteResult failed: %v", err)
	}
	text := buf.String()
	for _, want := range []string{
		"Mass loss rate:",
		"Corrosion rate:",
		"mm/y",
		"Cumulative loss:",
		"I = i_corr*A:",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("report missing %q:\n%s", want, text)
		}
	}
}

// TestWriteResultOmitsAreaLines checks that without an area the report
// does not fabricate a total current line.
func TestWriteResultOmitsAreaLines(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874}
	res, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteResult(&buf, res, "Fe"); err != nil {
		t.Fatalf("WriteResult failed: %v", err)
	}
	if strings.Contains(buf.String(), "I = i_corr*A:") {
		t.Errorf("report must not print total current without an area:\n%s", buf.String())
	}
}

// TestUsageMentionsSubcommands checks that the help text documents the
// rate subcommand and the example path.
func TestUsageMentionsSubcommands(t *testing.T) {
	var buf bytes.Buffer
	if err := Usage(&buf, "faraday-corr"); err != nil {
		t.Fatalf("Usage failed: %v", err)
	}
	text := buf.String()
	for _, want := range []string{"rate", "metals", "validate", "example/fe-seawater.json"} {
		if !strings.Contains(text, want) {
			t.Errorf("usage missing %q:\n%s", want, text)
		}
	}
}

// TestWriteMetalsListsEntries checks the registry listing shows both
// reference metals.
func TestWriteMetalsListsEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteMetals(&buf, metal.All()); err != nil {
		t.Fatalf("WriteMetals failed: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "Fe") || !strings.Contains(text, "Al") {
		t.Errorf("metals table missing entries:\n%s", text)
	}
}

// TestWriteComparisonShowsBothRates checks the iron-versus-aluminium
// comparison output.
func TestWriteComparisonShowsBothRates(t *testing.T) {
	cmp, err := metal.CompareRegistry(10.0)
	if err != nil {
		t.Fatalf("CompareRegistry failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteComparison(&buf, cmp); err != nil {
		t.Fatalf("WriteComparison failed: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "Fe") || !strings.Contains(text, "Al") {
		t.Errorf("comparison missing metals:\n%s", text)
	}
}

// TestFormatInvariantLine checks the one-line invariant rendering.
func TestFormatInvariantLine(t *testing.T) {
	got := FormatInvariantLine("i*2 => CR*2", 0.1159, 0.2318, 2)
	if !strings.Contains(got, "ratio=2") {
		t.Errorf("line = %q, want ratio=2", got)
	}
}

// TestPercentDiff checks the relative difference helper.
func TestPercentDiff(t *testing.T) {
	if got := PercentDiff(10, 12); got < 19 || got > 21 {
		t.Errorf("PercentDiff(10,12) = %v, want about 20", got)
	}
	if got := PercentDiff(0, 0); got != 0 {
		t.Errorf("PercentDiff(0,0) = %v, want 0", got)
	}
}
