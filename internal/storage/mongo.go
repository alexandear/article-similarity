package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	internalErrors "github.com/devchallenge/article-similarity/internal/errors"
	"github.com/devchallenge/article-similarity/internal/model"
)

const (
	maxArticles = 100

	collectionArticles      = "articles"
	collectionAutoincrement = "autoincrement"
)

type Storage struct {
	mc *mongo.Client
	db *mongo.Database
}

func New(mc *mongo.Client, database string) *Storage {
	return &Storage{
		mc: mc,
		db: mc.Database(database),
	}
}

func (s *Storage) NextArticleID(ctx context.Context) (int, error) {
	inc, err := s.autoincrement(ctx, collectionArticles)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return inc.Counter, nil
}

func (s *Storage) CreateArticle(ctx context.Context, id int, content string) error {
	a := article{
		ID:        id,
		Content:   content,
		CreatedAt: time.Now(),
	}

	ma, err := bson.Marshal(&a)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	if _, err := s.db.Collection(collectionArticles).InsertOne(ctx, ma); err != nil {
		return errors.Wrap(err, "failed to insert")
	}

	return nil
}

func (s *Storage) ArticleByID(ctx context.Context, id int) (model.Article, error) {
	res := s.db.Collection(collectionArticles).FindOne(ctx, bson.D{{Key: "id", Value: id}})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return model.Article{}, errors.Wrap(internalErrors.ErrNotFound, "not found")
	}

	if res.Err() != nil {
		return model.Article{}, errors.Wrap(res.Err(), "failed to find")
	}

	art := article{}
	if err := res.Decode(&art); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to decode")
	}

	return model.Article{
		ID:           art.ID,
		Content:      art.Content,
		DuplicateIDs: nil,
	}, nil
}

func (s *Storage) AllArticles(ctx context.Context) ([]model.Article, error) {
	articles := make([]model.Article, 0, maxArticles)
	collection := s.db.Collection(collectionArticles)

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all documents")
	}

	for cur.TryNext(ctx) {
		art := &article{}

		if err := cur.Decode(art); err != nil {
			return nil, errors.Wrap(err, "failed to cursor decode")
		}

		articles = append(articles, model.Article{
			ID:           art.ID,
			Content:      art.Content,
			DuplicateIDs: nil,
		})
	}

	return articles, nil
}

func (s *Storage) autoincrement(ctx context.Context, collection string) (*autoincrement, error) {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(true)
	doc := &autoincrement{}
	filter := bson.M{"collection": collection}
	update := bson.M{
		"$inc": bson.M{"counter": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	if err := s.db.Collection(collectionAutoincrement).FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "failed to find one and update")
	}

	return doc, nil
}
