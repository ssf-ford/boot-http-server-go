package main

import (
	"net/http"
)

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/healthz" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
