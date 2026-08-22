package metal

import (
	"encoding/json"
	"io"
	"os"
)

// WriteJSON serialises a spec as indented JSON. JSON has no comment
// syntax, so the unit legend is kept out of the file; the header comment
// is printed separately by the CLI when generating a template.
func WriteJSON(w io.Writer, spec Spec) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(spec)
}

// WriteJSONFile writes a spec to a path on disk.
func WriteJSONFile(path string, spec Spec) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteJSON(f, spec)
}

// ToBytes renders the spec as indented JSON, used by tests.
func ToBytes(spec Spec) ([]byte, error) {
	return json.MarshalIndent(spec, "", "  ")
}
