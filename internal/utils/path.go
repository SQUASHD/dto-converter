// internal/utils/path.go

package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// NormalizePath converts a given path to an absolute path.
func NormalizePath(path string) (string, error) {
	// handle ~ in path
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// handle relative paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
