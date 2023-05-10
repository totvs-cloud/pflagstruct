package code

import (
	"path"

	"github.com/dave/jennifer/jen"
	changecase "github.com/ku/go-change-case"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type SetUpConstructor struct {
	FlagsBuilderName string
	Struct           *projscan.Struct
}

func (c *SetUpConstructor) MethodName() string {
	return changecase.Pascal(path.Join("SetUp", c.Struct.Name, "to", "flags"))
}

func (c *SetUpConstructor) Statement() *jen.Statement {
	args := []jen.Code{
		jen.Id("flags").Op("*").Qual("github.com/spf13/pflag", "FlagSet"),
	}

	methodCall := changecase.Camel(path.Join("setUp", c.Struct.Name))

	return jen.Func().Id(c.MethodName()).Params(args...).Block(
		jen.Parens(
			jen.Op("&").Id(c.FlagsBuilderName).Values(jen.Id("flags").Op(":").Id("flags")),
		).
			Dot(methodCall).Call(),
	)
}

type GetConstructor struct {
	FlagsBuilderName string
	Struct           *projscan.Struct
}

func (g *GetConstructor) MethodName() string {
	return changecase.Pascal(path.Join("Get", g.Struct.Name, "from", "flags"))
}

func (g *GetConstructor) Statement() *jen.Statement {
	args := []jen.Code{
		jen.Id("flags").Op("*").Qual("github.com/spf13/pflag", "FlagSet"),
	}
	returns := []jen.Code{
		jen.Op("*").Qual(g.Struct.Package.Path, g.Struct.Name),
		jen.Error(),
	}

	methodCall := changecase.Camel(path.Join("get", g.Struct.Name))

	return jen.Func().Id(g.MethodName()).Params(args...).Params(returns...).Block(
		jen.Return().Parens(
			jen.Op("&").Id(g.FlagsBuilderName).Values(jen.Id("flags").Op(":").Id("flags")),
		).
			Dot(methodCall).Call(),
	)
}
