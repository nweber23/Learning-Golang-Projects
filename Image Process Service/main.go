package main

import (
	"log"
	"net/http"

	"image-process-service/config"
	"image-process-service/handlers"
	"image-process-service/middleware"
	"image-process-service/models"
	"image-process-service/processor"
	"image-process-service/storage"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.LoadConfig()

	store := models.NewStore()
	fileStorage := storage.NewLocalStorage(cfg.StoragePath)
	proc := processor.NewProcessor()
	authHandler := handlers.NewAuthHandler(store, cfg.JWTSecret)
	imageHandler := handlers.NewImageHandler(store, fileStorage, proc)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitPerHour)

	router := mux.NewRouter()

	// Serve static files (CSS, JS)
	router.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "styles.css")
	})
	router.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "app.js")
	})

	// Serve index.html for root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}).Methods("GET")

	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	imageRoutes := router.PathPrefix("/images").Subrouter()
	imageRoutes.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	imageRoutes.HandleFunc("", imageHandler.UploadImage).Methods("POST")
	imageRoutes.HandleFunc("", imageHandler.ListImages).Methods("GET")
	imageRoutes.HandleFunc("/{id}", imageHandler.GetImage).Methods("GET")
	imageRoutes.HandleFunc("/{id}", imageHandler.DeleteImage).Methods("DELETE")

	transformRoutes := imageRoutes.PathPrefix("/{id}/transform").Subrouter()
	transformRoutes.Use(rateLimiter.RateLimitMiddleware)
	transformRoutes.HandleFunc("", imageHandler.TransformImage).Methods("POST")

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: corsMiddleware(router),
	}
	log.Printf("Starting server on http://localhost:%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
