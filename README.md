# pomelo

轻量化go api框架

## Quick start
```go
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

	s := pomelo.Default()
	s.Add("/", hello)
	s.Run()
```