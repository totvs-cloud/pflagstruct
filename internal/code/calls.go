package code

import (
	"fmt"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type SetterCall struct {
	Prefix string
	Struct *projscan.Struct
	Field  *projscan.Field
}

func (s *SetterCall) Flag() string {
	return changecase.Param(path.Join(s.Prefix, s.Field.Name))
}

func (s *SetterCall) CobraMethod() string {
	if s.Field.IsTCloudTags() {
		return "StringSlice"
	}

	var suffix string
	if s.Field.Array {
		suffix = "Slice"
	}

	return changecase.Pascal(path.Join(s.Field.Type.String(), suffix))
}

func (s *SetterCall) UsageMessage() string {
	doc := strings.TrimSpace(s.Field.Doc)

	if s.Field.IsTCloudTags() {
		msg := fmt.Sprintf("the desired key-value pairs separated by commas (%s key1=value1,key2=value2,key3=value3)", s.Flag())
		if len(doc) > 0 {
			doc += ". Provide " + msg
		} else {
			doc += "provide " + msg
		}
	}

	return doc
}

func (s *SetterCall) DefaultValue() *jen.Statement {
	if s.Field.Array {
		return jen.Nil()
	}

	switch s.Field.Type {
	case projscan.FieldTypeInt, projscan.FieldTypeInt8, projscan.FieldTypeInt16, projscan.FieldTypeInt32, projscan.FieldTypeInt64:
		return jen.Lit(0)
	case projscan.FieldTypeFloat32, projscan.FieldTypeFloat64:
		return jen.Lit(0.0)
	case projscan.FieldTypeString:
		return jen.Lit("")
	case projscan.FieldTypeBool:
		return jen.Lit(false)
	default:
		return jen.Nil()
	}
}

func (s *SetterCall) Statement() *jen.Statement {
	if s.Field.IsTCloudTags() || s.Field.Type.IsValid() {
		return jen.Id("cf").
			Dot("flags").Dot(s.CobraMethod()).
			Call(jen.Lit(s.Flag()), s.DefaultValue(), jen.Lit(s.UsageMessage()))
	}

	return nil
}

type GetterCall struct {
	Prefix  string
	Struct  *projscan.Struct
	Pointer bool
	Field   *projscan.Field
}

func (g *GetterCall) CobraMethod() string {
	var suffix string
	if g.Field.Array {
		suffix = "Slice"
	}

	return changecase.Pascal(path.Join("get", g.Field.Type.String(), suffix))
}

func (g *GetterCall) Flag() string {
	return changecase.Param(path.Join(g.Prefix, g.Field.Name))
}

func (g *GetterCall) Statement() *jen.Statement {
	id := jen.Id(changecase.Camel(g.Struct.Name))

	returnId := jen.Nil()
	if !g.Pointer {
		returnId = jen.Id(changecase.Camel(g.Struct.Name))
	}

	if g.Field.StructRef != nil && !g.Field.StructRef.FromStandardLibrary() && (!g.Field.Array || g.Field.IsTCloudTags()) {
		return jen.If(jen.List(id.Dot(g.Field.Name), jen.Err()).Op("=").
			Id("cf").Dot(changecase.Camel(path.Join("Get", g.Prefix, g.Field.Name))).Call(), jen.Err().Op("!=").Nil()).
			Block(jen.Return().List(returnId, jen.Err()))
	}

	if g.Field.Type.IsValid() {
		return jen.If(jen.List(id.Dot(g.Field.Name), jen.Err()).Op("=").
			Id("cf").Dot("flags").Dot(g.CobraMethod()).Call(jen.Lit(g.Flag())), jen.Err().Op("!=").Nil()).
			Block(
				jen.Return().List(returnId, jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+g.Flag()+"\" from command flags: %w"), jen.Err())),
			)
	}

	return nil
}
