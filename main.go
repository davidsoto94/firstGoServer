package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

const html = `
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>

`

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(html, cfg.fileserverHits)))
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	adminRouter.Get("/metrics/*", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/validate_chirp", apiCfg.validateChipr)
	apiRouter.Get("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK.\n")
	})

	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)
	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}
	log.Printf("starting server...\n")
	log.Printf("listening: http://localhost:8080\n")
	log.Fatal(server.ListenAndServe())

}
