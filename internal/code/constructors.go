package code

import "github.com/dave/jennifer/jen"

type ConstructorForFlags struct{}

func (c *ConstructorForFlags) Statement() *jen.Statement {
	args := []jen.Code{
		jen.Id("cmd").Op("*").Qual("github.com/spf13/cobra", "Command"),
	}
	returns := []jen.Code{
		jen.Op("*").Id("CommandFlags"),
	}

	return jen.Func().Id("persistentFlagsOf").Params(args...).Params(returns...).Block(
		jen.Return().Op("&").Id("CommandFlags").Values(
			jen.Id("flags").Op(":").Id("cmd").Dot("PersistentFlags").Call(),
		),
	)
}

type ConstructorForPersistentFlags struct{}

func (c *ConstructorForPersistentFlags) Statement() *jen.Statement {
	args := []jen.Code{
		jen.Id("cmd").Op("*").Qual("github.com/spf13/cobra", "Command"),
	}
	returns := []jen.Code{
		jen.Op("*").Id("CommandFlags"),
	}

	return jen.Func().Id("flagsOf").Params(args...).Params(returns...).Block(
		jen.Return().Op("&").Id("CommandFlags").Values(
			jen.Id("flags").Op(":").Id("cmd").Dot("Flags").Call(),
		),
	)
}
