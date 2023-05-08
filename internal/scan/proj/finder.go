package proj

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
	"golang.org/x/mod/modfile"
)

type Finder struct {
	scanner *syntree.Scanner
}

func NewFinder(scanner *syntree.Scanner) *Finder {
	return &Finder{scanner: scanner}
}

func (s *Finder) FindProjectByDirectory(directory string) (*projscan.Project, error) {
	mod, err := newModule(directory)
	if err != nil {
		return nil, err
	}

	modName, err := mod.Name()
	if err != nil {
		return nil, err
	}

	dependencies, err := mod.Dependencies()
	if err != nil {
		return nil, err
	}

	return &projscan.Project{
		ModuleName:   modName,
		Directory:    mod.Directory(),
		Dependencies: dependencies,
	}, nil
}

func (s *Finder) FindProjectByPackage(pkg *projscan.Package) (*projscan.Project, error) {
	directory := pkg.Directory
	return s.FindProjectByDirectory(directory)
}

func readGoModFile(path string) (*modfile.File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	file, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func findGoModFilePath(directory string) (string, error) {
	filename := "go.mod"

	for {
		files, err := os.ReadDir(directory)
		if err != nil {
			return "", errors.WithStack(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if file.Name() == filename {
				return filepath.Join(directory, file.Name()), nil
			}
		}

		// If the file is not found in the current directory, move up to the parent directory
		parent := filepath.Dir(directory)
		if parent == directory {
			// If we've reached the root directory and haven't found the file, return an error
			return "", errors.Errorf("file not found: %s", filename)
		}

		directory = parent
	}
}
