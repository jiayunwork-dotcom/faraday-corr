package metal

import (
	"encoding/json"
	"io"
	"os"
)

func WriteJSON(w io.Writer, spec Spec) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(spec)
}

func WriteJSONFile(path string, spec Spec) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteJSON(f, spec)
}

func ToBytes(spec Spec) ([]byte, error) {
	return json.MarshalIndent(spec, "", "  ")
}
