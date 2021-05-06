package interval

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// write a .zip file from bytes
// the contents of the .zip file will be unloaded to the root of a
// temporary directory on the filesystem, with both paths within the .zip
// and full path to the .zip preserved.
// That is, the path to the root of the .zip file's contents is:
// `os.TempDir()+"worker.<random-string>/filepath.Dir(<fPath>)`
func unzipFromBytes(fPath string, b []byte) (string, error) {

	tmpDir, err := os.MkdirTemp(os.TempDir(), defaultWorkerDirPattern)

	if err != nil {
		return "", err
	}

	fullTmpDir := filepath.Join(tmpDir, fPath)

	r := bytes.NewReader(b)

	zipReader, err := zip.NewReader(r, r.Size())

	if err != nil {
		return "", err
	}

	if err := unzipFromReader(zipReader, filepath.Dir(fullTmpDir)); err != nil {
		return "", err
	}

	return tmpDir, nil
}

func unzipFromReader(reader *zip.Reader, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(path, filepath.Clean(target)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

		err = targetFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
