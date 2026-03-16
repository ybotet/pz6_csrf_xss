package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/ybotet/pz6_csrf_xss/services/auth/internal/auth"
	grpcserver "github.com/ybotet/pz6_csrf_xss/services/auth/internal/grpc"
	"github.com/ybotet/pz6_csrf_xss/shared/logger"
	"github.com/ybotet/pz6_csrf_xss/shared/middleware"
)

func main() {

    log := logger.New(logger.Config{
        ServiceName: "auth",
        Environment: "development",
        LogLevel:    "debug",
        JSONFormat:  true,
    })
    
    // Router
    r := mux.NewRouter()
    
    // Middlewares в ПРАВИЛЬНОМ порядке
    r.Use(middleware.RequestID)
    r.Use(middleware.Logging(log))

    port := os.Getenv("AUTH_GRPC_PORT")
    if port == "" {
        port = "50051"
        log.Printf("AUTH_GRPC_PORT не настроен, используется порт по умолчанию: %s", port)
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "your-secret-key-change-in-production"
        log.Printf("JWT_SECRET не настроен, используется ключ по умолчанию (НЕ ИСПОЛЬЗОВАТЬ В ПРОДАКШЕНЕ)")
    }

    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("Ошибка прослушивания порта %s: %v", port, err)
    }

    s := grpc.NewServer()
    authService := auth.NewService(jwtSecret)
    grpcserver.Register(s, authService)

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        log.Println("Выключение gRPC сервера...")
        s.GracefulStop()
    }()

    log.Printf("gRPC сервер Auth слушает порт %s", port)
    
    // Генерация тестового токена
    testToken, _ := authService.GenerateToken("тестовый_пользователь", 24*time.Hour)
    log.Printf("Тестовый токен (действителен 24ч): %s", testToken)
    
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Ошибка при обслуживании: %v", err)
    }
}