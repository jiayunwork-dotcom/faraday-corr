package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const maxBodyBytes = 1 << 20

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

type RateResponse struct {
	ICorr            float64 `json:"i_corr"`
	MassLossRate     float64 `json:"mass_loss_rate"`
	AnnualMassLoss   float64 `json:"annual_mass_loss"`
	CorrosionRate    float64 `json:"corrosion_rate_mm_y"`
	CorrosionRateUmY float64 `json:"corrosion_rate_um_y"`
	TotalCurrentUA   float64 `json:"total_current_ua"`
	CumulativeMassG  float64 `json:"cumulative_mass_g"`
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("request body is not valid JSON: %v", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func readExample(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func fileServerHandler(staticDir string) http.Handler {
	inner := http.FileServer(http.Dir(staticDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		inner.ServeHTTP(w, r)
	})
}
