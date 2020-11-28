package articlesim

import (
	"errors"
)

type (
	ArticleID        int
	DuplicateGroupID int
)

type Article struct {
	ID               ArticleID
	Content          string
	DuplicateIDs     []ArticleID
	IsUnique         bool
	DuplicateGroupID DuplicateGroupID
}

type DuplicateGroup struct {
	DuplicateGroupID DuplicateGroupID
	ArticleID        ArticleID
}

var ErrArticleNotFound = errors.New("article not found")
