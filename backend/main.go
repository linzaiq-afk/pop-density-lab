package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type apiServer struct {
	entities []Entity
}

func main() {
	addr := envOr("LISTEN_ADDR", ":8081")
	entities := getDataset()
	mux := newMux(entities)

	// Опциональный режим без Hydrogen: Go-сервер сам раздаёт frontend.
	// В docker-compose статику отдаёт Hydrogen, поэтому STATIC_DIR не задан.
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			log.Fatalf("STATIC_DIR %q is not a directory", dir)
		}
		mux.Handle("GET /", http.FileServer(http.Dir(dir)))
		log.Printf("serving static files from %s", dir)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           withCORS(withLogging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Go backend started on %s (%d entities)", addr, len(entities))
	log.Fatal(srv.ListenAndServe())
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func newMux(entities []Entity) *http.ServeMux {
	api := &apiServer{entities: entities}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", api.handleHealth)
	mux.HandleFunc("GET /api/meta", api.handleMeta)
	mux.HandleFunc("GET /api/entities", api.handleEntities)
	mux.HandleFunc("GET /api/entities/{id}", api.handleEntity)
	return mux
}

func (a *apiServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "pop-density-lab",
	})
}

func (a *apiServer) handleMeta(w http.ResponseWriter, _ *http.Request) {
	years := make(map[int]struct{})
	countries, cities := 0, 0
	for _, entity := range a.entities {
		switch entity.Type {
		case "country":
			countries++
		case "city":
			cities++
		}
		for _, stat := range entity.History {
			years[stat.Year] = struct{}{}
		}
	}

	list := make([]int, 0, len(years))
	for year := range years {
		list = append(list, year)
	}
	sort.Ints(list)

	writeJSON(w, http.StatusOK, map[string]any{
		"service":   "pop-density-lab",
		"years":     list,
		"count":     len(a.entities),
		"countries": countries,
		"cities":    cities,
	})
}

func (a *apiServer) handleEntities(w http.ResponseWriter, r *http.Request) {
	filterType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

	result := make([]Entity, 0, len(a.entities))
	for _, entity := range a.entities {
		if filterType != "" && filterType != "all" && entity.Type != filterType {
			continue
		}
		if query != "" &&
			!strings.Contains(strings.ToLower(entity.Name), query) &&
			!strings.Contains(strings.ToLower(entity.Country), query) {
			continue
		}
		result = append(result, entity)
	}
	writeJSON(w, http.StatusOK, result)
}

func (a *apiServer) handleEntity(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, entity := range a.entities {
		if entity.ID == id {
			writeJSON(w, http.StatusOK, entity)
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{
		"error": "entity not found",
		"id":    id,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, recorder.status, time.Since(started).Round(time.Millisecond))
	})
}
