package temp_files

import (
	"errors"
	"os"
)

func GenerateTempFiles(dockercompose string) error {
	path := "/tmp/janus/"

	// Generate janus path on tmp
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) { // Check if path already exists
		err := os.Mkdir(path, os.ModePerm) // Create path
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Store docker-compose file in tmp/janus.
	// Dockercompose variable is a string with the full body of our docker-compose.yml
	err = os.WriteFile("/tmp/janus/docker-compose.yml", []byte(dockercompose), 0644)
	if err != nil {
		return err
	}
	return nil
}
