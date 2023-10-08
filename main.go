package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const PORT = "8080"
const filepathRoot = "."

type apiConfig struct {
	fileserverHits int
}

func main() {
	fmt.Println("App Starting")
	cfg := apiConfig{fileserverHits: 0}
	mainRouter := chi.NewRouter()

	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mainRouter.Handle("/app", fsHandler)
	mainRouter.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", HandleHealthPath)
	apiRouter.Get("/reset", cfg.HandleReset)
	apiRouter.Post("/validate_chirp", HandleValidateChirp)
	mainRouter.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.HandleMetrics)
	mainRouter.Mount("/admin", adminRouter)

	corsMux := middlewareCors(mainRouter)
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsMux,
	}
	err := server.ListenAndServe()
	log.Printf("Serving on port: %s\n", PORT)
	log.Fatal(err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	fmt.Println(string(dat))
	w.WriteHeader(code)
	w.Write(dat)
}

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type BodyDto struct {
		Body string
	}
	// kerfuffle sharbert fornax
	var bodyDto BodyDto
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bodyDto)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		SendError(w, 500, "Something went wrong")
		return
	}
	if len(bodyDto.Body) == 0 {
		SendError(w, 400, "No empty body bitte!")
		return
	}
	if len(bodyDto.Body) > 140 {
		SendError(w, 400, "Chirp is too long")
		return
	}
	notPermitted := []string{"kerfuffle", "sharbert", "fornax"}
	splited := strings.Split(bodyDto.Body, " ")
	for index, verb := range splited {
		for _, v := range notPermitted {
			if strings.ToLower(verb) == v {
				splited[index] = "****"
			}
		}
	}
	type Valide struct {
		Valid string `json:"cleaned_body"`
	}
	respondWithJSON(w, 200, Valide{Valid: strings.Join(splited, " ")})
}

func SendError(w http.ResponseWriter, code int, message string) {
	type Error struct {
		error string
	}
	w.WriteHeader(code)
	if len(message) == 0 {
		message = "Something went wrong"
	}
	response, err := json.Marshal(Error{error: message})
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return
	}
	w.Write(response)
}

func HandleHealthPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits)))
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
