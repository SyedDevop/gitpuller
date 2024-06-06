package server

import (

	// "log"

	"git_server/internal/server/handler"
	mware "git_server/internal/server/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	reposHan := &handler.Repos{}

	r.Route("/user", func(r chi.Router) {
		r.With(mware.Paginate).Get("/repos", reposHan.PagenatedRepos)
	})

	r.Route("/users", func(r chi.Router) {
		r.With(mware.Paginate).Get("/{user}/repos", reposHan.PagenatedUserRepos)
	})
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	render.JSON(w, r, resp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, s.db.Health())
}
