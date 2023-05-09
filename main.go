package main

import (
	"fmt"
	"go/token"
	"os"
	"time"

	"github.com/enescakir/emoji"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	directory, pkgPath, structName, destination string
	debug                                       bool
)

func NewCommand() (*cobra.Command, error) {
	const (
		directoryFlagName   = "directory"
		packageFlagName     = "package"
		structNameFlagName  = "struct-name"
		destinationFlagName = "destination"
		debugFlagName       = "debug"
	)

	cmd := &cobra.Command{
		Use:           "flagstruct",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := map[string]string{
				"--" + directoryFlagName:   directory,
				"--" + packageFlagName:     pkgPath,
				"--" + structNameFlagName:  structName,
				"--" + destinationFlagName: destination,
			}
			err := validation.Validate(flags,
				validation.Map(
					validation.Key("--"+directoryFlagName),
					validation.Key("--"+packageFlagName, validation.Required.When(directory == "").Error(fmt.Sprintf("either %s or %s is required.", "--"+packageFlagName, "--"+directoryFlagName))),
					validation.Key("--"+structNameFlagName, validation.Required),
					validation.Key("--"+destinationFlagName, validation.Required),
				),
			)
			if err != nil {
				return errors.WithStack(err)
			}

			scanner := syntree.NewScanner(token.NewFileSet())
			projects := scanproj.NewFinder(scanner)
			packages := scanpkg.NewFinder(scanner, projects)
			structs := scanst.NewFinder(scanner, projects, packages)
			fields := scanfld.NewFinder(packages, projects, structs)

			if directory != "" {
				return code.NewGenerator(fields, packages, projects, structs).Generate(directory, structName, destination)
			}

			proj, err := projects.FindProjectByDirectory(destination)
			if err != nil {
				return err
			}

			pkg, err := packages.FindPackageByPathAndProject(pkgPath, proj)
			if err != nil {
				return err
			}

			return code.NewGenerator(fields, packages, projects, structs).Generate(pkg.Directory, structName, destination)
		},
	}

	cmd.Flags().StringVar(&directory, directoryFlagName, "", "Specifies the path where the tool will search for the source file containing the struct definition")
	cmd.Flags().StringVar(&pkgPath, packageFlagName, "", "Specifies the path where the tool will search for the source file containing the struct definition")
	cmd.Flags().StringVar(&structName, structNameFlagName, "", "Specifies the name of the struct that will be generated in the code. This flag is required.")
	cmd.Flags().StringVar(&destination, destinationFlagName, ".", "Specifies the file name and path where the generated code will be saved. This flag is required.")
	cmd.Flags().BoolVar(&debug, debugFlagName, false, "Enables debug mode, which will print additional information during the code generation process to help with troubleshooting. If this flag is not set, the tool will run in normal mode without additional debug information.")

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
