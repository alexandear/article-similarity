package handler

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/devchallenge/article-similarity/internal/model"
	"github.com/devchallenge/article-similarity/internal/swagger/models"
	"github.com/devchallenge/article-similarity/internal/swagger/restapi/operations"
)

type ArticleServer interface {
	CreateArticle(ctx context.Context, content string) (model.Article, error)
	ArticleByID(ctx context.Context, id int) (model.Article, error)
}

type Handler struct {
	article ArticleServer
}

func New(article ArticleServer) *Handler {
	return &Handler{
		article: article,
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

	article, err := h.article.CreateArticle(params.HTTPRequest.Context(), content)
	if err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	return operations.NewPostArticlesCreated().WithPayload(h.modelsArticle(article))
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	article, err := h.article.ArticleByID(params.HTTPRequest.Context(), int(params.ID))

	if errors.Is(err, mongo.ErrNoDocuments) {
		return operations.NewGetArticlesIDNotFound()
	}

	if err != nil {
		return operations.NewGetArticlesIDInternalServerError()
	}

	return operations.NewGetArticlesIDOK().WithPayload(h.modelsArticle(article))
}

func (h *Handler) modelsArticle(article model.Article) *models.Article {
	const maxDuplicates = 100

	duplicateIDs := make([]int64, 0, maxDuplicates)
	for _, id := range article.DuplicateIDs {
		duplicateIDs = append(duplicateIDs, int64(id))
	}

	return &models.Article{
		ID:                  swag.Int64(int64(article.ID)),
		Content:             swag.String(article.Content),
		DuplicateArticleIds: duplicateIDs,
	}
}
