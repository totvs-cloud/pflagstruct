package code

import (
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

func (g *Generator) Generate(directory string, structName string, destination string) error {
	pkg, err := g.packages.FindPackageByDirectory(destination)
	if err != nil {
		return err
	}

	st, err := g.structs.FindStructByDirectoryAndName(directory, structName)
	if err != nil {
		return err
	}

	flags, err := g.structFlags(st)
	if err != nil {
		return err
	}

	blocks := []Block{
		&CommandFlagsStruct{},
		&ConstructorForFlags{},
		&ConstructorForPersistentFlags{},
		&SetterMethod{Struct: st, Flags: flags},
	}

	refs, err := g.structReferences(st)
	if err != nil {
		return err
	}

	fields, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return err
	}

	getterMethods := []*GetterMethod{{
		Prefix:  "",
		Struct:  st,
		Pointer: true,
		Fields:  fields,
	}}

	for prefix, field := range refs {
		if !field.StructRef.FromStandardLibrary() && !field.Array {
			subFields, err := g.fields.FindFieldsByStruct(field.StructRef)
			if err != nil {
				return err
			}

			getterMethods = append(getterMethods, &GetterMethod{
				Prefix:  prefix,
				Struct:  field.StructRef,
				Pointer: field.Pointer,
				Fields:  subFields,
			})
		}
	}

	for _, getterMethod := range getterMethods {
		blocks = append(blocks, getterMethod)
	}

	source := &FlagSource{
		Package: pkg,
		Blocks:  blocks,
	}

	imports, err := g.structImports(st)
	if err != nil {
		return err
	}

	for _, imp := range imports {
		source.ImportName(imp.Path, imp.Name)
	}

	source.Print()

	return nil
}
