package projscan

// Project represents a Go project that contains a module with a given name and resides in a specific directory.
type Project struct {
	ModuleName   string        // Name of the project's module
	Directory    string        // Path to the project directory
	Dependencies []*Dependency // List of the project's dependencies
}

// Dependency represents a dependency of a Go project.
type Dependency struct {
	Directory string // Path to the directory containing the dependency's source code
	Path      string // Import path of the dependency
	Version   string // Version of the dependency (if specified)
}

// ProjectFinder provides a way to find a project by its directory or by the package it belongs to.
type ProjectFinder interface {
	FindProjectByDirectory(directory string) (*Project, error)
	FindProjectByPackage(pkg *Package) (*Project, error)
}
