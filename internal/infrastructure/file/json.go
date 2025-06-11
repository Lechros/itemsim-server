package file

import (
	"encoding/json"
	"log"
	"os"
)

func ReadJson(filePath string, v any) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		log.Fatal(err)
	}
}
