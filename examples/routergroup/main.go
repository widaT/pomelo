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

func main() {
	s := pomelo.Default()
	s.Add("/", hello)
	group := s.Group("/api")
	group.Add("/abc", hellogroup)
	group.Add("/cdb", hellogroup)
	s.Run()
}
