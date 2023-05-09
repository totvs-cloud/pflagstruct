package code

import "github.com/dave/jennifer/jen"

type Source interface {
	File() *jen.File
	Print()
}

type Block interface {
	Statement() *jen.Statement
}

type MethodBlock interface {
	Block
	MethodName() string
}
