package mware

import (
	"context"
	"net/http"
	"strconv"
)

type (
	PerPage struct {
		name string
	}
	Page struct {
		name string
	}
)

func (p *PerPage) String() string {
	return p.name
}

func (p *Page) String() string {
	return p.name
}

var (
	PerPageKey = &PerPage{"per_page"}
	PageKey    = &Page{"page"}
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perPageStr := r.URL.Query().Get(PerPageKey.String())
		pageStr := r.URL.Query().Get(PageKey.String())

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage < 1 {
			perPage = 0
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, PerPageKey, perPage)
		ctx = context.WithValue(ctx, PageKey, page)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
