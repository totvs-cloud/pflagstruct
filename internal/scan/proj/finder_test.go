package proj_test

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/totvs-cloud/pflagstruct/internal/scan/proj"
	"github.com/totvs-cloud/pflagstruct/internal/syntree"
	"github.com/totvs-cloud/pflagstruct/projscan"
)

func TestFinder_FindProject(t *testing.T) {
	t.Run("", func(t *testing.T) {
		svc := newProjectFinder()
		project, err := svc.FindProjectByDirectory("../../../_test/testdata/foo")
		require.NoError(t, err)
		require.NotEmpty(t, project)
	})
	t.Run("", func(t *testing.T) {
		svc := newProjectFinder()
		path := "../../../_test/testdata/foo"
		project, err := svc.FindProjectByDirectory(path)
		require.NoError(t, err)
		require.NotEqual(t, path, project.Directory)
		require.Equal(t, "github.com/totvs-cloud/pflagstruct/_test/testdata", project.ModuleName)
	})
}

func newProjectFinder() projscan.ProjectFinder {
	scanner := syntree.NewScanner(token.NewFileSet())
	Finder := proj.NewFinder(scanner)

	return Finder
}
