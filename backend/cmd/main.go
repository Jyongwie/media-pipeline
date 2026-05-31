package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/Jyongwie/media-pipeline/backend/internal/infrastructure"
	"github.com/Jyongwie/media-pipeline/backend/internal/presentation"
	"github.com/Jyongwie/media-pipeline/backend/internal/worker"
)

// NEW: The CORS Middleware Function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This tells the browser: "Allow requests from any website (*)"
		// (In a real banking app, we would restrict this strictly to your Vercel URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Browsers send a secret "OPTIONS" request first to check if it's safe. 
		// We intercept it and say "Yes, it's safe!"
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Initialize the database connection
	ctx := context.Background()
	dbConnString := os.Getenv("postgresql://neondb_owner:npg_LuqCrKSihU68@ep-silent-rice-aqplyc8q-pooler.c-8.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require")

	if dbConnString == "" {
		dbConnString = "postgres://admin:secretpassword@localhost:5432/mediadb"
		fmt.Println("Using local Docker database connection...")
	} else {
		fmt.Println("Using Production database connection...")
	}
	repo, err := infrastructure.NewRepository(ctx, dbConnString)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// 2. DEFER THE CLOSE! 
	// This guarantees the pool closes when main() exits, preventing memory leaks.
	defer repo.Close()

	renderHandler := presentation.NewRenderHandler(repo)

	fmt.Println("Database connected successfully!")
	mux := http.NewServeMux()

	// A simple health check endpoint
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "API is running", "environment": "dev"}`))
	})

	mux.HandleFunc("POST /api/jobs", renderHandler.CreateJob)
	mux.HandleFunc("GET /api/jobs", renderHandler.GetJobs)
	worker.StartRenderPool(repo)

	fmt.Println("Backend server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}