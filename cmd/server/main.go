package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"myapp/internal/handler"
	"myapp/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func runMigrations(db *sql.DB, migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return err
			}

			log.Printf("Applied migration: %s", file.Name())
		}
	}

	return nil
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set up migrations
	migrationsPath := "./migrations/"
	if err := runMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	// Create instances of the service and handler
	commandService := service.NewCommandService(db)
	commandHandler := handler.NewCommandHandler(commandService)

	// Set up the HTTP server with Gorilla Mux router
	r := mux.NewRouter()

	// Define API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/commands", commandHandler.CreateCommand).Methods("POST")
	api.HandleFunc("/commands", commandHandler.GetCommands).Methods("GET")
	api.HandleFunc("/commands/{id}", commandHandler.GetCommand).Methods("GET")

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is listening on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}