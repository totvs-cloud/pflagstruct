package projscan

// Package represents a Go package.
type Package struct {
	Directory string // Path to the directory containing the package's source code
	Path      string // Import path of the package
	Name      string // Name of the package
}

// PackageFinder provides a way to find a Go package by its directory or by its path and project.
type PackageFinder interface {
	FindPackageByDirectory(directory string) (*Package, error)
	FindPackageByPathAndProject(pkgPath string, proj *Project) (*Package, error)
}
