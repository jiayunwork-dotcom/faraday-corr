package metal

import (
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
)

func f(v float64) *float64 { return &v }

func TestIronPreset(t *testing.T) {
	fe := Iron()
	if fe.Symbol != "Fe" || fe.Name != "iron" {
		t.Errorf("iron identity wrong: %+v", fe)
	}
	if fe.MolarMass != 55.845 {
		t.Errorf("iron M = %v, want 55.845", fe.MolarMass)
	}
	if fe.Valence != 2 {
		t.Errorf("iron n = %v, want 2", fe.Valence)
	}
	if fe.Density != 7.874 {
		t.Errorf("iron rho = %v, want 7.874", fe.Density)
	}
}

func TestAluminumPreset(t *testing.T) {
	al := Aluminum()
	if al.Symbol != "Al" || al.Name != "aluminium" {
		t.Errorf("aluminium identity wrong: %+v", al)
	}
	if al.MolarMass != 26.9815 {
		t.Errorf("aluminium M = %v, want 26.9815", al.MolarMass)
	}
	if al.Valence != 3 {
		t.Errorf("aluminium n = %v, want 3", al.Valence)
	}
	if al.Density != 2.70 {
		t.Errorf("aluminium rho = %v, want 2.70", al.Density)
	}
}

func TestLookupBySymbolAndName(t *testing.T) {
	cases := []string{"Fe", "fe", "FE", "iron", "IRON"}
	for _, key := range cases {
		m, err := Lookup(key)
		if err != nil {
			t.Errorf("Lookup(%q) failed: %v", key, err)
			continue
		}
		if m.Symbol != "Fe" {
			t.Errorf("Lookup(%q) = %q, want Fe", key, m.Symbol)
		}
	}
}

func TestLookupUnknownMetal(t *testing.T) {
	if _, err := Lookup("Cu"); err == nil {
		t.Error("Lookup(Cu) must fail for an unregistered metal")
	}
}

func TestLoadFileJSON(t *testing.T) {
	spec, err := LoadFile("../../example/fe-seawater.json")
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}
	if spec.Metal != "Fe" {
		t.Errorf("metal = %q, want Fe", spec.Metal)
	}
	if spec.ICorr != 10.0 {
		t.Errorf("i_corr = %v, want 10", spec.ICorr)
	}
	if *spec.MolarMass != 55.845 {
		t.Errorf("M = %v, want 55.845", *spec.MolarMass)
	}
	if *spec.Valence != 2 {
		t.Errorf("n = %v, want 2", *spec.Valence)
	}
	if *spec.Density != 7.874 {
		t.Errorf("rho = %v, want 7.874", *spec.Density)
	}
}

func TestLoadFileRejectsMalformedJSON(t *testing.T) {
	_, err := LoadBytes([]byte(`{"i_corr": `))
	if err == nil {
		t.Error("malformed JSON must fail to load")
	}
}

func TestFillFromRegistry(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 5}
	filled, err := spec.FillFromRegistry()
	if err != nil {
		t.Fatalf("FillFromRegistry failed: %v", err)
	}
	if *filled.MolarMass != 55.845 || *filled.Valence != 2 || *filled.Density != 7.874 {
		t.Errorf("fill incomplete: %+v", filled)
	}
}

func TestFillFromRegistryKeepsExplicitValues(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 5, Valence: f(3)}
	filled, err := spec.FillFromRegistry()
	if err != nil {
		t.Fatalf("FillFromRegistry failed: %v", err)
	}
	if *filled.Valence != 3 {
		t.Errorf("valence = %v, want explicit 3", *filled.Valence)
	}
	if *filled.MolarMass != 55.845 {
		t.Errorf("molar mass = %v, want registry 55.845", *filled.MolarMass)
	}
}

func TestFillFromRegistryUnknownMetal(t *testing.T) {
	spec := Spec{Metal: "Zz", ICorr: 5}
	if _, err := spec.FillFromRegistry(); err == nil {
		t.Error("unknown metal label must fail")
	}
}

func TestToInputAssemblesFaradayInput(t *testing.T) {
	spec := Spec{
		Metal: "Fe", ICorr: 10,
		MolarMass: f(55.845), Valence: f(2), Density: f(7.874),
		Area: f(2), DurationY: f(1),
	}
	in, err := ToInput(spec)
	if err != nil {
		t.Fatalf("ToInput failed: %v", err)
	}
	want := faraday.Input{ICorr: 10, MolarMass: 55.845, Valence: 2, Density: 7.874, Area: 2, DurationY: 1}
	if in != want {
		t.Errorf("ToInput = %+v, want %+v", in, want)
	}
}

func TestToInputRejectsNegativeCurrent(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: -3, MolarMass: f(55.845), Valence: f(2), Density: f(7.874)}
	_, err := ToInput(spec)
	if err == nil {
		t.Error("negative i_corr must be rejected")
	}
	if !strings.Contains(err.Error(), "i_corr") {
		t.Errorf("error %q should name i_corr", err.Error())
	}
}

func TestToInputAcceptsZeroCurrent(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 0, MolarMass: f(55.845), Valence: f(2), Density: f(7.874)}
	in, err := ToInput(spec)
	if err != nil {
		t.Fatalf("zero current must be accepted, got %v", err)
	}
	res, err := faraday.Compute(in)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	if res.CorrosionRate != 0 {
		t.Errorf("CR at i=0 = %v, want 0", res.CorrosionRate)
	}
}

func TestToInputRejectsZeroValence(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 10, MolarMass: f(55.845), Valence: f(0), Density: f(7.874)}
	if _, err := ToInput(spec); err == nil {
		t.Error("zero valence must be rejected")
	}
}

func TestToInputRejectsZeroDensity(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 10, MolarMass: f(55.845), Valence: f(2), Density: f(0)}
	if _, err := ToInput(spec); err == nil {
		t.Error("zero density must be rejected")
	}
}

func TestCompareIronAluminumDepthRates(t *testing.T) {
	cmp, err := CompareRegistry(10.0)
	if err != nil {
		t.Fatalf("CompareRegistry failed: %v", err)
	}
	if cmp.FirstCR == cmp.SecondCR {
		t.Errorf("Fe CR %v must differ from Al CR %v", cmp.FirstCR, cmp.SecondCR)
	}
	if cmp.FirstCR <= 0 || cmp.SecondCR <= 0 {
		t.Errorf("depth rates must be positive: Fe %v, Al %v", cmp.FirstCR, cmp.SecondCR)
	}
}

func TestValidateSpecReportsMissingICorr(t *testing.T) {
	issues := ValidateSpec(Spec{Metal: "Fe", MolarMass: f(55.845), Valence: f(2), Density: f(7.874)})
	if len(issues) == 0 {
		t.Error("spec without i_corr must produce an issue")
	}
	found := false
	for _, it := range issues {
		if it.Field == "i_corr" {
			found = true
		}
	}
	if !found {
		t.Errorf("issues %v should include i_corr", issues)
	}
}

func TestValidateSpecReportsMissingMaterialFields(t *testing.T) {
	issues := ValidateSpec(Spec{ICorr: 10})
	fields := map[string]bool{}
	for _, it := range issues {
		fields[it.Field] = true
	}
	for _, want := range []string{"M", "n", "rho"} {
		if !fields[want] {
			t.Errorf("issues %v should include %q", issues, want)
		}
	}
}
