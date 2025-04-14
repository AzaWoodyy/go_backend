package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AzaWoodyy/go_backend/internal/config"
	"github.com/AzaWoodyy/go_backend/internal/models"
	"github.com/AzaWoodyy/go_backend/internal/repositories"
	"github.com/AzaWoodyy/go_backend/internal/services"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func championsHandler(ddragonSvc *services.DDragonService, championRepo *repositories.ChampionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /champions")

		latestVersion, err := ddragonSvc.GetLatestVersion()
		if err != nil {
			log.Printf("ERROR: Failed to get latest version: %v", err)
			http.Error(w, "Internal Server Error: Could not fetch latest version", http.StatusInternalServerError)
			return
		}
		log.Printf("Latest version found: %s", latestVersion)

		champions, err := ddragonSvc.GetChampions(latestVersion)
		if err != nil {
			log.Printf("ERROR: Failed to get champion data for version %s: %v", latestVersion, err)
			http.Error(w, "Internal Server Error: Could not fetch champion data", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully fetched data for %d champions", len(champions))

		// Save champions to database
		if err := championRepo.SaveChampions(champions); err != nil {
			log.Printf("ERROR: Failed to save champions to database: %v", err)
			http.Error(w, "Internal Server Error: Could not save champion data", http.StatusInternalServerError)
			return
		}

		// Get champions from database to ensure we have the complete data
		dbChampions, err := championRepo.GetChampions()
		if err != nil {
			log.Printf("ERROR: Failed to get champions from database: %v", err)
			http.Error(w, "Internal Server Error: Could not retrieve champion data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := struct {
			Version   string            `json:"version"`
			Champions []models.Champion `json:"champions"`
		}{
			Version:   latestVersion,
			Champions: dbChampions,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode response: %v", err)
			http.Error(w, "Internal Server Error: Could not encode response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Initialize database connection
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}

	// Initialize services and repositories
	ddragonSvc := services.NewDDragonService()
	championRepo := repositories.NewChampionRepository(db)

	// Auto-migrate the schema
	log.Println("Starting database migration...")
	if err := db.AutoMigrate(&models.Champion{}, &models.Tag{}, &models.Version{}); err != nil {
		log.Fatalf("FATAL: Could not migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/champions", championsHandler(ddragonSvc, championRepo))
	log.Println("Registered HTTP handlers for / and /champions")

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	listenAddr := fmt.Sprintf(":%s", appPort)
	log.Printf("Server starting on port %s...", appPort)

	server := &http.Server{
		Addr:         listenAddr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("FATAL: Could not start server: %s", err)
	}
}
