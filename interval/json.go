package interval

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultWorkerDirPattern = "worker.*"
)

// write to a temporary JSON file from received bytes
// the path to the resulting file is `os.TempDir()+"worker.<random-string>/<fPath>`
func createJSONFromBytes(fPath string, b []byte) (string, error) {

	tmpDir, err := os.MkdirTemp(os.TempDir(), defaultWorkerDirPattern)

	if err != nil {
		return "", fmt.Errorf("failed to create temporary worker directory: %w", err)
	}

	tmpFile := filepath.Join(tmpDir, fPath)

	// if the file is nested, create necessary parent directories
	if err := os.MkdirAll(filepath.Dir(tmpFile), 0755); err != nil {
		return "", err
	}

	// create the file, or truncate if already existing
	dst, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}

	// write new content to the file
	_, err = dst.Write(b)

	if err != nil {
		return "", err
	}

	err = dst.Close()

	if err != nil {
		return "", err
	}

	return tmpDir, nil
}
