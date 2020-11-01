package similarity

import (
	"strings"

	"github.com/devchallenge/article-similarity/internal/util"
)

type Similarity struct {
	logger    func(format string, v ...interface{})
	threshold float64

	irregular IrregularVerb
}

func NewSimilarity(logger func(format string, v ...interface{}), threshold float64, irregular IrregularVerb,
) *Similarity {
	return &Similarity{
		logger:    logger,
		threshold: threshold,
		irregular: irregular,
	}
}

func (s *Similarity) IsSimilar(idA int, contentA string, idB int, contentB string) bool {
	s.logger("use similarity threshold: %f to compare %d and %d", s.threshold, idA, idB)

	sim := s.Similarity(idA, contentA, idB, contentB) >= s.threshold

	return sim
}

func (s *Similarity) Similarity(idA int, contentA string, idB int, contentB string) float64 {
	lev := NewLevenshtein()

	s.logger("normalizing %d", idA)

	normA := s.normalizeAndReturnWords(contentA)

	s.logger("normalizing %d", idB)

	normB := s.normalizeAndReturnWords(contentB)

	sim := lev.CompareSentence(normA, normB)

	return sim
}

// normalizeAndReturnWords removes non-alphanumeric character, splits by whitespace characters,
// removes articles (a, an, the), change verbs to infinitives and returns lowercase words.
func (s *Similarity) normalizeAndReturnWords(content string) []string {
	modContent := string(util.Strip([]byte(content)))
	modContent = strings.ToLower(modContent)
	fields := strings.Fields(modContent)

	res := make([]string, 0, len(fields))
	articles := []string{"a", "an", "the"}

	for _, t := range fields {
		if util.Contains(articles, t) {
			continue
		}

		verb := s.irregular.ToInfinitive(t)

		res = append(res, verb)
	}

	return res
}
