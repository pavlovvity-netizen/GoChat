package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"GoChat/handler"
	"GoChat/repository"
	"GoChat/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используются переменные окружения по умолчанию")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	repo := repository.NewPostgresRepository(db)
	svc := service.NewMessageService(repo)
	h := handler.NewMessageHandler(svc)

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewUserService(userRepo)
	authHandler := handler.NewUserHandler(authSvc)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("POST /register", authHandler.RegisterUser)
	http.HandleFunc("POST /login", authHandler.LoginUser)

	http.HandleFunc("GET /messages", h.GetAll)
	http.HandleFunc("POST /messages", h.Create)
	http.HandleFunc("PUT /messages/{id}", h.Update)
	http.HandleFunc("DELETE /messages/{id}", h.Delete)

	// 3. Запуск веб-сервера
	log.Println("Сервер успешно запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка при работе сервера: %v", err)
	}
}
