package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/widaT/pomelo"
	middleware_jwt "github.com/widaT/pomelo/middleware/jwt"
	"time"
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

	s.Add("/api", func(c *pomelo.Context) {
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims = jwt.MapClaims{
			"Id":  "wida",
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		}
		tokenString, err := token.SignedString(pomelo.Str2byte(middleware_jwt.PomeloJWTSecret))
		if err != nil {
			c.JSON(struct {
				errcode int    `json:"errcode"`
				msg     string `json:"msg"`
			}{
				errcode: -1,
				msg:     err.Error(),
			})
		}
		c.JSON(struct {
			Token string `json:"token"`
		}{Token: tokenString})
	})

	group := s.Group("/api")
	group.Use(middleware_jwt.Jwt)
	group.Add("/abc", hellogroup)
	group.Add("/cdb", hellogroup)
	s.Add("/dddd", hello)
	s.Run()
}
