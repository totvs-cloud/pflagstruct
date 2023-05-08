package pkg

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

// pkgListWrapper is a wrapper struct for a list of packages.
type pkgListWrapper struct {
	pkgs []*projscan.Package
}

// givenPackageList creates a new pkgListWrapper from a list of packages.
func givenPackageList(pkgs []*projscan.Package) *pkgListWrapper {
	return &pkgListWrapper{pkgs: pkgs}
}

// containsName checks if the given package name is in the list of packages.
func (w *pkgListWrapper) containsName(pkgname string) bool {
	for _, pkg := range w.pkgs {
		if pkg.Name == pkgname {
			return true
		}
	}

	return false
}

// projAndPkgWrapper is a wrapper struct for a project and package path.
type projAndPkgWrapper struct {
	proj    *projscan.Project
	pkgpath string
}

// givenProjectAndPackagePath creates a new projAndPkgWrapper from a project and package path.
func givenProjectAndPackagePath(proj *projscan.Project, pkgpath string) *projAndPkgWrapper {
	return &projAndPkgWrapper{proj: proj, pkgpath: pkgpath}
}

// isInternal checks if the package is part of the project.
func (w *projAndPkgWrapper) isInternal() bool {
	return strings.HasPrefix(w.pkgpath, w.proj.ModuleName)
}

// isExternal checks if the package is a dependency of the project.
// If it is, the function returns the dependency and true. Otherwise, it returns nil and false.
func (w *projAndPkgWrapper) isExternal() (*projscan.Dependency, bool) {
	for _, dependency := range w.proj.Dependencies {
		if strings.HasPrefix(w.pkgpath, dependency.Path) {
			return dependency, true
		}
	}

	return nil, false
}

// getPackageDirectory returns the directory path of the package.
func (w *projAndPkgWrapper) getPackageDirectory() (string, error) {
	if w.isInternal() {
		return replacePrefix(w.pkgpath, w.proj.ModuleName, w.proj.Directory), nil
	}

	if dep, ok := w.isExternal(); ok {
		return replacePrefix(w.pkgpath, dep.Path, dep.Directory), nil
	}

	goroot, ok := os.LookupEnv("GOROOT")
	if !ok {
		return "", errors.New(`the environment variable GOROOT has not been configured`)
	}

	return path.Join(goroot, "src", w.pkgpath), nil
}
