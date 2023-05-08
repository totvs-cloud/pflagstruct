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

func (m *Module) Name() (string, error) {
	mod := m.file.Module
	if mod == nil {
		return "", errors.New("unable to retrieve module name")
	}

	return mod.Mod.Path, nil
}

func (m *Module) Directory() string {
	return m.directory
}

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
