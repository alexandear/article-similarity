package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type article struct {
	ID           int    `bson:"id"`
	Content      string `bson:"content"`
	DuplicateIDs []int  `bson:"duplicate_ids"`
	IsUnique     bool   `bson:"is_unique"`
}

type autoincrement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Collection string             `bson:"collection"`
	Counter    int                `bson:"counter"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
