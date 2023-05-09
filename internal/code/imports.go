package code

import (
	"path"

	changecase "github.com/ku/go-change-case"
	"github.com/samber/lo"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

func (g *Generator) structImports(st *projscan.Struct) ([]*projscan.Package, error) {
	refs, err := g.structReferences(st)
	if err != nil {
		return nil, err
	}

	pkgsmap := make(map[string]*projscan.Package)

	for _, ref := range refs {
		if !ref.FromStandardLibrary() {
			pkgsmap[ref.StructRef.Package.Path] = ref.StructRef.Package
		}
	}

	pkgs := make([]*projscan.Package, 0)
	for _, pkg := range pkgsmap {
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func (g *Generator) structReferences(st *projscan.Struct) (map[string]*projscan.Field, error) {
	flds, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return nil, err
	}

	refs := make(map[string]*projscan.Field)

	for _, fld := range flds {
		if fld.StructRef != nil {
			extracted, err := g.fieldReferences(fld, changecase.Param(fld.Name))
			if err != nil {
				return nil, err
			}

			refs = lo.Assign(refs, extracted)
		}
	}

	return refs, nil
}

func (g *Generator) fieldReferences(st *projscan.Field, prefix string) (map[string]*projscan.Field, error) {
	refs := map[string]*projscan.Field{prefix: st}

	flds, err := g.fields.FindFieldsByStruct(st.StructRef)
	if err != nil { // TODO: warn here
		return refs, nil
	}

	for _, fld := range flds {
		if fld.StructRef != nil {
			p := changecase.Param(path.Join(prefix, fld.Name))
			refs[p] = fld

			extractedRefs, err := g.fieldReferences(fld, p)
			if err != nil {
				// TODO: warn here
				continue
			}

			refs = lo.Assign[string, *projscan.Field](
				refs,
				extractedRefs,
			)
		}
	}

	return refs, nil
}
