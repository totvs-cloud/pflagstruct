package code

import (
	"path"

	"github.com/dave/jennifer/jen"
	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type SetterMethod struct {
	FlagsBuilderName string
	Struct           *projscan.Struct
	Flags            map[string][]*projscan.Field
}

func (s *SetterMethod) MethodName() string {
	return changecase.Camel(path.Join("SetUp", s.Struct.Name))
}

func (s *SetterMethod) Statement() *jen.Statement {
	receiver := jen.Id("cf").Op("*").Id(s.FlagsBuilderName)

	calls := make([]jen.Code, 0)

	for prefix, fields := range s.Flags {
		for _, field := range fields {
			calls = append(calls, (&SetterCall{
				Prefix: prefix,
				Struct: s.Struct,
				Field:  field,
			}).Statement())
		}
	}

	return jen.Func().Params(receiver).Id(s.MethodName()).Params().Block(calls...)
}

type GetterMethod struct {
	FlagsBuilderName string
	Prefix           string
	Struct           *projscan.Struct
	Pointer          bool
	Fields           []*projscan.Field
}

func (g *GetterMethod) MethodName() string {
	if g.Prefix == "" {
		return changecase.Camel(path.Join("Get", g.Struct.Name))
	}

	return changecase.Camel(path.Join("Get", g.Prefix))
}

func (g *GetterMethod) Initialization() *jen.Statement {
	if !g.Pointer {
		return nil
	}

	id := jen.Id(changecase.Camel(g.Struct.Name))

	return id.Op("=").Id("new").Call(jen.Qual(g.Struct.Package.Path, g.Struct.Name))
}

func (g *GetterMethod) ReturnType() *jen.Statement {
	id := jen.Id(changecase.Camel(g.Struct.Name))
	if g.Pointer {
		return id.Op("*").Qual(g.Struct.Package.Path, g.Struct.Name)
	}

	return id.Qual(g.Struct.Package.Path, g.Struct.Name)
}

func (g *GetterMethod) ReturnCall() *jen.Statement {
	return jen.Return().List(jen.Id(changecase.Camel(g.Struct.Name)), jen.Nil())
}

func (g *GetterMethod) Statement() *jen.Statement {
	receiver := jen.Id("cf").Op("*").Id(g.FlagsBuilderName)
	returns := []jen.Code{
		g.ReturnType(), jen.Id("err").Id("error"),
	}

	calls := []jen.Code{g.Initialization()}
	for _, field := range g.Fields {
		calls = append(calls, (&GetterCall{
			Prefix:  g.Prefix,
			Struct:  g.Struct,
			Pointer: g.Pointer,
			Field:   field,
		}).Statement())
	}

	calls = append(calls, g.ReturnCall())

	return jen.Func().Params(receiver).Id(g.MethodName()).Params().Params(returns...).Block(calls...)
}

type TagsGetterMethod struct {
	FlagsBuilderName string
	Prefix           string
	Struct           *projscan.Struct
	Pointer          bool
}

func (t *TagsGetterMethod) MethodName() string {
	if t.Prefix == "" {
		return changecase.Camel(path.Join("Get", t.Struct.Name))
	}

	return changecase.Camel(path.Join("Get", t.Prefix))
}

func (t *TagsGetterMethod) Initialization() *jen.Statement {
	id := jen.Index()
	if t.Pointer {
		id = jen.Index().Op("*")
	}

	return jen.Id("resultingTags").Op(":=").Id("make").Call(id.Qual(t.Struct.Package.Path, t.Struct.Name), jen.Lit(0), jen.Id("len").Call(jen.Id("tagStrList")))
}

func (t *TagsGetterMethod) ReturnType() *jen.Statement {
	id := jen.Id(changecase.Camel(t.Struct.Name))
	if t.Pointer {
		return id.Index().Op("*").Qual(t.Struct.Package.Path, t.Struct.Name)
	}

	return id.Index().Qual(t.Struct.Package.Path, t.Struct.Name)
}

func (t *TagsGetterMethod) ResultAssignment() *jen.Statement {
	if t.Pointer {
		return jen.Op("&").Qual(t.Struct.Package.Path, t.Struct.Name)
	}

	return jen.Qual(t.Struct.Package.Path, t.Struct.Name)
}

func (t *TagsGetterMethod) Flag() string {
	return changecase.Param(path.Join(t.Prefix))
}

func (t *TagsGetterMethod) Statement() *jen.Statement {
	receiver := jen.Id("cf").Op("*").Id(t.FlagsBuilderName)
	returns := []jen.Code{
		jen.Index().Op("*").Qual(t.Struct.Package.Path, t.Struct.Name),
		jen.Error(),
	}

	calls := []jen.Code{
		jen.List(jen.Id("tagStrList"), jen.Id("err")).Op(":=").Id("cf").Dot("flags").Dot("GetStringSlice").Call(jen.Lit(t.Flag())),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+t.Flag()+"\" from command flags: %w"), jen.Id("err")))),
		t.Initialization(),
		jen.For(jen.List(jen.Id("_"),
			jen.Id("tagStr")).Op(":=").Range().Id("tagStrList")).Block(jen.Id("parts").Op(":=").Qual("strings", "Split").Call(jen.Id("tagStr"),
			jen.Lit("=")),
			jen.If(jen.Id("len").Call(jen.Id("parts")).Op("!=").Lit(2)).Block(jen.Return().List(jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+t.Flag()+"\" from command flags: invalid format: %s"), jen.Id("tagStr")))),
			jen.Id("resultingTags").Op("=").Id("append").Call(jen.Id("resultingTags"),
				t.ResultAssignment().Values(jen.Id("Name").Op(":").Id("parts").Index(jen.Lit(0)), jen.Id("Value").Op(":").Id("parts").Index(jen.Lit(1))))),
	}

	calls = append(calls, jen.Return().List(jen.Id("resultingTags"), jen.Nil()))

	return jen.Func().Params(receiver).Id(t.MethodName()).Params().Params(returns...).Block(calls...)
}

type MapGetterMethod struct {
	FlagsBuilderName string
	Prefix           string
	Pointer          bool
}

func (t *MapGetterMethod) MethodName() string {
	return changecase.Camel(path.Join("Get", t.Prefix))
}

func (t *MapGetterMethod) Flag() string {
	return changecase.Param(path.Join(t.Prefix))
}

func (t *MapGetterMethod) Statement() *jen.Statement {
	const (
		filterStrList   = "filterStrList"
		resultingFilter = "resultingFilter"
		filterStr       = "filterStr"
		parts           = "parts"
	)

	receiver := jen.Id("cf").Op("*").Id(t.FlagsBuilderName)
	returns := []jen.Code{
		jen.Map(jen.String()).String(),
		jen.Error(),
	}

	calls := []jen.Code{
		jen.List(jen.Id(filterStrList), jen.Id("err")).Op(":=").Id("cf").Dot("flags").Dot("GetStringSlice").Call(jen.Lit(t.Flag())),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+t.Flag()+"\" from command flags: %w"), jen.Id("err"))),
		),
		jen.Id(resultingFilter).Op(":=").Id("make").Call(jen.Map(jen.String()).String()),
		jen.For(jen.List(jen.Id("_"), jen.Id(filterStr)).Op(":=").Range().Id(filterStrList)).Block(
			jen.Id(parts).Op(":=").Qual("strings", "Split").Call(jen.Id(filterStr), jen.Lit("=")),
			jen.If(jen.Id("len").Call(jen.Id(parts)).Op("!=").Lit(2)).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error retrieving \""+t.Flag()+"\" from command flags: %w"), jen.Id("err"))),
			),
			jen.Id(resultingFilter).Index(jen.Id(parts).Index(jen.Lit(0))).Op("=").Id(parts).Index(jen.Lit(1)),
		),
		jen.Return().List(jen.Id(resultingFilter), jen.Nil()),
	}

	return jen.Func().Params(receiver).Id(t.MethodName()).Params().Params(returns...).Block(calls...)
}
