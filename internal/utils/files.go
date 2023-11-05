package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FileProcessor is a callback function type that processes a file.
// It should return an error if processing fails.
type FileProcessor func(path string, info os.FileInfo) error

// WalkAndProcess traverses the directory at the given path, finds files with the specified suffix,
// and processes them using the provided FileProcessor callback function.
func WalkAndProcess(rootDir, suffix string, processor FileProcessor) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), suffix) {
			if err := processor(path, info); err != nil {
				return err
			}
		}

		return nil
	})
}
