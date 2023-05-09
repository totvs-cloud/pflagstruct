package code

import (
	"path"

	changecase "github.com/ku/go-change-case"
	"github.com/samber/lo"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

func (g *Generator) structFlags(st *projscan.Struct) (map[string][]*projscan.Field, error) {
	flds, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return nil, err
	}

	refs := make(map[string][]*projscan.Field)

	for _, fld := range flds {
		if fld.StructRef != nil && fld.Array {
			continue // ignore array of structs
		} else if fld.StructRef != nil {
			extracted, err := g.fieldFlags(fld, changecase.Param(fld.Name))
			if err != nil {
				return nil, err
			}

			refs = lo.Assign(refs, extracted)
		} else {
			fields := refs[""]
			refs[""] = append(fields, fld)
		}
	}

	return refs, nil
}

func (g *Generator) fieldFlags(field *projscan.Field, prefix string) (map[string][]*projscan.Field, error) {
	refs := make(map[string][]*projscan.Field)

	flds, err := g.fields.FindFieldsByStruct(field.StructRef)
	if err != nil { // TODO: warn here
		return refs, nil
	}

	for _, fld := range flds {
		if fld.StructRef != nil && fld.Array {
			continue // ignore array of structs
		} else if fld.StructRef != nil {
			subRefs, err := g.fieldFlags(fld, changecase.Param(path.Join(prefix, fld.Name)))
			if err != nil {
				// TODO: warn here
				continue
			}

			refs = lo.Assign(refs, subRefs)
		} else {
			fields := refs[prefix]
			refs[prefix] = append(fields, fld)
		}
	}

	return refs, nil
}
