package metal

import "strconv"

func trimZeros(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "" || s == "-0" {
		s = "0"
	}
	return s
}
