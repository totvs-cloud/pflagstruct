package code

import (
	"path"

	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func (g *Generator) structImports(st *projscan.Struct) ([]*projscan.Package, error) {
	refs, err := g.structReferences(st)
	if err != nil {
		return nil, err
	}

	pkgsmap := map[string]*projscan.Package{st.Package.Path: st.Package}

	for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
		ref := pair.Value
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

func (g *Generator) structReferences(st *projscan.Struct) (*orderedmap.OrderedMap[string, *projscan.Field], error) {
	flds, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return nil, err
	}

	refs := orderedmap.New[string, *projscan.Field]()

	for _, fld := range flds {
		switch KindOf(fld) {
		case FieldKindStruct, FieldKindTCloudTag:
			extracted, err := g.fieldReferences(fld, changecase.Param(fld.Name))
			if err != nil {
				return nil, err
			}

			merged := orderedmap.New[string, *projscan.Field]()
			// Copy refs to merged map
			for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
				merged.Set(pair.Key, pair.Value)
			}
			// Merge extracted into merged map
			for pair := extracted.Oldest(); pair != nil; pair = pair.Next() {
				merged.Set(pair.Key, pair.Value)
			}

			refs = merged
		case FieldKindStringMap:
			refs.Set(changecase.Param(fld.Name), fld)
		}
	}

	return refs, nil
}

func (g *Generator) fieldReferences(st *projscan.Field, prefix string) (*orderedmap.OrderedMap[string, *projscan.Field], error) {
	refs := orderedmap.New[string, *projscan.Field]()
	refs.Set(prefix, st)

	flds, err := g.fields.FindFieldsByStruct(st.StructRef)
	if err != nil { // TODO: warn here
		return refs, nil
	}

	for _, fld := range flds {
		if fld.StructRef != nil {
			p := changecase.Param(path.Join(prefix, fld.Name))
			refs.Set(p, fld)

			extracted, err := g.fieldReferences(fld, p)
			if err != nil {
				// TODO: warn here
				continue
			}

			merged := orderedmap.New[string, *projscan.Field]()
			// Copy refs to merged map
			for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
				merged.Set(pair.Key, pair.Value)
			}
			// Merge extracted into merged map
			for pair := extracted.Oldest(); pair != nil; pair = pair.Next() {
				merged.Set(pair.Key, pair.Value)
			}

			refs = merged
		}
	}

	return refs, nil
}
