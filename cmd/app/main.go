package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/AzaWoodyy/go_backend/internal/services"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "database/sql"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func championsHandler(ddragonSvc *services.DDragonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /champions")

		latestVersion, err := ddragonSvc.GetLatestVersion()
		if err != nil {
			log.Printf("ERROR: Failed to get latest version: %v", err)
			http.Error(w, "Internal Server Error: Could not fetch latest version", http.StatusInternalServerError)
			return
		}
		log.Printf("Latest version found: %s", latestVersion)

		championData, err := ddragonSvc.GetChampions(latestVersion)
		if err != nil {
			log.Printf("ERROR: Failed to get champion data for version %s: %v", latestVersion, err)
			http.Error(w, "Internal Server Error: Could not fetch champion data", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully fetched data for %d champions", len(championData.Data))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(championData)
		if err != nil {
			log.Printf("ERROR: Failed to encode champion data to JSON: %v", err)
		}
	}
}

func main() {
	ddragonSvc := services.NewDDragonService()
	log.Println("DDragon service initialized.")

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/champions", championsHandler(ddragonSvc))
	log.Println("Registered HTTP handlers for / and /champions")

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	listenAddr := fmt.Sprintf(":%s", appPort)
	log.Printf("Server starting on port %s...", appPort)

	// Example: Reading DB config (keep commented if not used yet)
	// dbHost := os.Getenv("MYSQL_HOST")
	// dbPort := os.Getenv("MYSQL_PORT_INTERNAL")
	// dbUser := os.Getenv("MYSQL_USER")
	// dbPassword := os.Getenv("MYSQL_PASSWORD")
	// dbName := os.Getenv("MYSQL_DATABASE")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	// log.Printf("Database DSN (if used): %s", dsn)
	// // db, err := sql.Open("mysql", dsn) ... connect and handle error ...

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("FATAL: Could not start server: %s", err)
	}
}
