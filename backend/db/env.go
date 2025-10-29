package db

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// loadEnvUpwards searches for a .env file in the cwd and parent directories and loads it.
// It returns an error if the file is not found or fails to load.
func loadEnvUpwards(filename string, maxDepth int) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	p := cwd
	for i := 0; i <= maxDepth; i++ {
		candidate := filepath.Join(p, filename)
		if _, statErr := os.Stat(candidate); statErr == nil {
			// found
			return godotenv.Load(candidate)
		}
		p = filepath.Dir(p)
	}
	return os.ErrNotExist
}
