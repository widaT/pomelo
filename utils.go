package pomelo

import (
	"net/http"
	"strings"
	"unsafe"
)

func Str2byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func RealIp(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		return strings.Split(xf, ",")[0]
	}
	ip := r.RemoteAddr
	pos := strings.LastIndex(ip, ":")
	var clientIp string
	if pos > 0 {
		clientIp = ip[0:pos]
	} else {
		clientIp = ip
	}
	return clientIp
}
