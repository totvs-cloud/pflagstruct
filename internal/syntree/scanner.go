package syntree

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"github.com/totvs-cloud/pflagstruct/internal/dir"
)

// Scanner is a struct that contains a fileset and is used to scan directories for Go files.
type Scanner struct {
	fset *token.FileSet
}

// NewScanner creates a new instance of Scanner.
func NewScanner(fset *token.FileSet) *Scanner {
	return &Scanner{fset: fset}
}

// ScanDirectory scans a directory for Go files and returns a map with the file names as keys and the corresponding
// AST nodes as values.
func (s *Scanner) ScanDirectory(directory string) (map[string]*ast.File, error) {
	// Get the absolute path of the directory.
	directory, err := dir.AbsolutePath(directory)
	if err != nil {
		return nil, err
	}

	// Parse the directory using the file filter and with comments enabled.
	pkgs, err := parser.ParseDir(s.fset, directory, filterTestingFiles, parser.ParseComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Create a map of the parsed files.
	files := make(map[string]*ast.File)

	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			files[name] = file
		}
	}

	return files, nil
}

// filterTestingFiles is a function that is used to filter files during parsing.
func filterTestingFiles(info fs.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}
