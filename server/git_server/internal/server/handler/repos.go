package handler

import (
	"errors"
	"fmt"
	"git_server/internal/file"
	"git_server/internal/server/middleware"
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

	perPage, ok := r.Context().Value(middleware.PerPageKey).(int)
	if !ok {
		log.Error("ReposHandlear", "err", errors.New("perPage key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return
	}
	page, ok := r.Context().Value(middleware.PageKey).(int)
	if !ok {
		log.Error("ReposHandlear", "err", errors.New("page key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return
	}

	dataLen := len(data)
	if perPage != 0 {
		w.Header().Set("Link", genrateLinks(perPage, page, dataLen))
	}
	if perPage == 0 {
		perPage = 30
	}

	fromIdx := perPage * (page - 1)
	toIdx := perPage * page
	if toIdx > dataLen {
		toIdx = dataLen
	}

	render.JSON(w, r, data[fromIdx:toIdx])
}

func (re *Repos) PagenatedRepos(w http.ResponseWriter, r *http.Request) {
	data, err := file.GetReposJson()
	if err != nil {
		log.Error("ReposHandlear", "err", err)
		render.Render(w, r, reserr.ErrRender(err))
		return
	}

	perPage, ok := r.Context().Value(middleware.PerPageKey).(int)
	if !ok {
		log.Error("ReposHandlear", "err", errors.New("perPage key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return
	}
	page, ok := r.Context().Value(middleware.PageKey).(int)
	if !ok {
		log.Error("ReposHandlear", "err", errors.New("page key is missing or invalid"))
		render.Render(w, r, reserr.ErrInternalServer)
		return
	}

	dataLen := len(data)
	if perPage != 0 {
		w.Header().Set("Link", genrateLinks(perPage, page, dataLen))
	}
	if perPage == 0 {
		perPage = 30
	}

	fromIdx := perPage * (page - 1)
	toIdx := perPage * page
	if toIdx > dataLen {
		toIdx = dataLen
	}

	render.JSON(w, r, data[fromIdx:toIdx])
}

func genLink(page, per_page int, rel string) string {
	return fmt.Sprintf(`<http://localhost:4069/user/repos?per_page=%d&page=%d>; rel="%s"`, per_page, page, rel)
}

func genrateLinks(perPage, page, dataLen int) string {
	chunk := math.Ceil(float64(dataLen) / float64(perPage))
	nextPage := page + 1
	lastPage := int(chunk)

	linkBuffer := strings.Builder{}

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

func (re *Repos) ReposHandlear(w http.ResponseWriter, r *http.Request) {
	data := file.GetReposByte()
	// if err != nil {
	// 	log.Error("ReposHandlear", "err", err)
	// 	render.Render(w, r, ErrRender(err))
	// 	return
	// }
	// render.JSON(w, r, data)

	w.Header().Set("Content-Type", "application/json")
	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(data)
}
