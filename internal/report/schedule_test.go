package report

import (
	"bytes"
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
)

// TestWriteSummary prints the engineering units for a computed result.
func TestWriteSummary(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1, DurationY: 1}
	res, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteSummary(&buf, faraday.Summarize(res)); err != nil {
		t.Fatalf("WriteSummary failed: %v", err)
	}
	text := buf.String()
	for _, want := range []string{"mpy", "g/(m^2 y)", "mdd"} {
		if !strings.Contains(text, want) {
			t.Errorf("summary missing %q:\n%s", want, text)
		}
	}
}

// TestWriteReverse prints the solved current density.
func TestWriteReverse(t *testing.T) {
	out, err := faraday.ReverseCurrentDensity(faraday.ReverseInput{
		MolarMass: 55.845, Valence: 2, Density: 7.874,
		Target: faraday.TargetCorrosionRate, CR: 0.1159,
	})
	if err != nil {
		t.Fatalf("ReverseCurrentDensity failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteReverse(&buf, out); err != nil {
		t.Fatalf("WriteReverse failed: %v", err)
	}
	if !strings.Contains(buf.String(), "Solved i_corr:") {
		t.Errorf("reverse output missing solved current:\n%s", buf.String())
	}
}

// TestWriteLife prints the life-cycle figures.
func TestWriteLife(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteLife(&buf, 0.1159, 1.0, 10.0, 5.0); err != nil {
		t.Fatalf("WriteLife failed: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "Years to consume allowance:") ||
		!strings.Contains(text, "Remaining thickness after") {
		t.Errorf("life output incomplete:\n%s", text)
	}
}

// TestWriteSchedule prints a multi-year cumulative loss table.
func TestWriteSchedule(t *testing.T) {
	in := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 1}
	s, err := faraday.BuildSchedule(in, StandardScheduleYears)
	if err != nil {
		t.Fatalf("BuildSchedule failed: %v", err)
	}
	var buf bytes.Buffer
	if err := WriteSchedule(&buf, s); err != nil {
		t.Fatalf("WriteSchedule failed: %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "years") || !strings.Contains(text, "cumulative") {
		t.Errorf("schedule output incomplete:\n%s", text)
	}
	if len(s.Years) != len(StandardScheduleYears) {
		t.Errorf("schedule year count = %d, want %d", len(s.Years), len(StandardScheduleYears))
	}
}

// TestNumFixed keeps fixed-point rendering stable.
func TestNumFixed(t *testing.T) {
	if got := NumFixed(55.845); got != "55.8450" {
		t.Errorf("NumFixed(55.845) = %q, want 55.8450", got)
	}
}
