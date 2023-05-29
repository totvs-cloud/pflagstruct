package code

import (
	"fmt"
	"path"

	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/internal/dir"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type Generator struct {
	fields   projscan.FieldFinder
	packages projscan.PackageFinder
	projects projscan.ProjectFinder
	structs  projscan.StructFinder
}

func NewGenerator(fields projscan.FieldFinder, packages projscan.PackageFinder, projects projscan.ProjectFinder, structs projscan.StructFinder) *Generator {
	return &Generator{fields: fields, packages: packages, projects: projects, structs: structs}
}

func (g *Generator) Generate(directory string, structName string, destination string) (string, error) {
	pkg, err := g.packages.FindPackageByDirectory(destination)
	if err != nil {
		return "", err
	}

	st, err := g.structs.FindStructByDirectoryAndName(directory, structName)
	if err != nil {
		return "", err
	}

	flags, err := g.structFlags(st)
	if err != nil {
		return "", err
	}

	fbn := changecase.Camel(path.Join(st.Name, "flags", "builder"))
	blocks := []Block{
		&SetUpConstructor{FlagsBuilderName: fbn, Struct: st},
		&GetConstructor{FlagsBuilderName: fbn, Struct: st},
		&FlagsBuilderStruct{Name: fbn},
		&SetterMethod{FlagsBuilderName: fbn, Struct: st, Flags: flags},
	}

	refs, err := g.structReferences(st)
	if err != nil {
		return "", err
	}

	fields, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return "", err
	}

	getterMethods := []MethodBlock{&GetterMethod{
		FlagsBuilderName: fbn,
		Prefix:           "",
		Struct:           st,
		Pointer:          true,
		Fields:           fields,
	}}

	for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
		prefix, field := pair.Key, pair.Value
		switch KindOf(field) {
		case FieldKindTCloudTag:
			getterMethods = append(getterMethods, &TagsGetterMethod{
				FlagsBuilderName: fbn,
				Prefix:           prefix,
				Struct:           field.StructRef,
				Pointer:          field.Pointer,
			})
		case FieldKindStringMap:
			getterMethods = append(getterMethods, &MapGetterMethod{
				FlagsBuilderName: fbn,
				Prefix:           prefix,
				Pointer:          field.Pointer,
			})
		case FieldKindStruct:
			subFields, err := g.fields.FindFieldsByStruct(field.StructRef)
			if err != nil {
				return "", err
			}

			getterMethods = append(getterMethods, &GetterMethod{
				FlagsBuilderName: fbn,
				Prefix:           prefix,
				Struct:           field.StructRef,
				Pointer:          field.Pointer,
				Fields:           subFields,
			})
		default:
			fmt.Println("")
		}
	}

	for _, getterMethod := range getterMethods {
		blocks = append(blocks, getterMethod)
	}

	source := &FlagSource{
		Struct:  st,
		Package: pkg,
		Blocks:  blocks,
	}

	imports, err := g.structImports(st)
	if err != nil {
		return "", err
	}

	for _, imp := range imports {
		source.ImportName(imp.Path, imp.Name)
	}

	absolutePath, err := dir.AbsolutePath(destination)
	if err != nil {
		return "", err
	}

	filepath := path.Join(absolutePath, changecase.Snake(path.Join(st.Name, "flags"))+".go")
	if err = source.WriteFile(filepath); err != nil {
		return "", err
	}

	return filepath, nil
}
