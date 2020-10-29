package similarity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilarity_Similarity(t *testing.T) {
	for name, tc := range map[string]struct {
		articleA string
		articleB string
		expected float64
	}{
		"when empty strings": {
			articleA: "",
			articleB: "",
			expected: 1.0,
		},
		"when strings with only whitespaces": {
			articleA: " ",
			articleB: "\t\n\r\f",
			expected: 1.0,
		},
		"when equal strings": {
			articleA: "hello, world",
			articleB: "hello, world",
			expected: 1.0,
		},
		"when one word strings with different case": {
			articleA: "Hello",
			articleB: "HELLO",
			expected: 1.0,
		},
		"when articles with punctuation": {
			articleA: "hello world",
			articleB: "!? hello, - world,",
			expected: 1.0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			sim := NewSimilarity()

			res := sim.Similarity(tc.articleA, tc.articleB)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestSimilarity_IsSimilar(t *testing.T) {
	sim := Similarity{
		Threshold: 0.8,
	}

	res := sim.IsSimilar("hello a very beautiful world", "hello a beautiful world")

	assert.True(t, res)
}
