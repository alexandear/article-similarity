package similarity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilarity_Similarity(t *testing.T) {
	for name, tc := range map[string]struct {
		idA      int
		contentA string
		idB      int
		contentB string
		expected float64
	}{
		"when empty strings": {
			idA:      1,
			contentA: "",
			contentB: "",
			idB:      2,
			expected: 1.0,
		},
		"when strings with only whitespaces": {
			idA:      1,
			contentA: " ",
			idB:      2,
			contentB: "\t\n\r\f",
			expected: 1.0,
		},
		"when equal strings": {
			idA:      1,
			contentA: "hello, world",
			idB:      2,
			contentB: "hello, world",
			expected: 1.0,
		},
		"when one word strings with different case": {
			idA:      1,
			contentA: "Hello",
			idB:      2,
			contentB: "HELLO",
			expected: 1.0,
		},
		"when contents with punctuation": {
			idA:      1,
			contentA: "hello world",
			idB:      2,
			contentB: "!? hello, - world,",
			expected: 1.0,
		},
		"when contents with articles": {
			idA:      1,
			contentA: "hello the world",
			idB:      2,
			contentB: "hello a world,",
			expected: 1.0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sim := NewSimilarity(t.Logf, 0.95, IrregularVerb{})

			res := sim.Similarity(tc.idA, tc.contentA, tc.idB, tc.contentB)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestSimilarity_IsSimilar(t *testing.T) {
	sim := NewSimilarity(t.Logf, 0.7, IrregularVerb{})

	res := sim.IsSimilar(1, "hello a very beautiful world", 2, "hello beautiful world")

	assert.True(t, res)
}
