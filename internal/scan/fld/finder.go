package fld

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/pkg/errors"

	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

// Finder provides a way to find fields in Go struct definitions.
type Finder struct {
	packages projscan.PackageFinder
	projects projscan.ProjectFinder
	structs  projscan.StructFinder
}

// NewFinder creates a new instance of Finder.
func NewFinder(packages projscan.PackageFinder, projects projscan.ProjectFinder, structs projscan.StructFinder) *Finder {
	return &Finder{packages: packages, projects: projects, structs: structs}
}

// FindFieldsByStruct returns a slice of fields for the given struct.
func (f *Finder) FindFieldsByStruct(st *projscan.Struct) ([]*projscan.Field, error) {
	if st.AST == nil || st.AST.StructType == nil || st.AST.StructType.Fields == nil || st.AST.StructType.Fields.List == nil {
		return nil, errors.Errorf("%q struct fields were not found", st.Name)
	}

	proj, err := f.projects.FindProjectByDirectory(st.Package.Directory)
	if err != nil {
		return nil, err
	}

	result := make([]*projscan.Field, 0)

	for _, field := range st.AST.StructType.Fields.List {
		for _, name := range field.Names {
			built, err := f.buildField(field.Type, st, proj, &projscan.Field{
				Name:      name.String(),
				Type:      "",
				Doc:       extractDoc(field.Doc),
				StructRef: nil,
				Pointer:   false,
				Array:     false,
			})
			if err != nil {
				return nil, err
			}

			result = append(result, built)
		}
	}

	return result, nil
}

// buildField creates a new Field based on the given parameters.
func (f *Finder) buildField(expr ast.Expr, st *projscan.Struct, proj *projscan.Project, field *projscan.Field) (*projscan.Field, error) {
	switch x := expr.(type) {
	case *ast.StarExpr:
		// it means that the field type is a pointer
		return f.buildField(x.X, st, proj, &projscan.Field{
			Name:         field.Name,
			Type:         field.Type,
			Doc:          field.Doc,
			StructRef:    field.StructRef,
			Pointer:      true,
			Array:        field.Array,
			ArrayPointer: field.ArrayPointer,
		})
	case *ast.ArrayType:
		// it means that the field type is an array
		return f.buildField(x.Elt, st, proj, &projscan.Field{
			Name:         field.Name,
			Type:         field.Type,
			Doc:          field.Doc,
			StructRef:    field.StructRef,
			Pointer:      false,
			Array:        true,
			ArrayPointer: field.Pointer,
		})
	case *ast.Ident:
		// it means that the field type is either a built-in type or a struct from the same package
		if projscan.FieldType(x.Name).IsValid() {
			return &projscan.Field{
				Name:         field.Name,
				Type:         projscan.FieldType(x.Name),
				Doc:          field.Doc,
				StructRef:    field.StructRef,
				Pointer:      field.Pointer,
				Array:        field.Array,
				ArrayPointer: field.ArrayPointer,
			}, nil
		}

		structRef, err := f.structs.FindStructByDirectoryAndName(st.Package.Directory, x.Name)
		if err != nil {
			return nil, err
		}

		return &projscan.Field{
			Name:         field.Name,
			Type:         projscan.FieldType(x.Name),
			Doc:          field.Doc,
			StructRef:    structRef,
			Pointer:      field.Pointer,
			Array:        field.Array,
			ArrayPointer: field.ArrayPointer,
		}, nil
	case *ast.SelectorExpr:
		// it means that the field type is a struct from another package
		if ident, ok := x.X.(*ast.Ident); ok {
			path, err := syntree.WrapFile(st.AST.File).FindPackagePathByName(ident.Name)
			if err != nil {
				return nil, err
			}

			pkg, err := f.packages.FindPackageByPathAndProject(path, proj)
			if err != nil {
				return nil, err
			}

			structRef, err := f.structs.FindStructByDirectoryAndName(pkg.Directory, x.Sel.Name)
			if err != nil {
				return nil, err
			}

			return &projscan.Field{
				Name:         field.Name,
				Type:         projscan.FieldType(fmt.Sprintf("%s.%s", pkg.Name, x.Sel.Name)),
				Doc:          field.Doc,
				StructRef:    structRef,
				Pointer:      field.Pointer,
				Array:        field.Array,
				ArrayPointer: field.ArrayPointer,
			}, nil
		}
	case *ast.MapType:
		// it means that the field type is a map
		key, err := f.buildField(x.Key, st, proj, &projscan.Field{
			Name:         field.Name,
			Type:         "",
			Doc:          field.Doc,
			StructRef:    nil,
			Pointer:      false,
			Array:        false,
			ArrayPointer: false,
		})
		if err != nil {
			return nil, err
		}

		value, err := f.buildField(x.Value, st, proj, &projscan.Field{
			Name:         field.Name,
			Type:         "",
			Doc:          field.Doc,
			StructRef:    nil,
			Pointer:      false,
			Array:        false,
			ArrayPointer: false,
		})
		if err != nil {
			return nil, err
		}

		return &projscan.Field{
			Name:         field.Name,
			Type:         projscan.FieldType(fmt.Sprintf("map[%s]%s", key.Type, value.Type)),
			Doc:          field.Doc,
			StructRef:    field.StructRef,
			Pointer:      field.Pointer,
			Array:        field.Array,
			ArrayPointer: field.ArrayPointer,
		}, nil
	}

	// if the expression is of a different type, the function returns an error
	return nil, errors.New("field type not found")
}

// extractDoc returns the documentation text for the given ast.CommentGroup. If doc is nil, an empty string is returned.
func extractDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}

	comments := make([]string, 0)
	for _, d := range doc.List {
		comments = append(comments, strings.TrimSpace(strings.TrimPrefix(d.Text, "//")))
	}

	return strings.Join(comments, "\n")
}
