package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset/", apiCfg.handlerReset)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK.\n")
	})

	corsMux := middlewareCors(mux)
	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}
	log.Printf("starting server...\n")
	log.Printf("listening: http://localhost:8080\n")
	log.Fatal(server.ListenAndServe())

}
