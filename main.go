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

			if directory == "" {
				proj, err := projects.FindProjectByDirectory(destination)
				if err != nil {
					return err
				}

				pkg, err := packages.FindPackageByPathAndProject(pkgPath, proj)
				if err != nil {
					return err
				}

				directory = pkg.Directory
			}

			filepath, err := code.NewGenerator(fields, packages, projects, structs).Generate(directory, structName, destination)
			if err != nil {
				return err
			}
			fmt.Printf("%s Code generated successfully! Find it at: %s\n", emoji.CheckMark, filepath)
			return nil
		},
	}

	cmd.Flags().StringVar(&structName, structNameFlagName, "", "specifies the name of the struct. This flag is required to generate code based on the provided struct")
	cmd.Flags().StringVar(&pkgPath, packageFlagName, "", "specifies the package path of the struct definition. This flag is required if the --directory flag is not informed")
	cmd.Flags().StringVar(&directory, directoryFlagName, "", "specifies the path where the source file containing the struct definition is located. This flag is required if the --package flag is not informed")
	cmd.Flags().StringVar(&destination, destinationFlagName, ".", "specifies the path where the generated code will be saved")
	cmd.Flags().BoolVar(&debug, debugFlagName, false, "enables debug mode, which provides additional output for debugging purposes")

	return cmd, nil
}

func fatal(err error) {
	if debug {
		color.Redf("%s %+v\n", emoji.CrossMark, err)
		os.Exit(1)
	}

	color.Redf("%s %s\n", emoji.CrossMark, err)
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
