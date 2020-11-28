package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/devchallenge/article-similarity/internal/model"
)

type article struct {
	ID               model.ArticleID        `bson:"id"`
	Content          string                 `bson:"content"`
	DuplicateIDs     []model.ArticleID      `bson:"duplicate_ids"`
	IsUnique         bool                   `bson:"is_unique"`
	DuplicateGroupID model.DuplicateGroupID `bson:"duplicate_group_id"`
}

type duplicateGroup struct {
	ID        model.DuplicateGroupID `bson:"id"`
	ArticleID model.ArticleID        `bson:"article_id"`
}

type autoincrement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Collection string             `bson:"collection"`
	Counter    int                `bson:"counter"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
