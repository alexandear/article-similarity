package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	articlesim "github.com/devchallenge/article-similarity/internal"
)

const (
	maxArticles        = 1000
	maxDuplicateGroups = 1000

	collectionArticles        = "articles"
	collectionDuplicateGroups = "duplicate_groups"
	collectionAutoincrement   = "autoincrement"
)

type article struct {
	ID               articlesim.ArticleID        `bson:"id"`
	Content          string                      `bson:"content"`
	DuplicateIDs     []articlesim.ArticleID      `bson:"duplicate_ids"`
	IsUnique         bool                        `bson:"is_unique"`
	DuplicateGroupID articlesim.DuplicateGroupID `bson:"duplicate_group_id"`
}

type duplicateGroup struct {
	ID        articlesim.DuplicateGroupID `bson:"id"`
	ArticleID articlesim.ArticleID        `bson:"article_id"`
}

type autoincrement struct {
	ID         primitive.ObjectID `bson:"_id"`
	Collection string             `bson:"collection"`
	Counter    int                `bson:"counter"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

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

func (s *Storage) NextArticleID(ctx context.Context) (articlesim.ArticleID, error) {
	inc, err := s.autoincrement(ctx, collectionArticles)
	if err != nil {
		return 0, fmt.Errorf("failed to get autoicrement for articles: %w", err)
	}

	return articlesim.ArticleID(inc.Counter), nil
}

func (s *Storage) CreateArticle(ctx context.Context, id articlesim.ArticleID, content string,
	duplicateIDs []articlesim.ArticleID, isUnique bool, duplicateGroupID articlesim.DuplicateGroupID) error {
	art := article{
		ID:               id,
		Content:          content,
		DuplicateIDs:     duplicateIDs,
		IsUnique:         isUnique,
		DuplicateGroupID: duplicateGroupID,
	}

	ma, err := bson.Marshal(&art)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	if _, err := s.collectionArticle.InsertOne(ctx, ma); err != nil {
		return fmt.Errorf("failed to insert article: %w", err)
	}

	return nil
}

func (s *Storage) UpdateArticle(ctx context.Context, id articlesim.ArticleID,
	duplicateIDs []articlesim.ArticleID) error {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.M{
		"$set": bson.M{"duplicate_ids": duplicateIDs},
	}

	if err := s.collectionArticle.FindOneAndUpdate(ctx, filter, update, nil).Err(); err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	return nil
}

func (s *Storage) ArticleByID(ctx context.Context, id articlesim.ArticleID) (articlesim.Article, error) {
	res := s.collectionArticle.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return articlesim.Article{}, fmt.Errorf("not found: %w", articlesim.ErrArticleNotFound)
	}

	if res.Err() != nil {
		return articlesim.Article{}, fmt.Errorf("failed to find: %w", res.Err())
	}

	art := article{}
	if err := res.Decode(&art); err != nil {
		return articlesim.Article{}, fmt.Errorf("failed to decode: %w", err)
	}

	return toModelArticle(art), nil
}

func (s *Storage) AllArticles(ctx context.Context) ([]articlesim.Article, error) {
	return s.articles(ctx, bson.D{})
}

func (s *Storage) UniqueArticles(ctx context.Context) ([]articlesim.Article, error) {
	return s.articles(ctx, bson.D{{Key: "is_unique", Value: true}})
}

func (s *Storage) articles(ctx context.Context, filter bson.D) ([]articlesim.Article, error) {
	articles := make([]articlesim.Article, 0, maxArticles)

	cur, err := s.collectionArticle.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles: %w", err)
	}

	for cur.TryNext(ctx) && len(articles) != maxArticles {
		art := article{}
		if err := cur.Decode(&art); err != nil {
			return nil, fmt.Errorf("failed to cursor decode to article: %w", err)
		}

		articles = append(articles, toModelArticle(art))
	}

	return articles, nil
}

func (s *Storage) NextDuplicateGroupID(ctx context.Context) (articlesim.DuplicateGroupID, error) {
	inc, err := s.autoincrement(ctx, collectionDuplicateGroups)
	if err != nil {
		return 0, fmt.Errorf("failed to get autoicrement for duplicate groups: %w", err)
	}

	return articlesim.DuplicateGroupID(inc.Counter), nil
}

func (s *Storage) CreateDuplicateGroup(ctx context.Context, id articlesim.DuplicateGroupID,
	articleID articlesim.ArticleID) error {
	dg := duplicateGroup{
		ID:        id,
		ArticleID: articleID,
	}

	mdg, err := bson.Marshal(&dg)
	if err != nil {
		return fmt.Errorf("failed to marshal duplicate group: %w", err)
	}

	if _, err := s.collectionDuplicateGroup.InsertOne(ctx, mdg); err != nil {
		return fmt.Errorf("failed to insert duplicate group: %w", err)
	}

	return nil
}

func (s *Storage) AllDuplicateGroups(ctx context.Context) ([]articlesim.DuplicateGroup, error) {
	groups := make([]articlesim.DuplicateGroup, 0, maxDuplicateGroups)

	cur, err := s.collectionDuplicateGroup.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find duplicate groups: %w", err)
	}

	for cur.TryNext(ctx) && len(groups) != maxDuplicateGroups {
		group := duplicateGroup{}
		if err := cur.Decode(&group); err != nil {
			return nil, fmt.Errorf("failed to cursor decode to group: %w", err)
		}

		groups = append(groups, articlesim.DuplicateGroup{
			DuplicateGroupID: group.ID,
			ArticleID:        group.ArticleID,
		})
	}

	return groups, nil
}

func toModelArticle(art article) articlesim.Article {
	return articlesim.Article{
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
		return nil, fmt.Errorf("failed to find one and update: %w", err)
	}

	return doc, nil
}
