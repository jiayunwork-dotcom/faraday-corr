package metal

import (
	"fmt"
	"strings"
)

// SpecIssue is one problem found in a specification file by ValidateSpec.
type SpecIssue struct {
	Field   string
	Message string
}

func (s SpecIssue) Error() string {
	if s.Field == "" {
		return s.Message
	}
	return fmt.Sprintf("%s: %s", s.Field, s.Message)
}

// ValidateSpec performs the specification-level checks that apply before
// the numeric validation in faraday.Validate: a metal label that is given
// must resolve, a missing i_corr is reported as such, and a metal that is
// not named must supply the material fields explicitly. The returned
// slice is empty when the spec is acceptable.
func ValidateSpec(spec Spec) []SpecIssue {
	var issues []SpecIssue
	if spec.Metal != "" {
		if _, err := Lookup(spec.Metal); err != nil {
			issues = append(issues, SpecIssue{Field: "metal", Message: err.Error()})
		}
	}
	if spec.ICorr == 0 {
		issues = append(issues, SpecIssue{Field: "i_corr", Message: "must be set to a positive number of uA/cm^2"})
	}
	if spec.Metal == "" {
		if spec.MolarMass == nil {
			issues = append(issues, SpecIssue{Field: "M", Message: "must be set or supplied by the metal label"})
		}
		if spec.Valence == nil {
			issues = append(issues, SpecIssue{Field: "n", Message: "must be set or supplied by the metal label"})
		}
		if spec.Density == nil {
			issues = append(issues, SpecIssue{Field: "rho", Message: "must be set or supplied by the metal label"})
		}
	}
	return issues
}

// SummarizeIssues joins a list of issues into one error line suitable for
// the validate subcommand and for tests that assert on rejection.
func SummarizeIssues(issues []SpecIssue) error {
	if len(issues) == 0 {
		return nil
	}
	msgs := make([]string, 0, len(issues))
	for _, it := range issues {
		msgs = append(msgs, it.Error())
	}
	return fmt.Errorf("invalid specification: %s", strings.Join(msgs, "; "))
}
