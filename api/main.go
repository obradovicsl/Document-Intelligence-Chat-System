package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/handlers"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/logger"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/repository"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/services"
)

func main() {
    // Load env
    if err := godotenv.Load(); err != nil {
        slog.Warn("no .env file found, using system environment variables")
	}
    
    log := logger.New(os.Getenv("ENV"))
    slog.SetDefault(log)

    // Initialize database
    if err := repository.Init(); err != nil {
        slog.Error("Failed to initialize database:", "error", err)
        os.Exit(1)
    }
    
    s3Service, err := services.InitS3()
    if err != nil {
        slog.Error("Failed to initialize database:", "error", err)
        os.Exit(1)
    }

    auth.Init()
    
    docRepo := repository.NewDocumentRepository(repository.DB)
    docService := services.NewDocumentService(docRepo, s3Service)
    docHandler := handlers.NewDocumentHandler(docService)

    
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    }).Methods("GET")


    // Protected routes
    api := r.PathPrefix("/api").Subrouter()
    api.Use(auth.AuthMiddleware)

    // Document upload endpoints
    api.HandleFunc("/upload/init", docHandler.HandleInitUpload).Methods("POST")
    api.HandleFunc("/upload/complete", docHandler.HandleCompleteUpload).Methods("POST")
    
    // 
    // api.HandleFunc("/documents", handlers.GetDocuments).Methods("GET")
    // api.HandleFunc("/documents/{id}", handlers.GetDocument).Methods("GET")
    // api.HandleFunc("/chat", handlers.Chat).Methods("POST")

    // CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3001"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Authorization", "Content-Type"},
    })

    handler := c.Handler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    slog.Info("starting server", "port", os.Getenv("PORT"), "env", os.Getenv("ENV"))
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        slog.Error("server failed to start",
        "error", err,
        "port", port)
        os.Exit(1)
    }
}