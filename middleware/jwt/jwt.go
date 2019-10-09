package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/widaT/pomelo"
)

const PomeloJWTSecret = "pomelo_secret"

func Jwt(next pomelo.Handler) pomelo.Handler {
	return pomelo.HandlerFunc(func(ctx *pomelo.Context) {
		_, err := request.ParseFromRequest(ctx.Request, request.OAuth2Extractor,
			func(token *jwt.Token) (interface{}, error) {
				b := pomelo.Str2byte(PomeloJWTSecret)
				return b, nil
			})
		if err != nil {
			ctx.Write(401, pomelo.Str2byte(err.Error()))
			return
		}
		next.Serve(ctx)
	})
}
