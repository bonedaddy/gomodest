package pkg

import (
	"net/http"

	"github.com/go-chi/render"
)

type Error interface {
	error
	Status() int
}

// HTTPErr represents an http error
type HTTPErr struct {
	UserMessage string
	Code        int
}

// Error satisfies the error interface
func (h HTTPErr) Error() string {
	return h.UserMessage
}

// Status return the http status code
func (h HTTPErr) Status() int {
	return h.Code
}

var InternalErr = HTTPErr{
	UserMessage: "Internal Error",
	Code:        500,
}

var BadRequest = HTTPErr{
	UserMessage: "Invalid Request",
	Code:        400,
}

// copied from https://github.com/go-chi/chi/blob/master/_examples/rest/main.go#L389
//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal error.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
