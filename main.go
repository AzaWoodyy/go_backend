package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	
	_ "github.com/go-sql-driver/mysql"
	_ "database/sql"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	// Example: Reading DB config (though not connecting yet)
	// dbHost := os.Getenv("MYSQL_HOST") // This will be 'db' from compose.yml
	// dbPort := os.Getenv("MYSQL_PORT_INTERNAL") // Docker internal port
	// dbUser := os.Getenv("MYSQL_USER")
	// dbPassword := os.Getenv("MYSQL_PASSWORD")
	// dbName := os.Getenv("MYSQL_DATABASE")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	// fmt.Printf("Would connect with DSN: %s\n", dsn)
	// // db, err := sql.Open("mysql", dsn) ... connect and handle error ...

	http.HandleFunc("/", helloHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Default port if not set
	}

	listenAddr := fmt.Sprintf(":%s", appPort)
	log.Printf("Server starting on port %s\n", appPort)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
