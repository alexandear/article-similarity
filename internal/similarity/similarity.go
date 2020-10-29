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

func (s *Similarity) Similarity(contentA, contentB string) float64 {
	lev := NewLevenshtein()

	sim := lev.CompareSentence(normalize(contentA), normalize(contentB))

	return sim
}

// normalize removes non-alphanumeric character, splits by whitespace characters, removes articles (a, an, the) and
// returns lowercase words.
func normalize(content string) []string {
	modContent := string(util.Strip([]byte(content)))
	modContent = strings.ToLower(modContent)
	fields := strings.Fields(modContent)

	res := make([]string, 0, len(fields))
	articles := []string{"a", "an", "the"}

	for _, t := range fields {
		if util.Contains(articles, t) {
			continue
		}

		res = append(res, t)
	}

	return res
}

func (s *Similarity) IsSimilar(contentA, contentB string) bool {
	return s.Similarity(contentA, contentB) >= s.Threshold
}
