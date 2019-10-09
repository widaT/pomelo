package main

import (
	"github.com/widaT/pomelo"
)

func hello(c *pomelo.Context) {
	c.JSON(struct {
		A int    `json:"a"`
		B string `json:"b"`
	}{
		A: 1233,
		B: "dddddd",
	})
}

func hellogroup(c *pomelo.Context) {
	c.JSON(struct {
		A int    `json:"a"`
		B string `json:"b"`
	}{
		A: 1233,
		B: "hellogroup",
	})
}

func Test(next pomelo.Handler) pomelo.Handler {
	return pomelo.HandlerFunc(func(ctx *pomelo.Context) {
		if ctx.Param("a") != "b" {
			ctx.JSON(struct {
				A int    `json:"a"`
				B string `json:"b"`
			}{
				A: 3,
				B: "dddddddddddd warning",
			})
			return
		}
		next.Serve(ctx)
	})
}

func main() {
	s := pomelo.Default()
	s.Add("/", hello)

	group := s.Group("/api")

	group.Use(Test)
	group.Add("/abc", hellogroup)

	group.Add("/cdb", hellogroup)

	s.Add("/dddd", hello)
	s.Run()
}
