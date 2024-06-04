package reserr

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error  `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
	HTTPStatusCode int    `json:"-"`
	AppCode        int64  `json:"code,omitempty"`
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

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

var ErrInternalServer = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal server error."}
