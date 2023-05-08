package pkg

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/totvs-cloud/pflagstruct/internal/dir"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

// Finder is a struct that provides methods for finding Go packages within a project.
type Finder struct {
	scanner  *syntree.Scanner
	projects projscan.ProjectFinder
}

// NewFinder returns a pointer to a new Finder struct with a scanner and a projects provided as arguments.
func NewFinder(scanner *syntree.Scanner, projects projscan.ProjectFinder) *Finder {
	return &Finder{scanner: scanner, projects: projects}
}

// FindPackageByDirectory returns the Go package found in the specified directory.
func (s *Finder) FindPackageByDirectory(directory string) (*projscan.Package, error) {
	directory, err := dir.AbsolutePath(directory)
	if err != nil {
		return nil, err
	}

	proj, err := s.projects.FindProjectByDirectory(directory)
	if err != nil {
		return nil, err
	}

	files, err := s.scanner.ScanDirectory(directory)
	if err != nil {
		return nil, err
	}

	result := make([]*projscan.Package, 0)

	for filename, file := range files {
		pkgname := file.Name.String()
		if pkgname == "main" {
			continue
		}

		if givenPackageList(result).containsName(pkgname) {
			continue
		}

		d := path.Dir(filename)
		p := replacePrefix(d, proj.Directory, proj.ModuleName)
		result = append(result, &projscan.Package{
			Directory: d,
			Path:      p,
			Name:      pkgname,
		})
	}

	if len(result) > 1 {
		return nil, errors.Errorf("%d package names were found in the same path %q", len(result), directory)
	}

	if len(result) == 0 {
		return nil, errors.Errorf("no Go packages were found at the path %q", directory)
	}

	return result[0], nil
}

// FindPackageByPathAndProject returns the Go package found at the specified path within the specified project.
func (s *Finder) FindPackageByPathAndProject(pkgPath string, proj *projscan.Project) (*projscan.Package, error) {
	directory, err := givenProjectAndPackagePath(proj, pkgPath).getPackageDirectory()
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(directory); err != nil {
		return nil, errors.WithStack(err)
	}

	return s.FindPackageByDirectory(directory)
}

// replacePrefix replaces the specified prefix
func replacePrefix(original string, prefix string, newPrefix string) string {
	if strings.HasPrefix(original, prefix) {
		return newPrefix + original[len(prefix):]
	}

	return original
}
