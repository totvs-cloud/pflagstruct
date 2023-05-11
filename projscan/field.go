//go:generate go-enum

package projscan

import (
	"strings"
)

// Field represents a field of a struct.
type Field struct {
	Name      string    // Name of the field
	Type      FieldType // Type of the field
	Doc       string    // Documentation for the field
	StructRef *Struct   // Reference to the struct that contains this field
	Pointer   bool      // Indicates whether the field is a pointer type or not
	Array     bool      // Indicates whether the field is an array type or not
}

// FieldType defines the available field types in Go
// ENUM(string, bool, int, int8, int16, int32, int64, float32, float64)
type FieldType string

// FromStandardLibrary returns true if the field's containing struct is part of the Go standard library.
func (s *Field) FromStandardLibrary() bool {
	if s.StructRef == nil {
		return true
	}

	return strings.HasPrefix(s.StructRef.Package.Path, "std/")
}

func (s *Field) IsTCloudTags() bool {
	return s.HasStructRef("github.com/totvs-cloud/tcloud-iaas-sdk/pkg/tags", "Tags")
}

func (s *Field) HasStructRef(path, name string) bool {
	return s.StructRef != nil &&
		s.StructRef.Name == name &&
		s.StructRef.Package != nil &&
		s.StructRef.Package.Path == path
}

// PackageDirectory returns the directory of the package that this field belongs to.
func (s *Field) PackageDirectory() string {
	if s.StructRef == nil || s.StructRef.Package == nil {
		return ""
	}

	return s.StructRef.Package.Directory
}

// FieldFinder provides a way to find all fields of a struct.
type FieldFinder interface {
	FindFieldsByStruct(st *Struct) ([]*Field, error)
}
