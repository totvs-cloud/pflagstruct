package code

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/samber/lo"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

type FlagSource struct {
	Package *projscan.Package
	Blocks  []Block

	variables []string
	imports   map[string]string
}

func (f *FlagSource) ImportName(path, name string) {
	if f.imports == nil {
		f.imports = map[string]string{
			"github.com/spf13/cobra": "cobra",
			"github.com/spf13/pflag": "pflag",
		}
	}

	f.imports[path] = name
	f.variables = lo.Uniq(append(f.variables, name))
}

func (f *FlagSource) File() *jen.File {
	file := jen.NewFilePathName(f.Package.Path, f.Package.Name)
	file.ImportNames(f.imports)

	for _, block := range f.Blocks {
		file.Add(block.Statement())
	}

	return file
}

func (f *FlagSource) Print() {
	file := f.File()
	bytes := []byte(fmt.Sprintf("%#v", file))
	fmt.Println(string(bytes))
}
