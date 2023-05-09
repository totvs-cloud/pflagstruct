package code

import "github.com/dave/jennifer/jen"

type CommandFlagsStruct struct{}

func (cfs *CommandFlagsStruct) Statement() *jen.Statement {
	fields := []jen.Code{
		jen.Id("flags").Op("*").Qual("github.com/spf13/pflag", "FlagSet"),
	}

	return jen.Type().Id("CommandFlags").Struct(fields...)
}
