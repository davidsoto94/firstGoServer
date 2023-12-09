package main

import (
	"io"
	"log"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
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
