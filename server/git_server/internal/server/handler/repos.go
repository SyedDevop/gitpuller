package handler

import (
	"errors"
	"fmt"
	"git_server/internal/file"
	"git_server/internal/server/middleware"
	reserr "git_server/internal/server/res_err"
	"math"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
)

type Repos struct{}

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
	chunk := math.Ceil(float64(dataLen) / float64(perPage))
	fromIdx := perPage * (page - 1)
	toIdx := perPage * page
	if toIdx > dataLen {
		toIdx = dataLen
	}

	w.Header().Set("Link", fmt.Sprintf("%0.0f", chunk))
	render.JSON(w, r, data[fromIdx:toIdx])
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
