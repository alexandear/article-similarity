package similarity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devchallenge/article-similarity/similarity"
)

func TestIsSimilar(t *testing.T) {
	article1 := "hello world"
	article2 := "world hello"

	assert.True(t, similarity.IsSimilar(article1, article2))
}
