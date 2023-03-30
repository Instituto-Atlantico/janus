package temp_files

import (
	"errors"
	"os"
	"path/filepath"
)

func GenerateTempFiles(dockercompose string) error {
	path := filepath.Join("/", "tmp", "janus")

	// Generate janus path on tmp
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) { // Check if path already exists
		err := os.MkdirAll(path, os.ModePerm) // Create path
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Store docker-compose file in tmp/janus.
	// Dockercompose variable is a string with the full body of our docker-compose.yml
	filepath := filepath.Join("/", "tmp", "janus", "docker-compose.yml")
	err = os.WriteFile(filepath, []byte(dockercompose), 0644)
	if err != nil {
		return err
	}
	return nil
}
