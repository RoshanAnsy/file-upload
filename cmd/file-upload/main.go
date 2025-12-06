package main

import (
	"context"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/roshanansy/file-upload/internal/database"
	fileupload "github.com/roshanansy/file-upload/internal/http/handlers/fileUpload"
	User "github.com/roshanansy/file-upload/internal/http/handlers/user"
	"github.com/roshanansy/file-upload/internal/http/middleware/authorized"
	"github.com/roshanansy/file-upload/internal/model"
	// ws "github.com/roshanansy/file-upload/internal/websocket/server"
)


func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (or set specific domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found, using system environment")
	}

	database.ConnectDatabase()
	database.DB.AutoMigrate(&model.User{}, &model.FileStore{}, &model.GenerateApiKey{})

	router := http.NewServeMux()

	router.Handle("POST /api/upload", authorized.Authorized(http.HandlerFunc(fileupload.UploadHandler)))
	router.Handle("POST /api/generateapikey", authorized.Authorized(http.HandlerFunc(fileupload.GenerateApiKeyHandler)))
	router.HandleFunc("POST /api/signup", User.CreateUser)
	router.HandleFunc("POST /api/login", User.LoginUser)
	router.HandleFunc("POST /api/resetPassword", User.ResetPassword)
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to File Upload Service"))
	})

	// Add CORS middleware
	server := http.Server{
		Addr:    ":8080",
		Handler: CORSMiddleware(router),
	}

	slog.Info("Starting server on :8080")

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Server failed", "error", err)
		}
	}()

	<-done
	slog.Info("Server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped")
}
