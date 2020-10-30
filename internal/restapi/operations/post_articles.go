// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostArticlesHandlerFunc turns a function with the right signature into a post articles handler
type PostArticlesHandlerFunc func(PostArticlesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostArticlesHandlerFunc) Handle(params PostArticlesParams) middleware.Responder {
	return fn(params)
}

// PostArticlesHandler interface for that can handle valid post articles params
type PostArticlesHandler interface {
	Handle(PostArticlesParams) middleware.Responder
}

// NewPostArticles creates a new http.Handler for the post articles operation
func NewPostArticles(ctx *middleware.Context, handler PostArticlesHandler) *PostArticles {
	return &PostArticles{Context: ctx, Handler: handler}
}

/*PostArticles swagger:route POST /articles postArticles

Add an article.

*/
type PostArticles struct {
	Context *middleware.Context
	Handler PostArticlesHandler
}

func (o *PostArticles) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostArticlesParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostArticlesBody post articles body
//
// swagger:model PostArticlesBody
type PostArticlesBody struct {

	// Article content
	// Required: true
	Content *string `json:"content"`
}

// Validate validates this post articles body
func (o *PostArticlesBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateContent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostArticlesBody) validateContent(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"content", "body", o.Content); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostArticlesBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostArticlesBody) UnmarshalBinary(b []byte) error {
	var res PostArticlesBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
