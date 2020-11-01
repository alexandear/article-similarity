package article

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devchallenge/article-similarity/internal/model"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type Article struct {
	logger func(format string, v ...interface{})

	sim *similarity.Similarity

	mongo         *mongo.Client
	mongoDatabase string
}

func New(
	logger func(format string, v ...interface{}),
	sim *similarity.Similarity,
	mongo *mongo.Client, mongoDatabase string,
) *Article {
	return &Article{
		logger:        logger,
		sim:           sim,
		mongo:         mongo,
		mongoDatabase: mongoDatabase,
	}
}

func (a *Article) CreateArticle(ctx context.Context, content string) (model.Article, error) {
	autoincrement, err := a.autoincrement(ctx, collectionArticles)
	if err != nil {
		return model.Article{}, nil
	}

	id := autoincrement.Counter
	article := article{
		ID:        id,
		Content:   content,
		CreatedAt: time.Now(),
	}

	ma, err := bson.Marshal(&article)
	if err != nil {
		return model.Article{}, errors.Wrap(err, "failed to marshal")
	}

	if _, err := a.mongo.Database(a.mongoDatabase).Collection(collectionArticles).InsertOne(ctx, ma); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to insert")
	}

	return a.article(ctx, id, content), nil
}

func (a *Article) ArticleByID(ctx context.Context, id int) (model.Article, error) {
	res := a.mongo.Database(a.mongoDatabase).Collection(collectionArticles).
		FindOne(ctx, bson.D{{Key: "id", Value: id}})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return model.Article{}, errors.Wrap(mongo.ErrNoDocuments, "not found")
	}

	if res.Err() != nil {
		return model.Article{}, errors.Wrap(res.Err(), "failed to find")
	}

	art := article{}
	if err := res.Decode(&art); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to decode")
	}

	return a.article(ctx, id, art.Content), nil
}

func (a *Article) article(ctx context.Context, id int, content string) model.Article {
	duplicateIDs, err := a.duplicateArticleIDs(ctx, id, content)
	if err != nil {
		a.logger("failed to find duplicate articles ids: %v", err)
	}

	return model.Article{
		ID:           id,
		Content:      content,
		DuplicateIDs: duplicateIDs,
	}
}

func (a *Article) duplicateArticleIDs(ctx context.Context, id int, content string) ([]int, error) {
	collection := a.mongo.Database(a.mongoDatabase).Collection(collectionArticles)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all documents")
	}

	var duplicateIDs []int

	for cursor.TryNext(ctx) {
		art := &article{}

		if err := cursor.Decode(art); err != nil {
			return nil, errors.Wrap(err, "failed to cursor decode")
		}

		if art.ID == id {
			continue
		}

		if a.sim.IsSimilar(id, content, art.ID, art.Content) {
			duplicateIDs = append(duplicateIDs, art.ID)
		}
	}

	return duplicateIDs, nil
}
