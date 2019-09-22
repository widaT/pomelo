package pomelo

import (
	"fmt"
	"net/url"
	"time"
)

func AccessLog(next Handler) Handler {
	return HandlerFunc(func(ctx *Context) {
		startTime := ctx.GetStime()
		next.Serve(ctx)
		timeElapsed := time.Since(startTime)
		uri := ctx.Request.URL.Path
		status := ctx.GetHttpStatus()
		size := ctx.GetSize()

		paramsCopy := url.Values{}
		if len(ctx.params) > 0 {
			for key, param := range ctx.GetParams() {
				if len(param) > 500 {
					paramsCopy.Set(key, "-5h-") //param value size bigger than 500  then hidden
					continue
				}
				paramsCopy.Set(key, param)
			}
		}

		fmt.Printf("[pomelo] %v |%3d |%s | %13v | %15s | %s | %d |%v \n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			status, timeElapsed, RealIp(ctx.Request), ctx.Request.Method,
			uri, size, paramsCopy.Encode())
	})
}
