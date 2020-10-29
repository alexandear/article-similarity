package similarity

import (
	"github.com/devchallenge/article-similarity/internal/util"
)

// Levenshtein represents the Levenshtein metric for measuring the similarity between sequences.
//   For more information see https://en.wikipedia.org/wiki/Levenshtein_distance.
type Levenshtein struct {
	// InsertCost represents the Levenshtein cost of a character insertion.
	InsertCost int

	// InsertCost represents the Levenshtein cost of a character deletion.
	DeleteCost int

	// InsertCost represents the Levenshtein cost of a character substitution.
	ReplaceCost int
}

// NewLevenshtein returns a new Levenshtein metric.
//
// Default options:
//   InsertCost: 1
//   DeleteCost: 1
//   ReplaceCost: 1
func NewLevenshtein() *Levenshtein {
	return &Levenshtein{
		InsertCost:  1,
		DeleteCost:  1,
		ReplaceCost: 1,
	}
}

// Element is a sequence element.
type Element interface{}

// CompareFn is function to compare elements.
type CompareFn func(a, b Element) bool

var DefaultCompareFn = func(a, b Element) bool { return a == b }

// Compare returns the Levenshtein similarity of sequenceA and sequenceB. Sequences is comparing with compare function.
// The returned similarity is a number between 0 and 1. Larger similarity numbers indicate closer matches.
func (m *Levenshtein) Compare(sequenceA, sequenceB []Element, compare CompareFn) float64 {
	distance := m.Distance(sequenceA, sequenceB, compare)
	if distance == 0 {
		return 1.0
	}
	maxLen := util.Max(len(sequenceA), len(sequenceB))
	return 1 - float64(distance)/float64(maxLen)
}

// Distance returns the Levenshtein distance between sequenceA and sequenceB. Sequences is comparing with compare
// function. Lower distances indicate closer matches. A distance of 0 means the strings are identical.
func (m *Levenshtein) Distance(sequenceA, sequenceB []Element, compare CompareFn) int {
	lenA, lenB := len(sequenceA), len(sequenceB)
	if lenA == 0 && lenB == 0 {
		return 0
	}

	if lenA == 0 {
		return m.InsertCost * lenB
	}
	if lenB == 0 {
		return m.DeleteCost * lenA
	}

	prevCol := make([]int, lenB+1)
	for i := 0; i <= lenB; i++ {
		prevCol[i] = i
	}

	col := make([]int, lenB+1)
	for i := 0; i < lenA; i++ {
		col[0] = i + 1
		for j := 0; j < lenB; j++ {
			delCost := prevCol[j+1] + m.DeleteCost
			insCost := col[j] + m.InsertCost

			subCost := prevCol[j]
			if !compare(sequenceA[i], sequenceB[j]) {
				subCost += m.ReplaceCost
			}

			col[j+1] = util.Min(delCost, insCost, subCost)
		}

		col, prevCol = prevCol, col
	}

	return prevCol[lenB]
}

// CompareWord returns the Levenshtein similarity between wordA and wordB strings.
// The function is a specialization of Compare for characters.
func (m *Levenshtein) CompareWord(wordA, wordB string) float64 {
	return m.Compare(stringToElementSlice(wordA), stringToElementSlice(wordB), DefaultCompareFn)
}

// DistanceWord returns the Levenshtein distance between wordA and wordB strings.
// The function is a specialization of Distance for characters.
func (m *Levenshtein) DistanceWord(wordA, wordB string) int {
	return m.Distance(stringToElementSlice(wordA), stringToElementSlice(wordB), DefaultCompareFn)
}

// CompareSentence returns the Levenshtein similarity between sentenceA and sentenceB sentences.
// Sentence consists from words. The function is a specialization of Compare for strings with
// case sensitive strings comparing.
func (m *Levenshtein) CompareSentence(sentenceA, sentenceB []string) float64 {
	return m.Compare(stringSliceToElementSlice(sentenceA), stringSliceToElementSlice(sentenceB),
		DefaultCompareFn)
}

// DistanceSentence returns the Levenshtein distance between sentenceA and sentenceB sentences.
// Sentence consists from words. The function is a specialization of Distance for strings with
// case sensitive strings comparing.
func (m *Levenshtein) DistanceSentence(sentenceA, sentenceB []string) int {
	return m.Distance(stringSliceToElementSlice(sentenceA), stringSliceToElementSlice(sentenceB),
		DefaultCompareFn)
}

func stringToElementSlice(str string) []Element {
	res := make([]Element, len(str))

	for i, r := range str {
		res[i] = r
	}

	return res
}

func stringSliceToElementSlice(slice []string) []Element {
	res := make([]Element, len(slice))

	for i := range slice {
		res[i] = slice[i]
	}

	return res
}
