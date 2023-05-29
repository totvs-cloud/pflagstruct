package code

import (
	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func (g *Generator) structFlags(st *projscan.Struct) (*orderedmap.OrderedMap[string, []*projscan.Field], error) {
	flds, err := g.fields.FindFieldsByStruct(st)
	if err != nil {
		return nil, err
	}

	refs := orderedmap.New[string, []*projscan.Field]()

	for _, fld := range flds {
		switch KindOf(fld) {
		case FieldKindNative, FieldKindTCloudTag:
			fields, _ := refs.Get("")
			refs.Set("", append(fields, fld))
		case FieldKindStruct:
			extracted, err := g.fieldFlags(fld, changecase.Param(fld.Name))
			if err != nil {
				return nil, err
			}

			merged := orderedmap.New[string, []*projscan.Field]()
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

func (g *Generator) fieldFlags(field *projscan.Field, prefix string) (*orderedmap.OrderedMap[string, []*projscan.Field], error) {
	refs := orderedmap.New[string, []*projscan.Field]()

	flds, err := g.fields.FindFieldsByStruct(field.StructRef)
	if err != nil { // TODO: warn here
		return refs, nil
	}

	for _, fld := range flds {
		switch KindOf(fld) {
		case FieldKindNative, FieldKindTCloudTag:
			fields, _ := refs.Get(prefix)
			refs.Set(prefix, append(fields, fld))
		case FieldKindStruct:
			extracted, err := g.fieldFlags(fld, changecase.Param(fld.Name))
			if err != nil {
				return nil, err
			}

			merged := orderedmap.New[string, []*projscan.Field]()
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
