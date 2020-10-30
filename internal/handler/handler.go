package handler

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/kelseyhightower/memkv"
	"github.com/pkg/errors"

	"github.com/devchallenge/article-similarity/internal/models"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type Handler struct {
	store *memkv.Store
	sim   *similarity.Similarity
}

func New(store *memkv.Store, sim *similarity.Similarity) *Handler {
	return &Handler{
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

	id := h.nextID()
	h.store.Set(idKey(id), content)

	return operations.NewPostArticlesCreated().WithPayload(&models.Article{
		Content:             swag.String(content),
		DuplicateArticleIds: h.duplicateArticleIDs(id, content),
		ID:                  swag.Int64(int64(id)),
	})
}

func (h *Handler) GetArticleByID(params operations.GetArticlesIDParams) middleware.Responder {
	content, err := h.store.GetValue(idKey(int(params.ID)))

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
	idContents, err := h.store.GetAll(idKeyPrefix + "*")
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get all contents"))

		return nil
	}

	duplicateIDs := make([]int64, 0, len(idContents))

	for _, idContent := range idContents {
		articleID := keyToID(idContent.Key)
		if articleID == id {
			continue
		}

		if h.sim.IsSimilar(content, idContent.Value) {
			duplicateIDs = append(duplicateIDs, int64(articleID))
		}
	}

	return duplicateIDs
}

const (
	nextIDKey   = "next_id"
	idKeyPrefix = "id_"
)

func idKey(id int) string {
	return idKeyPrefix + strconv.Itoa(id)
}

func keyToID(idStr string) int {
	id, err := strconv.Atoi(idStr[len(idKeyPrefix):])
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to convert key to id"))

		return 0
	}

	return id
}

func (h *Handler) nextID() int {
	id := 0

	defer func() {
		h.store.Set(nextIDKey, strconv.Itoa(id+1))
	}()

	idStr, err := h.store.GetValue(nextIDKey)

	switch {
	case err == nil:
	case errors.As(err, &memkv.ErrNotExist):
		return id
	default:
		fmt.Println(errors.Wrap(err, "failed to get next id value"))

		return id
	}

	id, err = strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to convert string to id"))
	}

	return id
}
