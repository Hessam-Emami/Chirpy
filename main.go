package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const PORT = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileserverHits int
}

func main() {
	fmt.Println("App Starting")
	cfg := apiConfig{}
	mux := chi.NewRouter()
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))
	mux.HandleFunc("/healthz", HandleHealthPath)
	mux.HandleFunc("/metrics", cfg.HandleMetrics)
	mux.HandleFunc("/reset", cfg.HandleReset)
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}
	err := server.ListenAndServe()
	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(err)
}
func HandleHealthPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	st := "Hits: " + strconv.Itoa(cfg.fileserverHits)
	w.Write([]byte(st))
}

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

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
