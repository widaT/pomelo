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
	Request *http.Request

	startTime time.Time
	params    map[string]string
	kv        map[string]interface{}
	server    *Server
	http.ResponseWriter
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
	ctx.Write(b)
}

func (ctx *Context) STR(content string) {
	ctx.ResponseWriter.Write(Str2byte(content))
}

func (ctx *Context) BYTE(content []byte) {
	ctx.ResponseWriter.Write(content)
}

func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write(Str2byte(body))
}

func (ctx *Context) Redirect(status int, url string) {
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write(Str2byte("Redirecting to: " + url))
}

func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write(Str2byte(message))
}

func (ctx *Context) Unauthorized() {
	ctx.ResponseWriter.WriteHeader(401)
}

func (ctx *Context) Forbidden() {
	ctx.ResponseWriter.WriteHeader(403)
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
		ctx.Header().Set("Content-Type", ctype)
	}
	return ctype
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.Header().Set(hdr, val)
	} else {
		ctx.Header().Add(hdr, val)
	}
}
