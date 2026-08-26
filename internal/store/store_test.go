package store

import (
	"path/filepath"
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

func TestWriteReadStore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	entries := map[string]faraday.Input{"fe": sampleInput()}
	if err := WriteFile(path, entries); err != nil {
		t.Fatal(err)
	}
	f, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Records) != 1 || f.Records[0].CorrosionRateUmY <= 0 {
		t.Fatalf("records = %+v", f.Records)
	}
}

func TestVerifyFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "verify.json")
	entries := map[string]faraday.Input{"fe": sampleInput()}
	if err := WriteFile(path, entries); err != nil {
		t.Fatal(err)
	}
	if err := VerifyFile(path, entries, 1e-12); err != nil {
		t.Fatal(err)
	}
}
