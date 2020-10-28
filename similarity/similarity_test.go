package similarity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devchallenge/article-similarity/similarity"
)

func TestIsSimilar(t *testing.T) {
	assert.False(t, similarity.IsSimilar("times new roman", "times roman"))
	assert.False(t, similarity.IsSimilar("hello world", "world hello"))
	assert.True(t, similarity.IsSimilar("hello a beautiful world ever", "hello a beautiful great world ever"))
}
