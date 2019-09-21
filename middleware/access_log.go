package middleware

import (
	"fmt"
	"github.com/widaT/pomelo"
	"time"
)

func AccessLog(next pomelo.Handler) pomelo.Handler {
	return pomelo.HandlerFunc(func(ctx *pomelo.Context) {
		stime := ctx.GetStime()
		next.Serve(ctx)
		timeElapsed := time.Since(stime)
		uri := ctx.Request.URL.Path
		fmt.Printf("[pomelo]%s %s %s %s %s %s \n", pomelo.RealIp(ctx.Request), time.Now(), ctx.Request.Method, timeElapsed, uri, ctx.GetParams())
	})
}
