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
	conf      *Config
	r         *Router
	errLogger Elogger
}

func NewServer() *Server {
	s := &Server{}
	s.init()
	s.r = NewRouter(s)
	s.errLogger = NewErrLogger(s.conf.ErrLog)
	return s
}

//with accesslog middleware
func Default() *Server {
	s := NewServer()
	s.Use(AccessLog)
	return s
}

func (s *Server) Router() *Router {
	return s.r
}

func (s *Server) Use(m Middleware) {
	s.r.Use(m)
}

func (s *Server) Add(path string, h interface{}) {
	s.r.Add(path, h)
}

func (s *Server) init() {
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
		cli.StringFlag{
			Name:   "elog",
			Value:  "",
			Usage:  "err log path",
			EnvVar: "POMELO_ERRLOG_PATH",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		s.conf = NewConfig(
			Address(ctx.String("address")),
			EnableGzip(ctx.Bool("gzip")),
			ParseMultiForm(ctx.Bool("multiform")),
			ELog(ctx.String("elog")),
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
		responseWriter: c,
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
		s.errLogger.Error("server listen err %#v", err)
		log.Fatal(err)
	}
	err = http.Serve(l, s)
	if err != nil {
		s.errLogger.Error("http Serve err %#v", err)
		log.Fatal(err)
	}
}
