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
	switch KindOf(s.Field) {
	case FieldKindTCloudTag, FieldKindStringMap:
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

	switch KindOf(s.Field) {
	case FieldKindStringMap, FieldKindTCloudTag:
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
	switch KindOf(s.Field) {
	case FieldKindTCloudTag, FieldKindStringMap, FieldKindNative:
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

	switch KindOf(g.Field) {
	case FieldKindNative:
		return jen.If(jen.List(id.Dot(g.Field.Name), jen.Err()).Op("=").
			Id("cf").Dot("flags").Dot(g.CobraMethod()).Call(jen.Lit(g.Flag())), jen.Err().Op("!=").Nil()).
			Block(
				jen.Return().List(returnId, jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+g.Flag()+"\" from command flags: %w"), jen.Err())),
			)
	case FieldKindStruct, FieldKindTCloudTag, FieldKindStringMap:
		return jen.If(jen.List(id.Dot(g.Field.Name), jen.Err()).Op("=").
			Id("cf").Dot(changecase.Camel(path.Join("Get", g.Prefix, g.Field.Name))).Call(), jen.Err().Op("!=").Nil()).
			Block(jen.Return().List(returnId, jen.Err()))
	}

	return nil
}

type PointerGetterCall struct {
	Prefix  string
	Struct  *projscan.Struct
	Pointer bool
	Field   *projscan.Field
}

func (g *PointerGetterCall) CobraMethod() string {
	var suffix string
	if g.Field.Array {
		suffix = "Slice"
	}

	return changecase.Pascal(path.Join("get", g.Field.Type.String(), suffix))
}

func (g *PointerGetterCall) Flag() string {
	return changecase.Param(path.Join(g.Prefix, g.Field.Name))
}

func (g *PointerGetterCall) Statement() *jen.Statement {
	structName := changecase.Camel(g.Struct.Name)
	fieldName := g.Field.Name
	flagValue := "flagValue"

	returnId := jen.Nil()
	if !g.Pointer {
		returnId = jen.Id(structName)
	}

	switch KindOf(g.Field) {
	case FieldKindNative:
		return jen.If(jen.List(jen.Id(flagValue), jen.Err()).Op(":=").
			Id("cf").Dot("flags").Dot(g.CobraMethod()).Call(jen.Lit(g.Flag())), jen.Err().Op("!=").Nil()).
			Block(
				jen.Return().List(returnId, jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+g.Flag()+"\" from command flags: %w"), jen.Err())),
			).Else().If(g.CompareToDefaultValue(jen.Id(flagValue).Op("!=")).Op("&&").Id(structName).Op("==").Nil()).
			Block(
				jen.Id(structName).Op("=").Op("&").Qual(g.Struct.Package.Path, g.Struct.Name).Values(jen.Id(fieldName).Op(":").Id(flagValue)),
			).Else().If(g.CompareToDefaultValue(jen.Id(flagValue).Op("!="))).
			Block(
				jen.Id(structName).Dot(fieldName).Op("=").Id(flagValue),
			)
	case FieldKindStruct, FieldKindTCloudTag, FieldKindStringMap:
		if g.Field.Pointer {
			return jen.If(jen.List(jen.Id(flagValue), jen.Err()).Op(":=").
				Id("cf").Dot(changecase.Camel(path.Join("Get", g.Prefix, fieldName))).Call(), jen.Err().Op("!=").Nil()).
				Block(
					jen.Return().List(returnId, jen.Err()),
				).Else().If(g.CompareToDefaultValue(jen.Id(flagValue).Op("!=")).Op("&&").Id(structName).Op("==").Nil()).
				Block(
					jen.Id(structName).Op("=").Op("&").Qual(g.Struct.Package.Path, g.Struct.Name).Values(jen.Id(fieldName).Op(":").Id(flagValue)),
				).Else().If(g.CompareToDefaultValue(jen.Id(flagValue).Op("!="))).
				Block(
					jen.Id(structName).Dot(fieldName).Op("=").Id(flagValue),
				)
		} else {
			return jen.If(jen.List(jen.Id(flagValue), jen.Err()).Op(":=").
				Id("cf").Dot(changecase.Camel(path.Join("Get", g.Prefix, fieldName))).Call(), jen.Err().Op("!=").Nil()).
				Block(
					jen.Return().List(returnId, jen.Err()),
				).Else().If(jen.Id(structName).Op("==").Nil()).
				Block(
					jen.Id(structName).Op("=").Op("&").Qual(g.Struct.Package.Path, g.Struct.Name).Values(jen.Id(fieldName).Op(":").Id(flagValue)),
				).Else().
				Block(
					jen.Id(structName).Dot(fieldName).Op("=").Id(flagValue),
				)
		}
	}

	return nil
}

func (g *PointerGetterCall) CompareToDefaultValue(statement *jen.Statement) *jen.Statement {
	if g.Field.Array {
		return statement.Nil()
	}

	switch g.Field.Type {
	case projscan.FieldTypeString:
		return statement.Lit("")
	case projscan.FieldTypeBool:
		return statement.False()
	case projscan.FieldTypeFloat32, projscan.FieldTypeFloat64:
		return statement.Lit(0.0)
	case projscan.FieldTypeInt, projscan.FieldTypeInt8, projscan.FieldTypeInt16, projscan.FieldTypeInt32, projscan.FieldTypeInt64:
		return statement.Lit(0)
	default:
		return statement.Nil()
	}
}
