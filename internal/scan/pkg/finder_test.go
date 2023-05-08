package pkg_test

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/totvs-cloud/pflagstruct/internal/scan/pkg"
	"github.com/totvs-cloud/pflagstruct/internal/scan/proj"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

func TestFinder_FindPackage(t *testing.T) {
	t.Run("", func(t *testing.T) {
		svc := newPackageFinder()
		singlepkg, err := svc.FindPackageByDirectory("../../../_test/testdata/foo")

		require.NoError(t, err)
		require.NotEmpty(t, singlepkg)
	})
	t.Run("", func(t *testing.T) {
		svc := newPackageFinder()
		path := "../../../_test/testdata/foo"

		singlepkg, err := svc.FindPackageByDirectory(path)
		require.NoError(t, err)
		require.NotEqual(t, path, singlepkg.Directory)
		require.Equal(t, "github.com/totvs-cloud/pflagstruct/_test/testdata/foo", singlepkg.Path)
	})
}

func TestFinder_FindPackageByPathAndProject(t *testing.T) {
	t.Run("", func(t *testing.T) {
		scanner := syntree.NewScanner(token.NewFileSet())
		var projSvc projscan.ProjectFinder = proj.NewFinder(scanner)
		var pkgSvc projscan.PackageFinder = pkg.NewFinder(scanner, projSvc)

		project, err := projSvc.FindProjectByDirectory("../../../_test/testdata/foo")
		require.NoError(t, err)

		singlepkg, err := pkgSvc.FindPackageByPathAndProject("github.com/apirator/apirator/api/v1alpha1", project)
		require.NoError(t, err)
		require.NotEmpty(t, singlepkg)
	})
}

func newPackageFinder() projscan.PackageFinder {
	scanner := syntree.NewScanner(token.NewFileSet())
	Finder := proj.NewFinder(scanner)
	pkgFinder := pkg.NewFinder(scanner, Finder)

	return pkgFinder
}
