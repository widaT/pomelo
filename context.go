package pomelo

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var PARAM_NOT_FOUND = errors.New("param not found")

type Context struct {
	Request        *http.Request
	startTime      time.Time
	params         map[string]string
	kv             map[string]interface{}
	server         *Server
	size           int
	statusCode     int
	responseWriter http.ResponseWriter
}

func (ctx *Context) GetServerConf() *Config {
	return ctx.server.conf
}

func (ctx *Context) GetSize() int {
	return ctx.size
}

func (ctx *Context) GetStatusCode() int {
	return ctx.statusCode
}

func (ctx *Context) ParamGetInt(key string) (int, error) {
	if ret, found := ctx.params[key]; found {
		return strconv.Atoi(ret)
	}
	return 0, PARAM_NOT_FOUND
}

func (ctx *Context) GetStime() time.Time {
	return ctx.startTime
}

func (ctx *Context) ParamGetIntWithDefault(key string, defalutVal int) int {
	if ret, found := ctx.params[key]; found {
		n, err := strconv.Atoi(ret)
		if err != nil {
			return defalutVal
		}
		return n
	}
	return defalutVal
}

func (ctx *Context) GetParams() map[string]string {
	return ctx.params
}

func (ctx *Context) ParamGet(key string) string {
	return ctx.params[key]
}

func (ctx *Context) ParamGetWithDefault(key, defaultVal string) string {
	if value, found := ctx.params[key]; found {
		return value
	}
	return defaultVal
}

func (ctx *Context) AddKv(key string, v interface{}) {
	ctx.kv[key] = v
}

func (ctx *Context) GetKv(key string) interface{} {
	return ctx.kv[key]
}

func (ctx *Context) JSON(content interface{}) {
	ctx.ContentType("application/json")
	b, _ := json.Marshal(content)
	ctx.Write(http.StatusOK, b)
}

func (ctx *Context) Write(statusCode int, body []byte) {
	ctx.statusCode = statusCode
	ctx.responseWriter.WriteHeader(statusCode)
	size := 0
	if len(body) > 0 {
		size, _ = ctx.responseWriter.Write(body)
	}
	ctx.size += size
}

func (ctx *Context) STRING(content string) {
	ctx.Write(http.StatusOK, Str2byte(content))
}

func (ctx *Context) BYTE(content []byte) {
	ctx.Write(http.StatusOK, content)
}

func (ctx *Context) Redirect(status int, url string) {
	ctx.responseWriter.Header().Set("Location", url)
	ctx.Write(status, Str2byte("Redirecting to: "+url))
}

func (ctx *Context) NotModified() {
	ctx.Write(http.StatusNotModified, nil)
}

func (ctx *Context) NotFound(message string) {
	ctx.Write(http.StatusNotFound, Str2byte(message))
}

func (ctx *Context) Unauthorized() {
	ctx.Write(http.StatusUnauthorized, nil)
}

func (ctx *Context) Forbidden() {
	ctx.Write(http.StatusForbidden, nil)
}

func (ctx *Context) ContentType(val string) string {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		ctx.responseWriter.Header().Set("Content-Type", ctype)
	}
	return ctype
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.responseWriter.Header().Set(hdr, val)
	} else {
		ctx.responseWriter.Header().Add(hdr, val)
	}
}
