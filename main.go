package main

import (
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"time"
)

func mockServer(httpCode int) *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpCode)
		xray.NewSegmentFromHeader(r.Context(), "Subsegment", r, header.FromString("newdata"))
		time.Sleep(3*time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
	}))
	return s
}

func main() {
	router := mux.NewRouter()
	router.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("newsegment"), HelloServer()))
	router.PathPrefix("/debug/pprof").HandlerFunc(pprof.Index)

	log.Panic(http.ListenAndServe("localhost:8067", router))
}

func HelloServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client = xray.Client(http.DefaultClient)
		ctx := r.Context()
		srv := mockServer(200)
		for jj:=3; jj > 0; jj-- {
			go func() {
				newCtx := xray.DetachContext(ctx)
				for i:=0; i < 100; i++ {
					req, _ := http.NewRequestWithContext(newCtx, "GET", srv.URL, nil)
					client.Do(req)
				}
			}()
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}