package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/ybotet/pz6_csrf_xss/gen/proto/auth"
)

type AuthMiddleware struct {
    authClient authpb.AuthServiceClient
}

func NewAuthMiddleware(client authpb.AuthServiceClient) *AuthMiddleware {
    return &AuthMiddleware{
        authClient: client,
    }
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            log.Printf("Токен не предоставлен")
            http.Error(w, "Токен не предоставлен", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            log.Printf("Неверный формат токена")
            http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
            return
        }
        token := parts[1]

        // Crear contexto con deadline
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        log.Printf("Вызов gRPC Verify с дедлайном 2с")

        // Llamar a Auth service via gRPC
        resp, err := m.authClient.Verify(ctx, &authpb.VerifyRequest{
            Token: token,
        })

        if err != nil {
            // Mapear errores gRPC a HTTP
            st, ok := status.FromError(err)
            if !ok {
                log.Printf("Внутренняя ошибка сервера: %v", err)
                http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
                return
            }

            log.Printf("Ошибка gRPC: %v", st.Message())
            
            switch st.Code() {
            case codes.Unauthenticated:
                http.Error(w, "Недействительный токен", http.StatusUnauthorized)
            case codes.DeadlineExceeded:
                http.Error(w, "Таймаут при проверке аутентификации", http.StatusGatewayTimeout)
            case codes.Unavailable:
                http.Error(w, "Сервис аутентификации недоступен", http.StatusServiceUnavailable)
            default:
                http.Error(w, "Ошибка аутентификации", http.StatusInternalServerError)
            }
            return
        }

        if !resp.Valid {
            log.Printf("Недействительный токен (невалидный ответ)")
            http.Error(w, "Недействительный токен", http.StatusUnauthorized)
            return
        }

        log.Printf("Токен действителен для пользователя: %s", resp.Subject)
        
        // Añadir subject al contexto
        ctx = context.WithValue(r.Context(), "user", resp.Subject)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}