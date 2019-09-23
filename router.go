package pomelo

import (
	"log"
	"runtime"
)

type Handler interface {
	Serve(*Context)
}

type HandlerFunc func(ctx *Context)

func (f HandlerFunc) Serve(ctx *Context) {
	f(ctx)
}

type Middleware func(Handler) Handler

type Router struct {
	server *Server
	mux    map[string]Handler
	chain  []Middleware
}

func NewRouter(server *Server) *Router {
	return &Router{
		server: server,
		mux:    make(map[string]Handler),
	}
}

func (r *Router) Add(path string, h interface{}) {
	var handler Handler
	switch h.(type) {
	case func(*Context):
		handler = HandlerFunc(h.(func(*Context)))
	case Handler:
		handler = h.(Handler)
	default:
		r.server.errLogger.Log("add route %s error", path)
		log.Fatal("add route error")
	}
	for i := len(r.chain) - 1; i >= 0; i-- {
		handler = r.chain[i](handler)
	}
	r.mux[path] = handler
}

func (r *Router) Run(path string, ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			r.server.errLogger.Log("Handler crashed with error %#v", err)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				r.server.errLogger.Log(file, line)
			}

		}
	}()

	if h, found := r.mux[path]; found {
		h.Serve(ctx)
		return
	}
	r.HttpNotFound(ctx)
}

func (r *Router) HttpNotFound(ctx *Context) {
	r.server.errLogger.Log("The requested URL %s was not found", ctx.Request.URL.Path)
	h := HandlerFunc(NotFound)
	h.Serve(ctx)
}

func (r *Router) Use(m ...Middleware) {
	r.chain = append(r.chain, m...)
}
