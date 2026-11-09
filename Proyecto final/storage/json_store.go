package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dataDir = "data"

func SaveData(filename string, data interface{}) error {
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(dataDir, filename)
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, file, 0644)
}

func LoadData[T any](filename string, target *[]T) error {
	filePath := filepath.Join(dataDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		*target = []T{}
		return nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, target)
}