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
	CreateArticle(ctx context.Context, id int, content string, duplicateIDs []int, isUnique bool) error
	UpdateArticle(ctx context.Context, id int, duplicateIDs []int) error
	ArticleByID(ctx context.Context, id int) (model.Article, error)
	AllArticles(ctx context.Context) ([]model.Article, error)
	UniqueArticles(ctx context.Context) ([]model.Article, error)
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

	duplicateIDs, err := a.duplicateArticleIDs(ctx, id, content)
	if err != nil {
		a.logger("failed to find duplicate articles ids: %v", err)
	}

	isUnique := len(duplicateIDs) == 0
	if err := a.storage.CreateArticle(ctx, id, content, duplicateIDs, isUnique); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to create article")
	}

	if !isUnique {
		a.updateArticlesWithDuplicateID(ctx, duplicateIDs, id)
	}

	return model.Article{
		ID:           id,
		Content:      content,
		DuplicateIDs: duplicateIDs,
		IsUnique:     isUnique,
	}, nil
}

func (a *Article) updateArticlesWithDuplicateID(ctx context.Context, duplicateIDs []int, id int) {
	for _, did := range duplicateIDs {
		art, err := a.storage.ArticleByID(ctx, did)
		if err != nil {
			a.logger("failed to get article by id=%d: %v", did, err)

			continue
		}

		if art.IsUnique {
			continue
		}

		if err := a.storage.UpdateArticle(ctx, art.ID, append(art.DuplicateIDs, id)); err != nil {
			a.logger("failed to update article=%d: %v", art.ID, err)
		}
	}
}

func (a *Article) ArticleByID(ctx context.Context, id int) (model.Article, error) {
	article, err := a.storage.ArticleByID(ctx, id)
	if err != nil {
		return model.Article{}, errors.WithStack(err)
	}

	return article, nil
}

func (a *Article) UniqueArticles(ctx context.Context) ([]model.Article, error) {
	articles, err := a.storage.UniqueArticles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unique articles")
	}

	return articles, nil
}

func (a *Article) DuplicateGroups(ctx context.Context) ([]model.DuplicateGroup, error) {
	return []model.DuplicateGroup{
		{
			IDs: []int{1, 3, 5},
		},
		{
			IDs: []int{7, 8, 9, 10, 11},
		},
	}, nil
}

func (a *Article) duplicateArticleIDs(ctx context.Context, id int, content string) ([]int, error) {
	articles, err := a.storage.AllArticles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all articles")
	}

	duplicates := make([]int, 0, len(articles))

	for _, article := range articles {
		if a.sim.IsSimilar(id, content, article.ID, article.Content) {
			duplicates = append(duplicates, article.ID)
		}
	}

	return duplicates, nil
}
