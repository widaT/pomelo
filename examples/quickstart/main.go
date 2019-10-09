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

func main() {
	s := pomelo.Default()
	s.Add("/", hello)
	s.Add("/dddd", hello)
	s.Run()
}
