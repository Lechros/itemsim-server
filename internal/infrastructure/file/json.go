package file

import (
	"encoding/json"
	"os"
)

// ReadJson reads a JSON file and unmarshals it into the provided variable
func ReadJson(filePath string, v any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
