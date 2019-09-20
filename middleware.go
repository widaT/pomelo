package pomelo

import (
	"net/http"
)

type middleware func(http.Handler) http.Handler

type Router struct {
	mux   map[string]http.Handler
	chain []middleware
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}

func (r *Router) Add(path string, h http.Handler) {

	var mergedHandler = h

	for i := len(r.chain) - 1; i >= 0; i-- {
		mergedHandler = r.chain[i](mergedHandler)
	}

	r.mux[path] = mergedHandler
}

func (r *Router) Run(path string, wr http.ResponseWriter, req *http.Request) {

	if h, found := r.mux[path]; found {
		h.ServeHTTP(wr, req)
	}
}

func (r *Router) Use(m middleware) {
	r.chain = append(r.chain, m)
}
