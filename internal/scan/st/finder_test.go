package st_test

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/totvs-cloud/pflagstruct/internal/scan/pkg"
	"github.com/totvs-cloud/pflagstruct/internal/scan/proj"
	"github.com/totvs-cloud/pflagstruct/internal/scan/st"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

func TestFinder_FindAllStructs(t *testing.T) {
	t.Run("", func(t *testing.T) {
		svc := newFinder()
		structs, err := svc.FindStructByDirectoryAndName("../../../_test/testdata/foo", "Baz2")
		require.NoError(t, err)
		require.NotEmpty(t, structs)
	})
	t.Run("", func(t *testing.T) {
		svc := newFinder()
		structs, err := svc.FindStructByDirectoryAndName("../../../_test/testdata/foo", "Baz")
		require.NoError(t, err)
		require.NotEmpty(t, structs)
	})
	t.Run("", func(t *testing.T) {
		svc := newFinder()
		structs, err := svc.FindStructByDirectoryAndName("../../../_test/testdata/foo", "Corge")
		require.NoError(t, err)
		require.NotEmpty(t, structs)
	})
	t.Run("", func(t *testing.T) {
		svc := newFinder()
		structs, err := svc.FindStructByDirectoryAndName("../../../_test/testdata/foo", "Grault")
		require.NoError(t, err)
		require.NotEmpty(t, structs)
	})
}

func newFinder() projscan.StructFinder {
	scanner := syntree.NewScanner(token.NewFileSet())
	Finder := proj.NewFinder(scanner)
	pkgFinder := pkg.NewFinder(scanner, Finder)
	stFinder := st.NewFinder(scanner, Finder, pkgFinder)

	return stFinder
}
