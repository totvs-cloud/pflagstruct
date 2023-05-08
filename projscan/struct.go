package projscan

import (
	"go/ast"
	"strings"
)

// Struct represents a Go struct.
type Struct struct {
	Package *Package // Package that contains the struct
	Name    string   // Name of the struct
	AST     *AST     // AST syntax tree references
}

// AST store the syntax tree references of a given Go struct.
type AST struct {
	StructType *ast.StructType // StructType syntax tree representing the struct
	File       *ast.File       // File syntax tree containing the struct
}

// FromStandardLibrary returns true if the struct is defined in a Go standard library package.
func (s *Struct) FromStandardLibrary() bool {
	return strings.HasPrefix(s.Package.Path, "std/")
}

// StructFinder provides a way to find a struct by its directory and name.
type StructFinder interface {
	FindStructByDirectoryAndName(directory, structName string) (*Struct, error)
}
