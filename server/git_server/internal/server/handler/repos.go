package handler

import (
	"errors"
	"fmt"
	"git_server/internal/file"
	mware "git_server/internal/server/middleware"
	reserr "git_server/internal/server/res_err"
	"math"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Repos struct{}

func (re *Repos) PagenatedUserRepos(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	if user == "" {
		render.Render(w, r, reserr.ErrNotFound)
		return
	}

	data, err := file.ReadJson(user)
	if err != nil {
		log.Error("ReposHandlear", "err", err)
		render.Render(w, r, reserr.ErrRender(err))
		return
	}

	paginateAndRender(w, r, data)
}

// paginateAndRender handles the pagination of the data and renders the appropriate JSON response.
func (re *Repos) PagenatedRepos(w http.ResponseWriter, r *http.Request) {
	data, err := file.GetReposJson()
	if err != nil {
		log.Error("ReposHandlear", "err", err)
		render.Render(w, r, reserr.ErrRender(err))
		return
	}
	paginateAndRender(w, r, data)
}

func paginateAndRender(w http.ResponseWriter, r *http.Request, data []map[string]interface{}) {
	pagin := getPagination(w, r)
	if pagin == nil {
		return
	}
	dataLen := len(data)
	if pagin.PerPage != 0 {
		w.Header().Set("Link", genrateLinks(pagin.PerPage, pagin.Page, dataLen))
	}
	window := getDataWindow(w, r, dataLen, pagin.PerPage, pagin.Page)
	if window == nil {
		return
	}
	render.JSON(w, r, data[window.fromIdx:window.toIdx])
}

// Pagination holds pagination information such as items per page and current page number.
type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}

// getPagination extracts pagination parameters (perPage and page) from the request context.
// It returns a Pagination struct or nil if an error occurs.
func getPagination(w http.ResponseWriter, r *http.Request) *Pagination {
	perPage, ok := r.Context().Value(mware.PerPageKey).(int)
	if !ok {
		log.Error("ReposHandler", "err", errors.New("perPage key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return nil
	}
	page, ok := r.Context().Value(mware.PageKey).(int)
	if !ok {
		log.Error("ReposHandler", "err", errors.New("page key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return nil
	}
	return &Pagination{PerPage: perPage, Page: page}
}

// Window defines the index range for the current page's data slice.
type Window struct {
	fromIdx, toIdx int
}

// getDataWindow calculates the data slice indices for the current page based on the pagination parameters.
// It returns a Window struct or nil if an error occurs.
func getDataWindow(w http.ResponseWriter, r *http.Request, dataLen, perPage, page int) *Window {
	chunk := math.Ceil(float64(dataLen) / float64(perPage))
	if page > int(chunk) {
		render.JSON(w, r, []string{})
		return nil
	}
	if perPage == 0 {
		perPage = 30
	}

	fromIdx := perPage * (page - 1)
	toIdx := perPage * page
	if toIdx > dataLen {
		toIdx = dataLen
	}
	return &Window{fromIdx: fromIdx, toIdx: toIdx}
}

func genLink(page, per_page int, rel string) string {
	return fmt.Sprintf(`<http://localhost:4069/user/repos?per_page=%d&page=%d>; rel="%s"`, per_page, page, rel)
}

func genrateLinks(perPage, page, dataLen int) string {
	chunk := math.Ceil(float64(dataLen) / float64(perPage))
	nextPage := page + 1
	lastPage := int(chunk)

	linkBuffer := strings.Builder{}

	if page > lastPage {
		linkBuffer.WriteString(genLink(lastPage, perPage, "prev"))
		linkBuffer.WriteString(", ")
		linkBuffer.WriteString(genLink(lastPage, perPage, "last"))
		linkBuffer.WriteString(", ")
		linkBuffer.WriteString(genLink(1, perPage, "first"))
		return linkBuffer.String()
	}

	if page < lastPage {
		linkBuffer.WriteString(genLink(nextPage, perPage, "next"))
		linkBuffer.WriteString(", ")
		linkBuffer.WriteString(genLink(lastPage, perPage, "last"))
	}
	if page > 1 {
		if linkBuffer.Len() != 0 {
			linkBuffer.WriteString(", ")
		}
		linkBuffer.WriteString(genLink(page-1, perPage, "prev"))
		linkBuffer.WriteString(", ")
		linkBuffer.WriteString(genLink(1, perPage, "first"))
	}

	return linkBuffer.String()
}
