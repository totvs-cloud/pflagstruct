package syntree

import (
	"go/ast"
	"strings"

	"github.com/pkg/errors"
)

// FileWrapper is a wrapper around the ast.File struct that provides convenience methods.
type FileWrapper struct {
	file *ast.File
}

// WrapFile creates a new FileWrapper instance based on an ast.File struct.
func WrapFile(file *ast.File) *FileWrapper {
	return &FileWrapper{file: file}
}

// importWrapper is a wrapper around the ast.ImportSpec struct that provides convenience methods.
type importWrapper struct {
	imp *ast.ImportSpec
}

// FindPackagePathByName searches the file's import statements for a package with the specified name and returns its path.
func (w *FileWrapper) FindPackagePathByName(name string) (string, error) {
	for _, imp := range w.Imports() {
		if imp.EqualsName(name) {
			return imp.Path(), nil
		}
	}

	return "", errors.Errorf("package path %q not found", name)
}

// Imports returns a slice of importWrapper instances representing the file's import statements.
func (w *FileWrapper) Imports() []*importWrapper {
	imports := make([]*importWrapper, 0)
	for _, imp := range w.file.Imports {
		imports = append(imports, &importWrapper{imp: imp})
	}

	return imports
}

// Path returns the import path of the package specified in the import statement.
func (w *importWrapper) Path() string {
	val := w.imp.Path.Value
	val = strings.TrimPrefix(val, `"`)
	val = strings.TrimSuffix(val, `"`)

	return val
}

// Name returns the name specified in the import statement or the package name if not specified.
func (w *importWrapper) Name() string {
	if w.imp.Name != nil {
		return w.imp.Name.String()
	}

	parts := strings.Split(w.Path(), "/")

	return parts[len(parts)-1]
}

// EqualsName checks if the import statement specifies the package name.
func (w *importWrapper) EqualsName(name string) bool {
	return w.Name() == name
}
