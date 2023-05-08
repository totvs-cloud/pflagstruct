package fld

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/totvs-cloud/pflagstruct/internal/scan/pkg"
	"github.com/totvs-cloud/pflagstruct/internal/scan/proj"
	"github.com/totvs-cloud/pflagstruct/internal/scan/st"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
)

func TestFinder_FindFieldsByStruct(t *testing.T) {
	t.Run("", func(t *testing.T) {
		scanner := syntree.NewScanner(token.NewFileSet())
		projsvc := proj.NewFinder(scanner)
		pkgsvc := pkg.NewFinder(scanner, projsvc)
		stsvc := st.NewFinder(scanner, projsvc, pkgsvc)
		fldsvc := NewFinder(pkgsvc, projsvc, stsvc)

		st, err := stsvc.FindStructByDirectoryAndName("../../../_test/testdata/bar", "Quuz")
		require.NoError(t, err)
		require.NotEmpty(t, st)

		flds, err := fldsvc.FindFieldsByStruct(st)
		require.NoError(t, err)
		require.NotEmpty(t, st)
		for _, fld := range flds {
			require.NotEmpty(t, fld)
		}
	})
}
