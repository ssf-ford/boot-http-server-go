package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"santnas/boot-http-server-course/internal/database"
)

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlTemplate, c.fileServerHits.Load())
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
	db             *database.Queries
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
