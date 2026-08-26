package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"faraday-corr/internal/faraday"
)

const formatVersion = 1

type Record struct {
	Version          int     `json:"version"`
	Name             string  `json:"name"`
	CorrosionRateUmY float64 `json:"corrosion_rate_um_y"`
	AnnualMassLoss   float64 `json:"annual_mass_loss_g_m2"`
	CumulativeLoss   float64 `json:"cumulative_mass_loss_g"`
	TotalCurrentA    float64 `json:"total_current_a"`
}

type File struct {
	Version int      `json:"version"`
	Records []Record `json:"records"`
}

func ComputeRecord(name string, in faraday.Input) (Record, error) {
	res, err := faraday.Compute(in)
	if err != nil {
		return Record{}, err
	}
	return Record{
		Version:          formatVersion,
		Name:             name,
		CorrosionRateUmY: res.CorrosionRateUmY,
		AnnualMassLoss:   res.AnnualMassLoss,
		CumulativeLoss:   res.CumulativeMassLoss,
		TotalCurrentA:    res.TotalCurrentA,
	}, nil
}

func WriteFile(path string, entries map[string]faraday.Input) error {
	if path == "" {
		return fmt.Errorf("store: empty path")
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("store: mkdir: %w", err)
		}
	}
	f := File{Version: formatVersion, Records: make([]Record, 0, len(entries))}
	for name, in := range entries {
		rec, err := ComputeRecord(name, in)
		if err != nil {
			return fmt.Errorf("store: %s: %w", name, err)
		}
		f.Records = append(f.Records, rec)
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("store: write: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("store: rename: %w", err)
	}
	return nil
}

func ReadFile(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("store: read: %w", err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, fmt.Errorf("store: decode: %w", err)
	}
	if f.Version != formatVersion {
		return File{}, fmt.Errorf("store: unsupported version %d", f.Version)
	}
	return f, nil
}

func VerifyFile(path string, entries map[string]faraday.Input, tol float64) error {
	stored, err := ReadFile(path)
	if err != nil {
		return err
	}
	if len(stored.Records) != len(entries) {
		return fmt.Errorf("store: record count mismatch")
	}
	byName := make(map[string]Record, len(stored.Records))
	for _, rec := range stored.Records {
		byName[rec.Name] = rec
	}
	for name, in := range entries {
		got, ok := byName[name]
		if !ok {
			return fmt.Errorf("store: missing %q", name)
		}
		fresh, err := ComputeRecord(name, in)
		if err != nil {
			return err
		}
		if diff(got.CorrosionRateUmY, fresh.CorrosionRateUmY, tol) {
			return fmt.Errorf("store: %s corrosion rate drift", name)
		}
	}
	return nil
}

func diff(a, b, tol float64) bool {
	if tol <= 0 {
		tol = 1e-9
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d > tol
}
