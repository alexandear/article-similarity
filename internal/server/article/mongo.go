package article

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionArticles = "articles"
)

type article struct {
	ID        int       `bson:"id"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"created_at"`
}

type autoincrement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Collection string             `bson:"collection"`
	Counter    int                `bson:"counter"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

const (
	collectionAutoincrement = "autoincrement"
)

func (a *Article) autoincrement(ctx context.Context, collection string) (*autoincrement, error) {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(true)
	doc := &autoincrement{}
	filter := bson.M{"collection": collection}
	update := bson.M{
		"$inc": bson.M{"counter": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	if err := a.mongo.Database(a.mongoDatabase).Collection(collectionAutoincrement).
		FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}
