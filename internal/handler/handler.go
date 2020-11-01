package handler

import (
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"

	internalErrors "github.com/devchallenge/article-similarity/internal/errors"
	"github.com/devchallenge/article-similarity/internal/model"
	"github.com/devchallenge/article-similarity/internal/swagger/models"
	"github.com/devchallenge/article-similarity/internal/swagger/restapi/operations"
)

const (
	serverTimeout = 5 * time.Second
)

type ArticleServer interface {
	CreateArticle(ctx context.Context, content string) (model.Article, error)
	ArticleByID(ctx context.Context, id int) (model.Article, error)
	UniqueArticles(ctx context.Context) ([]model.Article, error)
	DuplicateGroups(ctx context.Context) ([]model.DuplicateGroup, error)
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
	api.GetArticlesHandler = operations.GetArticlesHandlerFunc(h.GetUniqueArticles)
	api.GetDuplicateGroupsHandler = operations.GetDuplicateGroupsHandlerFunc(h.GetDuplicateGroups)
}

func (h *Handler) PostArticles(params operations.PostArticlesParams) middleware.Responder {
	content := *params.Body.Content
	if content == "" {
		return operations.NewPostArticlesBadRequest().WithPayload(&models.Error{
			Message: swag.String("empty content"),
			Code:    0,
		})
	}

	ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), serverTimeout)
	defer cancel()

	article, err := h.article.CreateArticle(ctx, content)
	if err != nil {
		return operations.NewPostArticlesInternalServerError()
	}

	return operations.NewPostArticlesCreated().WithPayload(modelsArticle(article))
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), serverTimeout)
	defer cancel()

	article, err := h.article.ArticleByID(ctx, int(params.ID))

	if errors.Is(err, internalErrors.ErrNotFound) {
		return operations.NewGetArticlesIDNotFound()
	}

	if err != nil {
		return operations.NewGetArticlesIDInternalServerError()
	}

	return operations.NewGetArticlesIDOK().WithPayload(modelsArticle(article))
}

func (h *Handler) GetUniqueArticles(params operations.GetArticlesParams) middleware.Responder {
	ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), serverTimeout)
	defer cancel()

	articles, err := h.article.UniqueArticles(ctx)
	if err != nil {
		return operations.NewGetArticlesInternalServerError()
	}

	modelsArticles := make([]*models.Article, 0, len(articles))
	for _, article := range articles {
		modelsArticles = append(modelsArticles, modelsArticle(article))
	}

	return operations.NewGetArticlesOK().WithPayload(&operations.GetArticlesOKBody{
		Articles: modelsArticles,
	})
}

func (h *Handler) GetDuplicateGroups(params operations.GetDuplicateGroupsParams) middleware.Responder {
	ctx, cancel := context.WithTimeout(params.HTTPRequest.Context(), serverTimeout)
	defer cancel()

	groups, err := h.article.DuplicateGroups(ctx)
	if err != nil {
		return operations.NewGetDuplicateGroupsInternalServerError()
	}

	modelsGroups := make([][]models.ArticleID, 0, len(groups))

	for _, g := range groups {
		mg := make([]models.ArticleID, 0, len(g.IDs))

		for _, id := range g.IDs {
			mg = append(mg, models.ArticleID(id))
		}

		modelsGroups = append(modelsGroups, mg)
	}

	return operations.NewGetDuplicateGroupsOK().WithPayload(&operations.GetDuplicateGroupsOKBody{
		DuplicateGroups: modelsGroups,
	})
}

func modelsArticle(article model.Article) *models.Article {
	const maxDuplicates = 100

	duplicateIDs := make([]int64, 0, maxDuplicates)
	for _, id := range article.DuplicateIDs {
		duplicateIDs = append(duplicateIDs, int64(id))
	}

	return &models.Article{
		ID:                  models.ArticleID(int64(article.ID)),
		Content:             swag.String(article.Content),
		DuplicateArticleIds: duplicateIDs,
	}
}
