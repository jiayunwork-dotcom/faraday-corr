package metal

import (
	"strings"
	"testing"

	"faraday-corr/internal/faraday"
)

// f returns a pointer to a float for building Spec values.
func f(v float64) *float64 { return &v }

// TestIronPreset pins the built-in iron entry.
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

// TestAluminumPreset pins the built-in aluminium entry.
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

// TestLookupBySymbolAndName checks case-insensitive symbol and name lookup.
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

// TestLookupUnknownMetal verifies that an unknown label returns an error.
func TestLookupUnknownMetal(t *testing.T) {
	if _, err := Lookup("Cu"); err == nil {
		t.Error("Lookup(Cu) must fail for an unregistered metal")
	}
}

// TestLoadFileJSON loads the checked-in iron example and verifies the
// parsed fields match the specification.
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

// TestLoadFileRejectsMalformedJSON verifies that a broken file is
// reported as an error rather than silently defaulted.
func TestLoadFileRejectsMalformedJSON(t *testing.T) {
	_, err := LoadBytes([]byte(`{"i_corr": `))
	if err == nil {
		t.Error("malformed JSON must fail to load")
	}
}

// TestFillFromRegistry fills omitted material fields from the metal label.
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

// TestFillFromRegistryKeepsExplicitValues verifies explicit fields win
// over the registry.
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

// TestFillFromRegistryUnknownMetal propagates the lookup error.
func TestFillFromRegistryUnknownMetal(t *testing.T) {
	spec := Spec{Metal: "Zz", ICorr: 5}
	if _, err := spec.FillFromRegistry(); err == nil {
		t.Error("unknown metal label must fail")
	}
}

// TestToInputAssemblesFaradayInput checks the conversion to faraday.Input.
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

// TestToInputRejectsNegativeCurrent pins the i<0 rule at the input layer.
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

// TestToInputAcceptsZeroCurrent pins i=0 => CR=0 at the input layer.
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

// TestToInputRejectsZeroValence pins the n<=0 rule at the input layer.
// An explicit zero in the file must not be silently replaced by the
// registry value.
func TestToInputRejectsZeroValence(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 10, MolarMass: f(55.845), Valence: f(0), Density: f(7.874)}
	if _, err := ToInput(spec); err == nil {
		t.Error("zero valence must be rejected")
	}
}

// TestToInputRejectsZeroDensity pins the rho<=0 rule at the input layer.
func TestToInputRejectsZeroDensity(t *testing.T) {
	spec := Spec{Metal: "Fe", ICorr: 10, MolarMass: f(55.845), Valence: f(2), Density: f(0)}
	if _, err := ToInput(spec); err == nil {
		t.Error("zero density must be rejected")
	}
}

// TestCompareIronAluminumDepthRates verifies the two reference metals
// produce different depth rates at one shared current density.
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

// TestValidateSpecReportsMissingICorr ensures the spec-level validator
// flags a spec without a current density.
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

// TestValidateSpecReportsMissingMaterialFields ensures a spec without a
// metal label and without material fields is flagged.
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
