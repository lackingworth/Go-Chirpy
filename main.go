package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	PORT         = "8080"
	FilepathRoot = "."
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load()))); err != nil {
		log.Printf("Failed to write response from handlerMetrics: %v", err)
	}
}

func main() {

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FilepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", FilepathRoot, PORT)
	log.Fatal(server.ListenAndServe())
}
