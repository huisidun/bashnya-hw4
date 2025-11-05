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

func TestProcessDuplicates(t *testing.T) {
	result, err := Process([]string{"a", "a", "b", "c", "c"}, Options{Duplicates: true})
	require.NoError(t, err)
	require.Equal(t, []string{"a", "c"}, result)
}

func TestProcessUniqueOnly(t *testing.T) {
	result, err := Process([]string{"a", "a", "b", "c", "c"}, Options{UniqueOnly: true})
	require.NoError(t, err)
	require.Equal(t, []string{"b"}, result)
}

func TestProcessIgnoreCase(t *testing.T) {
	result, err := Process([]string{"A", "a", "B"}, Options{IgnoreCase: true})
	require.NoError(t, err)
	require.Equal(t, []string{"A", "B"}, result)
}

func TestProcessSkipFields(t *testing.T) {
	result, err := Process([]string{"x a", "y a", "z b"}, Options{SkipFields: 1})
	require.NoError(t, err)
	require.Equal(t, []string{"x a", "z b"}, result)
}

func TestProcessSkipChars(t *testing.T) {
	result, err := Process([]string{"1a", "2a", "3b"}, Options{SkipChars: 1})
	require.NoError(t, err)
	require.Equal(t, []string{"1a", "3b"}, result)
}

func TestProcessConflictFlags(t *testing.T) {
	_, err := Process([]string{"a"}, Options{Count: true, Duplicates: true})
	require.Error(t, err)
}