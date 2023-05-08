package dir

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// AbsolutePath returns the absolute path of a given file or directory path and checks if it's a directory.
// If the given path is a file, the function returns the directory containing the file.
func AbsolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !info.IsDir() {
		path = filepath.Dir(path)
	}

	return path, nil
}
