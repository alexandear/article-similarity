package handler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Autoincrement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Collection string             `bson:"collection"`
	Counter    int                `bson:"counter"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

const (
	collectionAutoincrement = "autoincrement"
)

func (h *Handler) autoincrement(ctx context.Context, collection string) (*Autoincrement, error) {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(true)
	doc := &Autoincrement{}
	filter := bson.M{"collection": collection}
	update := bson.M{
		"$inc": bson.M{"counter": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	if err := h.mongo.Database("dev").Collection(collectionAutoincrement).
		FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}
