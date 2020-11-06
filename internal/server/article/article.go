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
	NextArticleID(ctx context.Context) (model.ArticleID, error)
	CreateArticle(ctx context.Context, id model.ArticleID, content string, duplicateIDs []model.ArticleID, isUnique bool,
		duplicateGroupID model.DuplicateGroupID) error
	UpdateArticle(ctx context.Context, id model.ArticleID, duplicateIDs []model.ArticleID) error
	ArticleByID(ctx context.Context, id model.ArticleID) (model.Article, error)
	AllArticles(ctx context.Context) ([]model.Article, error)
	UniqueArticles(ctx context.Context) ([]model.Article, error)
	NextDuplicateGroupID(ctx context.Context) (model.DuplicateGroupID, error)
	CreateDuplicateGroup(ctx context.Context, duplicateGroupID model.DuplicateGroupID, articleID model.ArticleID) error
	AllDuplicateGroups(ctx context.Context) ([]model.DuplicateGroup, error)
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

	duplicateIDs, duplicateGroupID, err := a.duplicateArticleIDsWithDuplicateGroupID(ctx, id, content)
	if err != nil {
		a.logger("failed to find duplicate articles ids: %v", err)
	}

	isUnique := len(duplicateIDs) == 0
	if err := a.storage.CreateArticle(ctx, id, content, duplicateIDs, isUnique, duplicateGroupID); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to create article")
	}

	if err := a.storage.CreateDuplicateGroup(ctx, duplicateGroupID, id); err != nil {
		return model.Article{}, errors.Wrap(err, "failed to create duplicate group")
	}

	if !isUnique {
		a.updateArticlesWithDuplicateID(ctx, duplicateIDs, id)
	}

	return model.Article{
		ID:               id,
		Content:          content,
		DuplicateIDs:     duplicateIDs,
		IsUnique:         isUnique,
		DuplicateGroupID: duplicateGroupID,
	}, nil
}

func (a *Article) updateArticlesWithDuplicateID(ctx context.Context, duplicateIDs []model.ArticleID,
	id model.ArticleID) {
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

func (a *Article) ArticleByID(ctx context.Context, id model.ArticleID) (model.Article, error) {
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

func (a *Article) DuplicateGroups(ctx context.Context) (map[model.DuplicateGroupID][]model.ArticleID, error) {
	groups, err := a.storage.AllDuplicateGroups(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all duplicate groups")
	}

	duplicateGroups := make(map[model.DuplicateGroupID][]model.ArticleID, len(groups))
	for _, g := range groups {
		duplicateGroups[g.DuplicateGroupID] = append(duplicateGroups[g.DuplicateGroupID], g.ArticleID)
	}

	return duplicateGroups, nil
}

func (a *Article) duplicateArticleIDsWithDuplicateGroupID(ctx context.Context, id model.ArticleID, content string,
) ([]model.ArticleID, model.DuplicateGroupID, error) {
	articles, err := a.storage.AllArticles(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get all articles")
	}

	duplicates := make([]model.ArticleID, 0, len(articles))

	var duplicateGroupID model.DuplicateGroupID

	for _, article := range articles {
		if a.sim.IsSimilar(int(id), content, int(article.ID), article.Content) {
			duplicates = append(duplicates, article.ID)
			duplicateGroupID = article.DuplicateGroupID
		}
	}

	if duplicateGroupID == 0 {
		gid, err := a.storage.NextDuplicateGroupID(ctx)
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to get next duplicate group id")
		}

		return nil, gid, nil
	}

	return duplicates, duplicateGroupID, nil
}
