package pomelo

import (
	"github.com/urfave/cli"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	conf *Config
	r    *Router
}

func NewServer() *Server {
	return &Server{
		r: NewRouter(),
	}
}

func (s *Server) Router() *Router {
	return s.r
}

func (s *Server) Use(m Middleware) {
	s.r.Use(m)
}

func (s *Server) Add(path string, h Handler) {
	s.r.Add(path, h)
}

func (s *Server) Init() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "address",
			Value:  "0.0.0.0:8080",
			Usage:  "listen ip",
			EnvVar: "POMELO_LISTEN_IP",
		},
		cli.BoolFlag{
			Name:   "gzip",
			Usage:  "enable gzip",
			EnvVar: "POMELO_ENABLE_GZIP",
		},
		cli.BoolFlag{
			Name:   "multiform",
			Usage:  "parse multipart form",
			EnvVar: "POMELO_MULTI_FORM",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		s.conf = NewConfig(
			Address(ctx.String("address")),
			EnableGzip(ctx.Bool("gzip")),
			ParseMultiForm(ctx.Bool("multiform")),
		)
		return nil
	}

	app.Run(os.Args)
}

func (s *Server) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:        req,
		startTime:      time.Now(),
		server:         s,
		ResponseWriter: c,
		params:         make(map[string]string),
		kv:             make(map[string]interface{}),
	}

	//copy data from get and post(x-www-form-urlencoded)
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			ctx.params[k] = v[0]
		}
	}

	if s.conf.ParseMultiForm {
		//copy dat from post(multipart/form-data)
		req.ParseMultipartForm(32 << 20)
		if len(req.PostForm) > 0 {
			for k, v := range req.PostForm {
				if len(ctx.params[k]) > 0 {
					ctx.params[k] = v[0]
				}
			}
		}
	}

	s.r.Run(req.URL.Path, ctx)
}

func (s *Server) Run() {
	l, err := net.Listen("tcp", s.conf.Address)
	if err != nil {
		log.Fatal(err)
	}
	err = http.Serve(l, s)
	if err != nil {
		log.Fatal(err)
	}
}
