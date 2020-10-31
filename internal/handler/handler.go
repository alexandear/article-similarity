package handler

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/kelseyhightower/memkv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devchallenge/article-similarity/internal/models"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type Handler struct {
	mongo *mongo.Client
	store *memkv.Store
	sim   *similarity.Similarity
}

func New(mongo *mongo.Client, store *memkv.Store, sim *similarity.Similarity) *Handler {
	return &Handler{
		mongo: mongo,
		store: store,
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

	autoincrement, err := h.autoincrement(context.TODO(), "articles")
	if err != nil {
		panic(err)
	}

	id := autoincrement.Counter
	article := &Article{
		ID:      id,
		Content: content,
	}

	ma, err := bson.Marshal(article)
	if err != nil {
		panic(err)
	}

	if _, err := h.mongo.Database("dev").Collection("articles").InsertOne(context.TODO(), ma); err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	return operations.NewPostArticlesCreated().WithPayload(&models.Article{
		Content:             swag.String(content),
		DuplicateArticleIds: h.duplicateArticleIDs(id, content),
		ID:                  swag.Int64(int64(id)),
	})
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	content, err := h.store.GetValue("1")

	switch {
	case err == nil:
	case errors.As(err, &memkv.ErrNotExist):
		return operations.NewGetArticlesIDNotFound()
	default:
		fmt.Println(errors.Wrap(err, "failed to get content"))

		return operations.NewGetArticlesIDInternalServerError().WithPayload(&models.Error{
			Code:    0,
			Message: swag.String("failed to get article"),
		})
	}

	return operations.NewGetArticlesIDOK().WithPayload(&models.Article{
		ID:                  swag.Int64(params.ID),
		Content:             swag.String(content),
		DuplicateArticleIds: h.duplicateArticleIDs(int(params.ID), content),
	})
}

func (h *Handler) duplicateArticleIDs(id int, content string) []int64 {
	collection := h.mongo.Database("dev").Collection("articles")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get all documents"))

		return nil
	}

	duplicateIDs := make([]int64, 0, 1)

	for cursor.TryNext(context.TODO()) {
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
