package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Init() {
    clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
    println(os.Getenv("CLERK_SECRET_KEY"))
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }

        // Verify JWT with Clerk
        claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
            Token: token,
        })
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Add user ID to context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper to get user ID from context
func GetUserID(r *http.Request) string {
    userID, _ := r.Context().Value(UserIDKey).(string)
    return userID
}