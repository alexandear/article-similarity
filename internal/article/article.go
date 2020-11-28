package article

import (
	"context"
	"fmt"
	"log"

	articlesim "github.com/devchallenge/article-similarity/internal"
)

type Similarity interface {
	IsSimilar(idA int, contentA string, idB int, contentB string) bool
	Similarity(idA int, contentA string, idB int, contentB string) float64
}

type Storage interface {
	NextArticleID(ctx context.Context) (articlesim.ArticleID, error)
	CreateArticle(ctx context.Context, id articlesim.ArticleID, content string, duplicateIDs []articlesim.ArticleID,
		isUnique bool, duplicateGroupID articlesim.DuplicateGroupID) error
	UpdateArticle(ctx context.Context, id articlesim.ArticleID, duplicateIDs []articlesim.ArticleID) error
	ArticleByID(ctx context.Context, id articlesim.ArticleID) (articlesim.Article, error)
	AllArticles(ctx context.Context) ([]articlesim.Article, error)
	UniqueArticles(ctx context.Context) ([]articlesim.Article, error)
	NextDuplicateGroupID(ctx context.Context) (articlesim.DuplicateGroupID, error)
	CreateDuplicateGroup(ctx context.Context, duplicateGroupID articlesim.DuplicateGroupID,
		articleID articlesim.ArticleID) error
	AllDuplicateGroups(ctx context.Context) ([]articlesim.DuplicateGroup, error)
}

type Service struct {
	similar Similarity
	storage Storage
}

func New(similar Similarity, storage Storage) *Service {
	return &Service{
		similar: similar,
		storage: storage,
	}
}

func (a *Service) CreateArticle(ctx context.Context, content string) (articlesim.Article, error) {
	id, err := a.storage.NextArticleID(ctx)
	if err != nil {
		return articlesim.Article{}, fmt.Errorf("failed to get next article id: %w", err)
	}

	duplicateIDs, duplicateGroupID, err := a.duplicateArticleIDsWithDuplicateGroupID(ctx, id, content)
	if err != nil {
		log.Printf("failed to find duplicate articles ids: %v", err)
	}

	isUnique := len(duplicateIDs) == 0
	if err := a.storage.CreateArticle(ctx, id, content, duplicateIDs, isUnique, duplicateGroupID); err != nil {
		return articlesim.Article{}, fmt.Errorf("failed to create article: %w", err)
	}

	if err := a.storage.CreateDuplicateGroup(ctx, duplicateGroupID, id); err != nil {
		return articlesim.Article{}, fmt.Errorf("failed to create duplicate group: %w", err)
	}

	if !isUnique {
		a.updateArticlesWithDuplicateID(ctx, duplicateIDs, id)
	}

	return articlesim.Article{
		ID:               id,
		Content:          content,
		DuplicateIDs:     duplicateIDs,
		IsUnique:         isUnique,
		DuplicateGroupID: duplicateGroupID,
	}, nil
}

func (a *Service) updateArticlesWithDuplicateID(ctx context.Context, duplicateIDs []articlesim.ArticleID,
	id articlesim.ArticleID) {
	for _, did := range duplicateIDs {
		art, err := a.storage.ArticleByID(ctx, did)
		if err != nil {
			log.Printf("failed to get article by id=%d: %v", did, err)

			continue
		}

		if art.IsUnique {
			continue
		}

		if err := a.storage.UpdateArticle(ctx, art.ID, append(art.DuplicateIDs, id)); err != nil {
			log.Printf("failed to update article=%d: %v", art.ID, err)
		}
	}
}

func (a *Service) ArticleByID(ctx context.Context, id articlesim.ArticleID) (articlesim.Article, error) {
	article, err := a.storage.ArticleByID(ctx, id)
	if err != nil {
		return articlesim.Article{}, fmt.Errorf("failed to get article from storage: %w", err)
	}

	return article, nil
}

func (a *Service) UniqueArticles(ctx context.Context) ([]articlesim.Article, error) {
	articles, err := a.storage.UniqueArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique articles: %w", err)
	}

	return articles, nil
}

func (a *Service) DuplicateGroups(ctx context.Context) (map[articlesim.DuplicateGroupID][]articlesim.ArticleID, error) {
	groups, err := a.storage.AllDuplicateGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all duplicate groups: %w", err)
	}

	duplicateGroups := make(map[articlesim.DuplicateGroupID][]articlesim.ArticleID, len(groups))
	for _, g := range groups {
		duplicateGroups[g.DuplicateGroupID] = append(duplicateGroups[g.DuplicateGroupID], g.ArticleID)
	}

	return duplicateGroups, nil
}

func (a *Service) duplicateArticleIDsWithDuplicateGroupID(ctx context.Context, id articlesim.ArticleID, content string,
) ([]articlesim.ArticleID, articlesim.DuplicateGroupID, error) {
	articles, err := a.storage.AllArticles(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all articles: %w", err)
	}

	duplicates := make([]articlesim.ArticleID, 0, len(articles))

	var duplicateGroupID articlesim.DuplicateGroupID

	for _, article := range articles {
		if a.similar.IsSimilar(int(id), content, int(article.ID), article.Content) {
			duplicates = append(duplicates, article.ID)
			duplicateGroupID = article.DuplicateGroupID
		}
	}

	if duplicateGroupID == 0 {
		gid, err := a.storage.NextDuplicateGroupID(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get next duplicate group id: %w", err)
		}

		return nil, gid, nil
	}

	return duplicates, duplicateGroupID, nil
}
