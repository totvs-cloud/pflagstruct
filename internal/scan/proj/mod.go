package proj

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/totvs-cloud/pflagstruct/internal/dir"
	"github.com/totvs-cloud/pflagstruct/projscan"
	"golang.org/x/mod/modfile"
)

type Module struct {
	directory string
	file      *modfile.File
}

// newModule returns a new instance of the Module struct with the go.mod file located in the given directory.
func newModule(directory string) (*Module, error) {
	directory, err := dir.AbsolutePath(directory)
	if err != nil {
		return nil, err
	}

	gmfp, err := findGoModFilePath(directory)
	if err != nil {
		return nil, err
	}

	gmf, err := readGoModFile(gmfp)
	if err != nil {
		return nil, err
	}

	return &Module{
		directory: path.Dir(gmfp),
		file:      gmf,
	}, nil
}

// Name returns the name of the module, or an error if it cannot be determined.
func (m *Module) Name() (string, error) {
	mod := m.file.Module
	if mod == nil {
		return "", errors.New("unable to retrieve module name")
	}

	return mod.Mod.Path, nil
}

// Directory returns the directory containing the go.mod file for the module.
func (m *Module) Directory() string {
	return m.directory
}

// Dependencies returns a list of dependencies for the module, or an error if they cannot be determined.
func (m *Module) Dependencies() ([]*projscan.Dependency, error) {
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return nil, errors.New(`the environment variable GOPATH has not been configured`)
	}

	dependencies := make([]*projscan.Dependency, 0)
	for _, req := range m.file.Require {
		dependencies = append(dependencies, &projscan.Dependency{
			Directory: path.Join(gopath, "pkg/mod", req.Mod.Path+"@"+req.Mod.Version),
			Path:      req.Mod.Path,
			Version:   req.Mod.Version,
		})
	}

	return dependencies, nil
}
