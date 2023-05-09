package main

import (
	"go/token"
	"os"
	"time"

	"github.com/enescakir/emoji"
	"github.com/gookit/color"
	"github.com/lmittmann/tint"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/totvs-cloud/pflagstruct/internal/code"
	scanfld "github.com/totvs-cloud/pflagstruct/internal/scan/fld"
	scanpkg "github.com/totvs-cloud/pflagstruct/internal/scan/pkg"
	scanproj "github.com/totvs-cloud/pflagstruct/internal/scan/proj"
	scanst "github.com/totvs-cloud/pflagstruct/internal/scan/st"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"golang.org/x/exp/slog"
)

var (
	directory, structName, destination string
	debug                              bool
)

type Generator interface {
	Generate(directory string, structName string, destination string) error
}

func NewGenerator() Generator {
	scanner := syntree.NewScanner(token.NewFileSet())
	projects := scanproj.NewFinder(scanner)
	packages := scanpkg.NewFinder(scanner, projects)
	structs := scanst.NewFinder(scanner, projects, packages)
	fields := scanfld.NewFinder(packages, projects, structs)

	return code.NewGenerator(fields, packages, projects, structs)
}

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "flagstruct",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			generator := NewGenerator()
			return generator.Generate(directory, structName, destination)
		},
	}

	const (
		directoryFlagName   = "directory"
		structNameFlagName  = "struct-name"
		destinationFlagName = "destination"
		debugFlagName       = "debug"
	)

	cmd.Flags().StringVar(&directory, directoryFlagName, "", "Specifies the path where the tool will search for the source file containing the struct definition")
	cmd.Flags().StringVar(&structName, structNameFlagName, "", "Specifies the name of the struct that will be generated in the code. This flag is required.")
	cmd.Flags().StringVar(&destination, destinationFlagName, "", "Specifies the file name and path where the generated code will be saved. This flag is required.")
	cmd.Flags().BoolVar(&debug, debugFlagName, false, "Enables debug mode, which will print additional information during the code generation process to help with troubleshooting. If this flag is not set, the tool will run in normal mode without additional debug information.")

	if err := cmd.MarkFlagRequired(directoryFlagName); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := cmd.MarkFlagRequired(structNameFlagName); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := cmd.MarkFlagRequired(destinationFlagName); err != nil {
		return nil, errors.WithStack(err)
	}

	return cmd, nil
}

func fatal(err error) {
	if debug {
		color.Redf("%s %+v", emoji.CrossMark, err)
		os.Exit(1)
	}

	color.Redf("%s %s", emoji.CrossMark, err)
	os.Exit(1)
}

func main() {
	options := tint.Options{TimeFormat: time.Kitchen}
	if debug {
		options.AddSource = true
		options.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(options.NewHandler(os.Stderr)))

	cmd, err := NewCommand()
	if err != nil {
		fatal(err)
	}

	if err = cmd.Execute(); err != nil {
		fatal(err)
	}
}
