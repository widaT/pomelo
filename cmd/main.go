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
	s := pomelo.Default(pomelo.LogMaxSize(2000), pomelo.ALog("log/ac.log"))
	s.Add("/", hello)
	s.Run()
}
