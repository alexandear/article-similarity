package handler

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devchallenge/article-similarity/internal/models"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type Handler struct {
	logger func(format string, v ...interface{})

	mongo         *mongo.Client
	mongoDatabase string

	sim *similarity.Similarity
}

func New(
	logger func(format string, v ...interface{}),
	mongo *mongo.Client, mongoDatabase string,
	sim *similarity.Similarity,
) *Handler {
	return &Handler{
		logger: logger,

		mongo:         mongo,
		mongoDatabase: mongoDatabase,

		sim: sim,
	}
}

func (h *Handler) ConfigureHandlers(api *operations.ArticleSimilarityAPI) {
	api.PostArticlesHandler = operations.PostArticlesHandlerFunc(h.PostArticles)
	api.GetArticlesIDHandler = operations.GetArticlesIDHandlerFunc(h.GetArticleByID)
}

const (
	collectionArticles = "articles"
)

func (h *Handler) PostArticles(params operations.PostArticlesParams) middleware.Responder {
	content := *params.Body.Content
	if content == "" {
		return operations.NewPostArticlesBadRequest().WithPayload(&models.Error{
			Message: swag.String("empty content"),
			Code:    0,
		})
	}

	ctx := params.HTTPRequest.Context()

	autoincrement, err := h.autoincrement(ctx, collectionArticles)
	if err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	id := autoincrement.Counter
	article := Article{
		ID:      id,
		Content: content,
	}

	ma, err := bson.Marshal(&article)
	if err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	if _, err := h.mongo.Database(h.mongoDatabase).Collection(collectionArticles).InsertOne(ctx, ma); err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	modelsArticle := h.modelsArticle(ctx, article)

	return operations.NewPostArticlesCreated().WithPayload(modelsArticle)
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	res := h.mongo.Database(h.mongoDatabase).Collection(collectionArticles).
		FindOne(ctx, bson.D{{Key: "id", Value: params.ID}})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return operations.NewGetArticlesIDNotFound()
	}

	if res.Err() != nil {
		return operations.NewGetArticlesIDInternalServerError()
	}

	article := Article{}
	if err := res.Decode(&article); err != nil {
		return operations.NewGetArticlesIDInternalServerError()
	}

	modelsArticle := h.modelsArticle(ctx, article)

	return operations.NewGetArticlesIDOK().WithPayload(modelsArticle)
}

func (h *Handler) modelsArticle(ctx context.Context, article Article) *models.Article {
	duplicateIDs := h.duplicateArticleIDs(ctx, article.ID, article.Content)

	return &models.Article{
		ID:                  swag.Int64(int64(article.ID)),
		Content:             swag.String(article.Content),
		DuplicateArticleIds: duplicateIDs,
	}
}

const maxDuplicateIDs = 100

func (h *Handler) duplicateArticleIDs(ctx context.Context, id int, content string) []int64 {
	collection := h.mongo.Database(h.mongoDatabase).Collection(collectionArticles)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		h.logger("failed to get all documents: %v", err)

		return nil
	}

	duplicateIDs := make([]int64, 0, maxDuplicateIDs)

	for cursor.TryNext(ctx) {
		a := &Article{}
		if err := cursor.Decode(a); err != nil {
			h.logger("failed to cursor decode: %v", err)

			continue
		}

		if a.ID == id {
			continue
		}

		if h.sim.IsSimilar(id, content, a.ID, a.Content) {
			duplicateIDs = append(duplicateIDs, int64(a.ID))
		}
	}

	return duplicateIDs
}
