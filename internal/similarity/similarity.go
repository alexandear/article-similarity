package similarity

import (
	"strings"

	"github.com/devchallenge/article-similarity/internal/util"
)

type Similarity struct {
	Threshold float64
}

const defaultSimilarityThreshold = 0.95

// NewSimilarity returns Similarity with 0.95 threshold.
func NewSimilarity() *Similarity {
	return &Similarity{
		Threshold: defaultSimilarityThreshold,
	}
}

func (s *Similarity) Similarity(articleA, articleB string) float64 {
	lev := NewLevenshtein()

	sim := lev.CompareSentence(tokens(articleA), tokens(articleB))

	return sim
}

// tokens removes non-alphanumeric character, splits by whitespace characters and returns lowercase words.
func tokens(article string) []string {
	a := string(util.Strip([]byte(article)))
	a = strings.ToLower(a)

	return strings.Fields(a)
}

func (s *Similarity) IsSimilar(articleA, articleB string) bool {
	return s.Similarity(articleA, articleB) >= s.Threshold
}
