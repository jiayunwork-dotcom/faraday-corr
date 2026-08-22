package faraday

import "fmt"

// validationError describes a rejected input value. The message names the
// field, the offending value and the required constraint so that a user
// can fix the input file without guessing.
type validationError struct {
	field  string
	value  float64
	reason string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s = %v: %s", e.field, e.value, e.reason)
}

func errNonPositive(field string, v float64) error {
	return &validationError{field: field, value: v, reason: "must be strictly positive"}
}

func errNegative(field string, v float64) error {
	return &validationError{field: field, value: v, reason: "must not be negative"}
}
