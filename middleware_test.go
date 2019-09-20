package pomelo

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func hello(wr http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	wr.Write([]byte("hello"))
	timeElapsed := time.Since(timeStart)
	fmt.Println(timeElapsed)
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		// next handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		//logger.Println(timeElapsed)
		fmt.Println(timeElapsed)
	})
}

func TestMiddleware(t *testing.T) {
	/*	r := NewRouter()

		r.Use(timeMiddleware)
		r.Add("/", http.HandlerFunc(hello))

		fmt.Println(r)*/

	s := NewServer()
	s.Init()
	s.r.Use(timeMiddleware)
	s.r.Add("/", http.HandlerFunc(hello))
	s.Run()
}
