package pomelo

import (
	"github.com/urfave/cli"
	"log"
	"net"
	"net/http"
	"os"
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

func (s *Server) Use(m middleware) {
	s.r.Use(m)
}

func (s *Server) Add(path string, h http.Handler) {
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
		cli.Int64Flag{
			Name:   "level",
			Value:  3,
			Usage:  "log level",
			EnvVar: "POMELO_LOG_LEVEL",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		s.conf = NewConfig(
			Address(ctx.String("address")),
			LogLevel(ctx.Int("level")),
			EnableGzip(ctx.Bool("gzip")),
		)
		return nil
	}
	app.Run(os.Args)
}

func (s *Server) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	s.r.Run(req.URL.Path, c, req)
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
