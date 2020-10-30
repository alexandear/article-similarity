package similarity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilarity_Similarity(t *testing.T) {
	for name, tc := range map[string]struct {
		contentA string
		contentB string
		expected float64
	}{
		"when empty strings": {
			contentA: "",
			contentB: "",
			expected: 1.0,
		},
		"when strings with only whitespaces": {
			contentA: " ",
			contentB: "\t\n\r\f",
			expected: 1.0,
		},
		"when equal strings": {
			contentA: "hello, world",
			contentB: "hello, world",
			expected: 1.0,
		},
		"when one word strings with different case": {
			contentA: "Hello",
			contentB: "HELLO",
			expected: 1.0,
		},
		"when contents with punctuation": {
			contentA: "hello world",
			contentB: "!? hello, - world,",
			expected: 1.0,
		},
		"when contents with articles": {
			contentA: "hello the world",
			contentB: "hello a world,",
			expected: 1.0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sim := NewSimilarity(func(string, ...interface{}) {}, 0.95)

			res := sim.Similarity(tc.contentA, tc.contentB)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestSimilarity_IsSimilar(t *testing.T) {
	sim := Similarity{
		Threshold: 0.7,
	}

	res := sim.IsSimilar("hello a very beautiful world", "hello beautiful world")

	assert.True(t, res)
}
