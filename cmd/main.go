package main

import (
	"github.com/widaT/pomelo"
	"github.com/widaT/pomelo/middleware"
)

func hello(c *pomelo.Context) {
	c.JSON(struct {
		A int    `json:"a"`
		B string `json:"b"`
	}{
		A: 1233,
		B: "eeee",
	})

}

func timeMiddleware(next pomelo.Handler) pomelo.Handler {
	return pomelo.HandlerFunc(func(ctx *pomelo.Context) {
		//timeStart := time.Now()
		if ctx.ParamGet("name") == "h" {
			ctx.STR("aaaa")
			return
		}

		next.Serve(ctx)
		//timeElapsed := time.Since(timeStart)
		//logger.Println(timeElapsed)
		//fmt.Println(timeElapsed)
	})
}

func main() {
	s := pomelo.NewServer()
	s.Init()
	s.Use(middleware.AccessLog)
	s.Use(timeMiddleware)
	s.Add("/", pomelo.HandlerFunc(hello))
	s.Run()
}
