package code

import (
	"fmt"
	"os"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
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

func (f *FlagSource) Bytes() []byte {
	file := f.File()
	return []byte(fmt.Sprintf("%#v", file))
}

func (f *FlagSource) Print() {
	bytes := f.Bytes()
	fmt.Println(string(bytes))
}

func (f *FlagSource) WriteFile(directory string, file string) error {
	bytes := f.Bytes()
	if err := os.WriteFile(path.Join(directory, file), bytes, 0o644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
