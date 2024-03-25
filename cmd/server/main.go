package main

//Все необходимые библиотеки, а также импортируем внутренние зависимости проекта
import (
	"database/sql"
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

// Функция для запуска миграци(делается при запуске кода)
func runMigrations(db *sql.DB, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return err
			}

			log.Printf("Миграции выполнены по файлу: %s", file.Name())
		}
	}

	return nil
}

func main() {
	// Загружаем зависимости из .env файла
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Установили миграции
	migrationsPath := "./migrations/"
	if err := runMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Создили зависимости для сервиса и ручек
	commandService := service.NewCommandService(db)
	commandHandler := handler.NewCommandHandler(commandService)

	// Задали роуты с помощью Gorilla Mux
	r := mux.NewRouter()

	// Определеяем наши роуты
	api := r.PathPrefix("/api/v1").Subrouter()
	go api.HandleFunc("/commands", commandHandler.CreateCommand).Methods("POST")
	go api.HandleFunc("/commands", commandHandler.GetCommands).Methods("GET")
	go api.HandleFunc("/commands/{id}", commandHandler.GetCommand).Methods("GET")

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is listening on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
