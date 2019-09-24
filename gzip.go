package pomelo

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipWriter struct {
	writer *gzip.Writer
	http.ResponseWriter
}

func (g *GzipWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func Gzip(next Handler) Handler {
	return HandlerFunc(func(ctx *Context) {
		if needCompress(ctx.Request) {
			gw, err := gzip.NewWriterLevel(ctx.responseWriter, gzip.DefaultCompression)
			if err != nil {
				next.Serve(ctx)
				return
			}
			ctx.responseWriter = &GzipWriter{gw, ctx.responseWriter}
			ctx.SetHeader("Content-Encoding", "gzip", true)
			ctx.SetHeader("Vary", "Accept-Encoding", true)
			next.Serve(ctx)
			ctx.SetHeader("Content-Length", "0", true)
			gw.Close()
			return
		}
		next.Serve(ctx)
	})
}

func needCompress(req *http.Request) bool {
	//@todo should filter image files
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}
