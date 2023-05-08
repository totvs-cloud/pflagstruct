package syntree

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner_ScanDirectory(t *testing.T) {
	scanner := NewScanner(token.NewFileSet())
	directory := "../../_test/testdata/foo"
	_, err := scanner.ScanDirectory(directory)
	require.NoError(t, err)
}
