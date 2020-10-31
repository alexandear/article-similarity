package handler

import (
	"context"
	"fmt"

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
	mongo *mongo.Client
	sim   *similarity.Similarity
}

func New(mongo *mongo.Client, sim *similarity.Similarity) *Handler {
	return &Handler{
		mongo: mongo,
		sim:   sim,
	}
}

func (h *Handler) ConfigureHandlers(api *operations.ArticleSimilarityAPI) {
	api.PostArticlesHandler = operations.PostArticlesHandlerFunc(h.PostArticles)
	api.GetArticlesIDHandler = operations.GetArticlesIDHandlerFunc(h.GetArticleByID)
}

func (h *Handler) PostArticles(params operations.PostArticlesParams) middleware.Responder {
	content := *params.Body.Content
	if content == "" {
		return operations.NewPostArticlesBadRequest().WithPayload(&models.Error{
			Message: swag.String("empty content"),
			Code:    0,
		})
	}

	ctx := params.HTTPRequest.Context()

	autoincrement, err := h.autoincrement(ctx, "articles")
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

	if _, err := h.mongo.Database("dev").Collection("articles").InsertOne(ctx, ma); err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	modelsArticle := h.modelsArticle(ctx, article)

	return operations.NewPostArticlesCreated().WithPayload(modelsArticle)
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	res := h.mongo.Database("dev").Collection("articles").
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
	collection := h.mongo.Database("dev").Collection("articles")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get all documents"))

		return nil
	}

	duplicateIDs := make([]int64, 0, maxDuplicateIDs)

	for cursor.TryNext(ctx) {
		a := &Article{}
		if err := cursor.Decode(a); err != nil {
			fmt.Println(errors.Wrap(err, "failed to cursor decode"))

			continue
		}

		if a.ID == id {
			continue
		}

		if h.sim.IsSimilar(content, a.Content) {
			duplicateIDs = append(duplicateIDs, int64(a.ID))
		}
	}

	return duplicateIDs
}
