package similarity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevenshtein_Compare(t *testing.T) {
	for name, tc := range map[string]struct {
		sequenceA []Element
		sequenceB []Element
		expected  float64
	}{
		"when empty sequences": {
			sequenceA: []Element{},
			sequenceB: []Element{},
			expected:  1.0,
		},
		"when one empty sequence": {
			sequenceA: []Element{},
			sequenceB: []Element{Element(5)},
			expected:  0.0,
		},
		"when non empty sequences": {
			sequenceA: []Element{Element(2), Element(3)},
			sequenceB: []Element{Element(1), Element(2), Element(3), Element(4)},
			expected:  0.5,
		},
	} {
		t.Run(name, func(t *testing.T) {
			lev := NewLevenshtein()

			res := lev.Compare(tc.sequenceA, tc.sequenceB, DefaultCompareFn)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestLevenshtein_Distance(t *testing.T) {
	for name, tc := range map[string]struct {
		sequenceA []Element
		sequenceB []Element
		expected  int
	}{
		"when empty sequences": {
			sequenceA: []Element{},
			sequenceB: []Element{},
			expected:  0,
		},
		"when one empty sequence": {
			sequenceA: []Element{},
			sequenceB: []Element{Element(5)},
			expected:  1,
		},
		"when non empty sequences": {
			sequenceA: []Element{Element(2), Element(3)},
			sequenceB: []Element{Element(1), Element(2), Element(3), Element(4)},
			expected:  2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			lev := NewLevenshtein()

			res := lev.Distance(tc.sequenceA, tc.sequenceB, DefaultCompareFn)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestLevenshtein_CompareWord(t *testing.T) {
	lev := NewLevenshtein()

	res := lev.CompareWord("carry", "bark")

	assert.Equal(t, 0.4, res)
}

func TestLevenshtein_DistanceWord(t *testing.T) {
	lev := NewLevenshtein()

	res := lev.DistanceWord("carry", "bark")

	assert.Equal(t, 3, res)
}

func TestLevenshtein_CompareSentence(t *testing.T) {
	lev := NewLevenshtein()

	res := lev.CompareSentence([]string{"one", "two", "three", "three", "four"},
		[]string{"five", "two", "three", "Three"})

	assert.Equal(t, 0.4, res)
}

func TestLevenshtein_DistanceSentence(t *testing.T) {
	lev := NewLevenshtein()

	res := lev.DistanceSentence([]string{"one", "two", "three", "three", "four"},
		[]string{"five", "two", "three", "Three"})

	assert.Equal(t, 3, res)
}
