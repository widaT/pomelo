package pomelo

type Handler interface {
	Serve(*Context)
}

type HandlerFunc func(ctx *Context)

func (f HandlerFunc) Serve(ctx *Context) {
	f(ctx)
}

type Middleware func(Handler) Handler

type Router struct {
	mux   map[string]Handler
	chain []Middleware
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]Handler),
	}
}

func (r *Router) Add(path string, h Handler) {
	var mergedHandler = h
	for i := len(r.chain) - 1; i >= 0; i-- {
		mergedHandler = r.chain[i](mergedHandler)
	}
	r.mux[path] = mergedHandler
}

func (r *Router) Run(path string, ctx *Context) {
	if h, found := r.mux[path]; found {
		h.Serve(ctx)
		return
	}
	r.HttpNotFound(ctx)
}

func (r *Router) HttpNotFound(ctx *Context) {
	//@todo write a log to err_log
	h := HandlerFunc(NotFound)
	h.Serve(ctx)
}

func (r *Router) Use(m Middleware) {
	r.chain = append(r.chain, m)
}
