
package main

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestProcessBasic(t *testing.T) {
	result, err := Process([]string{"a", "a", "b"}, Options{})
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b"}, result)
}

func TestProcessCount(t *testing.T) {
	result, err := Process([]string{"a", "a", "b"}, Options{Count: true})
	require.NoError(t, err)
	require.Equal(t, []string{"2 a", "1 b"}, result)
}

func TestProcessConflict(t *testing.T) {
	_, err := Process([]string{"a"}, Options{Count: true, Duplicates: true})
	require.Error(t, err)
}