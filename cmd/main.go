package main

import (
	"fmt"
	"github.com/widaT/pomelo"
	"net/http"
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

func main() {
	s := pomelo.NewServer()
	s.Init()
	s.Use(timeMiddleware)
	s.Add("/", http.HandlerFunc(hello))
	s.Run()
}
