package metal

import (
	"fmt"
	"strings"
)

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
