package http

import (
	"context"
	"errors"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	articlesim "github.com/devchallenge/article-similarity/internal"
	"github.com/devchallenge/article-similarity/internal/http/models"
	"github.com/devchallenge/article-similarity/internal/http/restapi/operations"
)

const (
	serverTimeout = 5 * time.Second
)

type ArticleServer interface {
	CreateArticle(ctx context.Context, content string) (articlesim.Article, error)
	ArticleByID(ctx context.Context, id articlesim.ArticleID) (articlesim.Article, error)
	UniqueArticles(ctx context.Context) ([]articlesim.Article, error)
	DuplicateGroups(ctx context.Context) ([]articlesim.DuplicateGroupResp, error)
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

	article, err := h.article.ArticleByID(ctx, articlesim.ArticleID(params.ID))

	if errors.Is(err, articlesim.ErrArticleNotFound) {
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

	modelsDuplicateGroups := make([][]models.ArticleID, 0, len(groups))

	for _, g := range groups {
		ids := make([]models.ArticleID, 0, len(g.ArticleIDs))

		for _, id := range g.ArticleIDs {
			ids = append(ids, models.ArticleID(id))
		}

		modelsDuplicateGroups = append(modelsDuplicateGroups, ids)
	}

	return operations.NewGetDuplicateGroupsOK().WithPayload(&operations.GetDuplicateGroupsOKBody{
		DuplicateGroups: modelsDuplicateGroups,
	})
}

func modelsArticle(article articlesim.Article) *models.Article {
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
