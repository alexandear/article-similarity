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
	maxArticles        = 1000
	maxDuplicateGroups = 1000

	collectionArticles        = "articles"
	collectionDuplicateGroups = "duplicate_groups"
	collectionAutoincrement   = "autoincrement"
)

type Storage struct {
	collectionArticle        *mongo.Collection
	collectionDuplicateGroup *mongo.Collection
	collectionAutoincrement  *mongo.Collection
}

func New(mc *mongo.Client, database string) *Storage {
	db := mc.Database(database)

	return &Storage{
		collectionArticle:        db.Collection(collectionArticles),
		collectionDuplicateGroup: db.Collection(collectionDuplicateGroups),
		collectionAutoincrement:  db.Collection(collectionAutoincrement),
	}
}

func (s *Storage) NextArticleID(ctx context.Context) (model.ArticleID, error) {
	inc, err := s.autoincrement(ctx, collectionArticles)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return model.ArticleID(inc.Counter), nil
}

func (s *Storage) CreateArticle(ctx context.Context, id model.ArticleID, content string, duplicateIDs []model.ArticleID,
	isUnique bool, duplicateGroupID model.DuplicateGroupID) error {
	art := article{
		ID:               id,
		Content:          content,
		DuplicateIDs:     duplicateIDs,
		IsUnique:         isUnique,
		DuplicateGroupID: duplicateGroupID,
	}

	ma, err := bson.Marshal(&art)
	if err != nil {
		return errors.Wrap(err, "failed to marshal article")
	}

	if _, err := s.collectionArticle.InsertOne(ctx, ma); err != nil {
		return errors.Wrap(err, "failed to insert article")
	}

	return nil
}

func (s *Storage) UpdateArticle(ctx context.Context, id model.ArticleID, duplicateIDs []model.ArticleID) error {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.M{
		"$set": bson.M{"duplicate_ids": duplicateIDs},
	}

	if err := s.collectionArticle.FindOneAndUpdate(ctx, filter, update, nil).Err(); err != nil {
		return errors.Wrap(err, "failed to update article")
	}

	return nil
}

func (s *Storage) ArticleByID(ctx context.Context, id model.ArticleID) (model.Article, error) {
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
	return s.articles(ctx, bson.D{})
}

func (s *Storage) UniqueArticles(ctx context.Context) ([]model.Article, error) {
	return s.articles(ctx, bson.D{{Key: "is_unique", Value: true}})
}

func (s *Storage) articles(ctx context.Context, filter bson.D) ([]model.Article, error) {
	articles := make([]model.Article, 0, maxArticles)

	cur, err := s.collectionArticle.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find articles")
	}

	for cur.TryNext(ctx) && len(articles) != maxArticles {
		art := article{}
		if err := cur.Decode(&art); err != nil {
			return nil, errors.Wrap(err, "failed to cursor decode to article")
		}

		articles = append(articles, toModelArticle(art))
	}

	return articles, nil
}

func (s *Storage) NextDuplicateGroupID(ctx context.Context) (model.DuplicateGroupID, error) {
	inc, err := s.autoincrement(ctx, collectionDuplicateGroups)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return model.DuplicateGroupID(inc.Counter), nil
}

func (s *Storage) CreateDuplicateGroup(ctx context.Context, id model.DuplicateGroupID, articleID model.ArticleID,
) error {
	dg := duplicateGroup{
		ID:        id,
		ArticleID: articleID,
	}

	mdg, err := bson.Marshal(&dg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal duplicate group")
	}

	if _, err := s.collectionDuplicateGroup.InsertOne(ctx, mdg); err != nil {
		return errors.Wrap(err, "failed to insert duplicate group")
	}

	return nil
}

func (s *Storage) AllDuplicateGroups(ctx context.Context) ([]model.DuplicateGroup, error) {
	groups := make([]model.DuplicateGroup, 0, maxDuplicateGroups)

	cur, err := s.collectionDuplicateGroup.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find duplicate groups")
	}

	for cur.TryNext(ctx) && len(groups) != maxDuplicateGroups {
		group := duplicateGroup{}
		if err := cur.Decode(&group); err != nil {
			return nil, errors.Wrap(err, "failed to cursor decode to group")
		}

		groups = append(groups, model.DuplicateGroup{
			DuplicateGroupID: group.ID,
			ArticleID:        group.ArticleID,
		})
	}

	return groups, nil
}

func toModelArticle(art article) model.Article {
	return model.Article{
		ID:               art.ID,
		Content:          art.Content,
		DuplicateIDs:     art.DuplicateIDs,
		IsUnique:         art.IsUnique,
		DuplicateGroupID: art.DuplicateGroupID,
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
