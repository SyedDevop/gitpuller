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

type GitErrRes struct {
	Message          string `json:"message"`
	Status           string `json:"status"`
	DocumentationUrl string `json:"documentation_url"`
	HTTPStatusCode   int    `json:"-"`
}

func (g *GitErrRes) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, g.HTTPStatusCode)
	return nil
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

var (
	ErrNotFound       = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
	ErrInternalServer = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal server error."}
	ErrGitTree        = &GitErrRes{
		HTTPStatusCode:   404,
		Status:           "404",
		Message:          "Not Found",
		DocumentationUrl: "https://docs.github.com/rest/git/trees#get-a-tree",
	}
)
