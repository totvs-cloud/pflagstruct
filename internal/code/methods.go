package code

import (
	"path"

	"github.com/dave/jennifer/jen"
	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type SetterMethod struct {
	Struct *projscan.Struct
	Flags  map[string][]*projscan.Field
}

func (s *SetterMethod) MethodName() string {
	return changecase.Pascal(path.Join("SetUp", s.Struct.Name))
}

func (s *SetterMethod) Statement() *jen.Statement {
	receiver := jen.Id("cf").Op("*").Id("CommandFlags")

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
	Prefix  string
	Struct  *projscan.Struct
	Pointer bool
	Fields  []*projscan.Field
}

func (g *GetterMethod) MethodName() string {
	if g.Prefix == "" {
		return changecase.Pascal(path.Join("Get", g.Struct.Name))
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
	return jen.Return().List(jen.Id(changecase.Camel(g.Struct.Name)), jen.Id("nil"))
}

func (g *GetterMethod) Statement() *jen.Statement {
	receiver := jen.Id("cf").Op("*").Id("CommandFlags")
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
