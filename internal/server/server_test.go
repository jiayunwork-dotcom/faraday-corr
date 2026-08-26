package server_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
	"faraday-corr/internal/server"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	src := filepath.Join("..", "..", "example", "fe-seawater.json")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	dir := t.TempDir()
	ex := filepath.Join(dir, "fe-seawater.json")
	if err := os.WriteFile(ex, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	srv := httptest.NewServer(server.New(dir, ex))
	t.Cleanup(srv.Close)
	return srv
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var out server.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !out.OK || out.Service != "faraday-corr" {
		t.Errorf("health=%+v", out)
	}
}

func TestRateEndpointMatchesCompute(t *testing.T) {
	srv := newTestServer(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "example", "fe-seawater.json"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	spec, err := metal.LoadBytes(raw)
	if err != nil {
		t.Fatalf("LoadBytes: %v", err)
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		t.Fatalf("ToInput: %v", err)
	}
	want, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	resp, err := http.Post(srv.URL+"/api/rate", "application/json", strings.NewReader(string(raw)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var out server.RateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if math.Abs(out.CorrosionRate-want.CorrosionRate)/want.CorrosionRate > 1e-12 {
		t.Errorf("CR api=%v compute=%v", out.CorrosionRate, want.CorrosionRate)
	}
	if math.Abs(out.MassLossRate-want.MassLossRate)/want.MassLossRate > 1e-12 {
		t.Errorf("mdot api=%v compute=%v", out.MassLossRate, want.MassLossRate)
	}
}

func TestRateEndpointRejectsBadInput(t *testing.T) {
	srv := newTestServer(t)
	body := `{"metal":"Fe","i_corr":-1,"M":55.845,"n":2,"rho":7.874}`
	resp, err := http.Post(srv.URL+"/api/rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400", resp.StatusCode)
	}
	var er server.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if strings.TrimSpace(er.Error) == "" {
		t.Error("empty error")
	}
}

func TestExampleEndpoint(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/api/example")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var spec metal.Spec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if spec.ICorr != 10 {
		t.Errorf("i_corr=%v", spec.ICorr)
	}
}

func TestPolarizationEndpoint(t *testing.T) {
	srv := newTestServer(t)
	body := `{
		"spec": {"metal":"Fe","i_corr":10,"M":55.845,"n":2,"rho":7.874},
		"rp_ohm_cm2": 1737.177,
		"beta_a_vdec": 0.06,
		"beta_c_vdec": 0.12
	}`
	resp, err := http.Post(srv.URL+"/api/polarization", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	var out server.RateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.CorrosionRate <= 0 {
		t.Errorf("CR=%v", out.CorrosionRate)
	}
}
