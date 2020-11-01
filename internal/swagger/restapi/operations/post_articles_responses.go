// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/devchallenge/article-similarity/internal/swagger/models"
)

// PostArticlesCreatedCode is the HTTP code returned for type PostArticlesCreated
const PostArticlesCreatedCode int = 201

/*PostArticlesCreated Article added.

swagger:response postArticlesCreated
*/
type PostArticlesCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Article `json:"body,omitempty"`
}

// NewPostArticlesCreated creates PostArticlesCreated with default headers values
func NewPostArticlesCreated() *PostArticlesCreated {

	return &PostArticlesCreated{}
}

// WithPayload adds the payload to the post articles created response
func (o *PostArticlesCreated) WithPayload(payload *models.Article) *PostArticlesCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post articles created response
func (o *PostArticlesCreated) SetPayload(payload *models.Article) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostArticlesCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostArticlesBadRequestCode is the HTTP code returned for type PostArticlesBadRequest
const PostArticlesBadRequestCode int = 400

/*PostArticlesBadRequest Invalid arguments

swagger:response postArticlesBadRequest
*/
type PostArticlesBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostArticlesBadRequest creates PostArticlesBadRequest with default headers values
func NewPostArticlesBadRequest() *PostArticlesBadRequest {

	return &PostArticlesBadRequest{}
}

// WithPayload adds the payload to the post articles bad request response
func (o *PostArticlesBadRequest) WithPayload(payload *models.Error) *PostArticlesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post articles bad request response
func (o *PostArticlesBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostArticlesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostArticlesInternalServerErrorCode is the HTTP code returned for type PostArticlesInternalServerError
const PostArticlesInternalServerErrorCode int = 500

/*PostArticlesInternalServerError Internal server error

swagger:response postArticlesInternalServerError
*/
type PostArticlesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostArticlesInternalServerError creates PostArticlesInternalServerError with default headers values
func NewPostArticlesInternalServerError() *PostArticlesInternalServerError {

	return &PostArticlesInternalServerError{}
}

// WithPayload adds the payload to the post articles internal server error response
func (o *PostArticlesInternalServerError) WithPayload(payload *models.Error) *PostArticlesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post articles internal server error response
func (o *PostArticlesInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostArticlesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
