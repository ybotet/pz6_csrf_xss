package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ybotet/pz3_logr/services/tasks/internal/clients"
	"github.com/ybotet/pz3_logr/services/tasks/internal/handlers"

	// Алиас для внутреннего middleware (tasks)
	internalMiddleware "github.com/ybotet/pz3_logr/services/tasks/internal/middleware"

	// Общие пакеты
	"github.com/ybotet/pz3_logr/shared/logger"
	"github.com/ybotet/pz3_logr/shared/middleware" // общий middleware
)

func main() {
    tasksPort := os.Getenv("TASKS_PORT")
    if tasksPort == "" {
        tasksPort = "8082"
    }

    authAddr := os.Getenv("AUTH_GRPC_ADDR")
    if authAddr == "" {
        authAddr = "localhost:50051"
    }

    log := logger.New(logger.Config{
        ServiceName: "tasks",
        Environment: "development",
        LogLevel:    "debug",
        JSONFormat:  true,
    })
    
    // Роутер
    r := mux.NewRouter()
    
    // Middleware из SHARED (RequestID и Logging)
    r.Use(middleware.RequestID)      // ✅ Из общего middleware
    r.Use(middleware.Logging(log))   // ✅ Из общего middleware

    // Подключение к Auth service
    authClient, err := clients.NewAuthClient(authAddr)
    if err != nil {
        log.Fatalf("Ошибка подключения к Auth service: %v", err)
    }
    defer authClient.Close()

    // Создание middleware и обработчиков
    // Используем internalMiddleware (с алиасом) для NewAuthMiddleware
    authMiddleware := internalMiddleware.NewAuthMiddleware(authClient.GetClient())  // ✅ Изменено
    taskHandler := handlers.NewTaskHandler()

    // Настройка маршрутов
    r.HandleFunc("/tasks", authMiddleware.Authenticate(taskHandler.GetTasks)).Methods("GET")
    r.HandleFunc("/tasks", authMiddleware.Authenticate(taskHandler.CreateTask)).Methods("POST")

    log.Printf("Сервер Tasks слушает порт %s", tasksPort)
    log.Fatal(http.ListenAndServe(":"+tasksPort, r))
}