package restapi

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/errors"
)

// DefaultHTTPCode is used when the error Code cannot be used as an HTTP code.
const DefaultHTTPCode = http.StatusUnprocessableEntity

// ServeError the error handler interface implementation.
// It is a modified copy of (github.com/go-openapi/errors).ServeError and used to serve errors in go-swagger API.
// The reason for modification is to replace go-swagger default params validation errors.
// HTTP status and code (422 Unprocessable entity, 600+) with 400 BadArgument.
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")
	switch e := err.(type) {
	case *errors.CompositeError:
		er := flattenComposite(e)
		// strips composite errors to first element only
		if len(er.Errors) > 0 {
			ServeError(rw, r, er.Errors[0])
		} else {
			// guard against empty CompositeError (invalid construct)
			ServeError(rw, r, nil)
		}
	case *errors.MethodNotAllowedError:
		rw.Header().Add("Allow", strings.Join(err.(*errors.MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(e))
		}
	case errors.Error:
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write(errorAsJSON(errors.New(http.StatusInternalServerError, "Unknown error")))
			return
		}
		if e.Code() >= errors.InvalidTypeCode {
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			rw.WriteHeader(asHTTPCode(int(e.Code())))
		}
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(e))
		}
	case nil:
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write(errorAsJSON(errors.New(http.StatusInternalServerError, "Unknown error")))
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != http.MethodHead {
			_, _ = rw.Write(errorAsJSON(errors.New(http.StatusInternalServerError, err.Error())))
		}
	}
}

func errorAsJSON(err errors.Error) []byte {
	b, _ := json.Marshal(struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}{err.Code(), err.Error()})
	return b
}

func flattenComposite(errs *errors.CompositeError) *errors.CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *errors.CompositeError:
			if len(e.Errors) > 0 {
				flat := flattenComposite(e)
				if len(flat.Errors) > 0 {
					res = append(res, flat.Errors...)
				}
			}
		default:
			if e != nil {
				res = append(res, e)
			}
		}
	}
	return errors.CompositeValidationError(res...)
}

func asHTTPCode(input int) int {
	if input >= 600 {
		return DefaultHTTPCode
	}
	return input
}
