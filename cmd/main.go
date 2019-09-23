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
		B: "eeee",
	})

}

func main() {
	s := pomelo.Default()
	s.Add("/", hello)
	s.Run()
}
