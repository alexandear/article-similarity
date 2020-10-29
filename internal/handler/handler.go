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
}

func New(store *memkv.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) ConfigureHandlers(api *operations.ArticleSimilarityAPIAPI) {
	api.PostArticlesHandler = operations.PostArticlesHandlerFunc(h.PostArticles)
}

func (h *Handler) PostArticles(params operations.PostArticlesParams) middleware.Responder {
	content := *params.Body.Content
	if content == "" {
		return operations.NewPostArticlesBadRequest().WithPayload(&models.Error{
			Message: swag.String("empty content"),
			Code:    0,
		})
	}

	idContents, err := h.store.GetAll(idKeyPrefix + "*")
	if err != nil {
		return operations.NewPostArticlesInternalServerError().WithPayload(&models.Error{
			Message: swag.String("failed to get contents"),
			Code:    0,
		})
	}

	duplicateIDs := make([]int64, 0, len(idContents))
	sim := similarity.NewSimilarity()

	for _, idContent := range idContents {
		if sim.IsSimilar(content, idContent.Value) {
			duplicateIDs = append(duplicateIDs, int64(keyToID(idContent.Key)))
		}
	}

	id := h.nextID()
	h.store.Set(idKey(id), content)

	return operations.NewPostArticlesOK().WithPayload(&operations.PostArticlesOKBody{
		Content:             swag.String(content),
		DuplicateArticleIds: duplicateIDs,
		ID:                  swag.Int64(int64(id)),
	})
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
