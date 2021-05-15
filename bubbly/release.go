package bubbly

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/valocode/bubbly/env"
	"github.com/valocode/bubbly/parser"
)

// createReleaseSpec takes a filename pointing to a file or directory and creates
// a release spec based on the definition found in the file/directory.
// If no definition exists, an error is returned.
func createReleaseSpec(bCtx *env.BubblyContext, filename string) (*ReleaseSpec, error) {
	// Get the release by filename
	var fileParser BubblyFileParser
	err := parser.ParseFilename(bCtx, filename, &fileParser)
	if err != nil {
		return nil, fmt.Errorf("error parsing bubbly configs: %w", err)
	}

	if fileParser.Release == nil {
		return nil, fmt.Errorf("no release definition found")
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	// The base directory is needed to resolve relative file paths.
	// Set it up here and add it to the release
	var baseDir = filename
	if !stat.IsDir() {
		baseDir = filepath.Dir(baseDir)
	}

	release := fileParser.Release
	release.BaseDir = baseDir

	return release, nil
}
