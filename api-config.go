package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d\n", c.fileServerHits.Load())
}

func (c *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	c.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(fmt.Sprintf("Hits: %d\n", c.fileServerHits.Load())))
	fmt.Fprintf(w, "Hits: %d\n", c.fileServerHits.Load())
}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
