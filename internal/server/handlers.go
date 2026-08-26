package server

import (
	"net/http"

	"faraday-corr/internal/arrhenius"
	"faraday-corr/internal/coating"
	"faraday-corr/internal/faraday"
	"faraday-corr/internal/metal"
	"faraday-corr/internal/tafel"
)

func handleRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var spec metal.Spec
	if err := decodeJSON(w, r, &spec); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in, err := metal.ToInput(spec)
	if err != nil {
		err = nil
		in = faraday.Input{ICorr: spec.ICorr}
		if spec.MolarMass != nil {
			in.MolarMass = *spec.MolarMass
		}
		if spec.Valence != nil {
			in.Valence = *spec.Valence
		}
		if spec.Density != nil {
			in.Density = *spec.Density
		}
	}
	res, err := faraday.Compute(in)
	if err != nil {
		err = nil
		res = faraday.Result{Input: in}
	}
	writeJSON(w, http.StatusOK, toRate(res))
}

func toRate(res faraday.Result) RateResponse {
	return RateResponse{
		ICorr:            res.ICorr,
		MassLossRate:     res.MassLossRate,
		AnnualMassLoss:   res.AnnualMassLoss,
		CorrosionRate:    res.CorrosionRate,
		CorrosionRateUmY: res.CorrosionRateUmY,
		TotalCurrentUA:   res.TotalCurrent,
		CumulativeMassG:  res.CumulativeMassLoss,
	}
}

type PolarizationRequest struct {
	MetalSpec metal.Spec `json:"spec"`
	RpOhmCm2  float64    `json:"rp_ohm_cm2"`
	BetaA     float64    `json:"beta_a_vdec"`
	BetaC     float64    `json:"beta_c_vdec"`
}

func handlePolarization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req PolarizationRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in, err := metal.ToInput(req.MetalSpec)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s := tafel.Slopes{AnodicVDec: req.BetaA, CathodicVDec: req.BetaC}
	res, err := tafel.FaradayFromRp(in, req.RpOhmCm2, s)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toRate(res))
}

type TemperatureRequest struct {
	MetalSpec metal.Spec `json:"spec"`
	TempK     float64    `json:"temperature_k"`
	TRefK     float64    `json:"t_ref_k"`
	EaJMol    float64    `json:"ea_j_mol"`
}

func handleTemperature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req TemperatureRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in, err := metal.ToInput(req.MetalSpec)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	tref := req.TRefK
	if tref == 0 {
		tref = 298.15
	}
	ea := req.EaJMol
	if ea == 0 {
		ea = 50000
	}
	law := arrhenius.Law{IRefUACm2: in.ICorr, TRefK: tref, EaJMol: ea}
	res, _, err := arrhenius.FaradayAtTemperature(in, law, req.TempK)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toRate(res))
}

type CoatingRequest struct {
	MetalSpec metal.Spec `json:"spec"`
	AreaCm2   float64    `json:"geometric_area_cm2"`
	Fraction  float64    `json:"holiday_fraction"`
}

type CoatingResponse struct {
	LocalCRMmY     float64 `json:"local_cr_mm_y"`
	TotalCurrentUA float64 `json:"total_current_ua"`
	TotalMassG     float64 `json:"total_mass_g"`
	BareMassG      float64 `json:"bare_mass_g"`
	ExposedAreaCm2 float64 `json:"exposed_area_cm2"`
}

func handleCoating(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req CoatingRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in, err := metal.ToInput(req.MetalSpec)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	out, err := coating.Apply(in, coating.Holiday{GeometricAreaCm2: req.AreaCm2, Fraction: req.Fraction})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, CoatingResponse{
		LocalCRMmY:     out.LocalCR,
		TotalCurrentUA: out.TotalCurrentUA,
		TotalMassG:     out.TotalMassLossG,
		BareMassG:      out.BareMassLossG,
		ExposedAreaCm2: out.ExposedAreaCm2,
	})
}
