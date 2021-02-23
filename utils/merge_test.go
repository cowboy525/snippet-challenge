package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/topoface/snippet-challenge/utils"
)

// Test merging maps alone. This isolates the complexity of merging maps from merging maps recursively in
// a struct/ptr/etc.
// Remember that for our purposes, "merging" means replacing base with patch if patch is /anything/ other than nil.
func TestMergeWithMaps(t *testing.T) {
	t.Run("merge maps where patch is longer", func(t *testing.T) {
		m1 := map[string]int{"this": 1, "is": 2, "a map": 3}
		m2 := map[string]int{"this": 1, "is": 3, "a second map": 3, "another key": 4}

		expected := map[string]int{"this": 1, "is": 3, "a second map": 3, "another key": 4}
		merged, err := mergeStringIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})

	t.Run("merge maps where base is longer", func(t *testing.T) {
		m1 := map[string]int{"this": 1, "is": 2, "a map": 3, "with": 4, "more keys": -12}
		m2 := map[string]int{"this": 1, "is": 3, "a second map": 3}
		expected := map[string]int{"this": 1, "is": 3, "a second map": 3}

		merged, err := mergeStringIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})

	t.Run("merge maps where base is empty", func(t *testing.T) {
		m1 := make(map[string]int)
		m2 := map[string]int{"this": 1, "is": 3, "a second map": 3, "another key": 4}

		expected := map[string]int{"this": 1, "is": 3, "a second map": 3, "another key": 4}
		merged, err := mergeStringIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})

	t.Run("merge maps where patch is empty", func(t *testing.T) {
		m1 := map[string]int{"this": 1, "is": 3, "a map": 3, "another key": 4}
		var m2 map[string]int
		expected := map[string]int{"this": 1, "is": 3, "a map": 3, "another key": 4}

		merged, err := mergeStringIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})

	t.Run("merge map[string]*int patch with different keys and values", func(t *testing.T) {
		m1 := map[string]*int{"this": newInt(1), "is": newInt(3), "a key": newInt(3)}
		m2 := map[string]*int{"this": newInt(2), "is": newInt(3), "a key": newInt(4)}
		expected := map[string]*int{"this": newInt(2), "is": newInt(3), "a key": newInt(4)}

		merged, err := mergeStringPtrIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})

	t.Run("merge map[string]*int patch has nil keys -- doesn't matter, maps overwrite completely", func(t *testing.T) {
		m1 := map[string]*int{"this": newInt(1), "is": newInt(3), "a key": newInt(3)}
		m2 := map[string]*int{"this": newInt(1), "is": nil, "a key": newInt(3)}
		expected := map[string]*int{"this": newInt(1), "is": nil, "a key": newInt(3)}

		merged, err := mergeStringPtrIntMap(m1, m2)
		require.NoError(t, err)

		assert.Equal(t, expected, merged)
	})
}

func mergeStringIntMap(base, patch map[string]int) (map[string]int, error) {
	ret, err := utils.Merge(base, patch, nil)
	if err != nil {
		return nil, err
	}
	retTS := ret.(map[string]int)
	return retTS, nil
}

func mergeStringPtrIntMap(base, patch map[string]*int) (map[string]*int, error) {
	ret, err := utils.Merge(base, patch, nil)
	if err != nil {
		return nil, err
	}
	retTS := ret.(map[string]*int)
	return retTS, nil
}

func newInt(n int) *int { return &n }
