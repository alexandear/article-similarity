package article

import (
	"context"

	"github.com/pkg/errors"

	"github.com/devchallenge/article-similarity/internal/model"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type Article struct {
	logger  func(format string, v ...interface{})
	sim     *similarity.Similarity
	storage Storage
}

type Storage interface {
	NextArticleID(ctx context.Context) (int, error)
	CreateArticle(ctx context.Context, id int, content string) error
	ArticleByID(ctx context.Context, id int) (model.Article, error)
	AllArticles(ctx context.Context) ([]model.Article, error)
}

func New(
	logger func(format string, v ...interface{}),
	sim *similarity.Similarity,
	storage Storage,
) *Article {
	return &Article{
		logger:  logger,
		sim:     sim,
		storage: storage,
	}
}

func (a *Article) CreateArticle(ctx context.Context, content string) (model.Article, error) {
	id, err := a.storage.NextArticleID(ctx)
	if err != nil {
		return model.Article{}, errors.Wrap(err, "failed to get next article id")
	}

	if err := a.storage.CreateArticle(ctx, id, content); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to create article")
	}

	return a.article(ctx, id, content), nil
}

func (a *Article) ArticleByID(ctx context.Context, id int) (model.Article, error) {
	article, err := a.storage.ArticleByID(ctx, id)
	if err != nil {
		return model.Article{}, errors.WithStack(err)
	}

	return a.article(ctx, id, article.Content), nil
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
	articles, err := a.storage.AllArticles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all articles")
	}

	duplicates := make([]int, 0, len(articles))

	for _, article := range articles {
		if id == article.ID {
			continue
		}

		if a.sim.IsSimilar(id, content, article.ID, article.Content) {
			duplicates = append(duplicates, article.ID)
		}
	}

	return duplicates, nil
}
