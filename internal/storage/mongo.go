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
	maxArticles = 1000

	collectionArticles      = "articles"
	collectionAutoincrement = "autoincrement"
)

type Storage struct {
	collectionArticle       *mongo.Collection
	collectionAutoincrement *mongo.Collection
}

func New(mc *mongo.Client, database string) *Storage {
	db := mc.Database(database)
	return &Storage{
		collectionArticle:       db.Collection(collectionArticles),
		collectionAutoincrement: db.Collection(collectionAutoincrement),
	}
}

func (s *Storage) NextArticleID(ctx context.Context) (int, error) {
	inc, err := s.autoincrement(ctx, collectionArticles)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return inc.Counter, nil
}

func (s *Storage) CreateArticle(ctx context.Context, id int, content string, duplicateIDs []int) error {
	art := article{
		ID:           id,
		Content:      content,
		DuplicateIDs: duplicateIDs,
		CreatedAt:    time.Now(),
	}

	ma, err := bson.Marshal(&art)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	if _, err := s.collectionArticle.InsertOne(ctx, ma); err != nil {
		return errors.Wrap(err, "failed to insert")
	}

	return nil
}

func (s *Storage) ArticleByID(ctx context.Context, id int) (model.Article, error) {
	res := s.collectionArticle.FindOne(ctx, bson.D{{Key: "id", Value: id}})
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

	return toModelArticle(art), nil
}

func (s *Storage) AllArticles(ctx context.Context) ([]model.Article, error) {
	articles := make([]model.Article, 0, maxArticles)

	cur, err := s.collectionArticle.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all documents")
	}

	for cur.TryNext(ctx) && len(articles) != maxArticles {
		art := article{}
		if err := cur.Decode(&art); err != nil {
			return nil, errors.Wrap(err, "failed to cursor decode")
		}

		articles = append(articles, toModelArticle(art))
	}

	return articles, nil
}

func toModelArticle(art article) model.Article {
	return model.Article{
		ID:           art.ID,
		Content:      art.Content,
		DuplicateIDs: art.DuplicateIDs,
	}
}

func (s *Storage) autoincrement(ctx context.Context, collection string) (*autoincrement, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	doc := &autoincrement{}
	filter := bson.M{"collection": collection}
	update := bson.M{
		"$inc": bson.M{"counter": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	if err := s.collectionAutoincrement.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "failed to find one and update")
	}

	return doc, nil
}
