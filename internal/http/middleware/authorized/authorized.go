package authorized

// // func Authorized(){

// // }

// import (
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v5"
// 	// "github.com/joho/godotenv"
// )

// func Authorized(next http.Handler) http.Handler {
// 	// Load .env once (optional: you can move this to your main.go)
	
// 	secretKey := []byte(os.Getenv("JWT_SECRET"))

	

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
// 			return
// 		}

// 		// Expected format: "Bearer <token>"
// 		parts := strings.SplitN(authHeader, " ", 2)
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			http.Error(w, `Expected format: "Bearer <token>"`, http.StatusUnauthorized)
// 			return
// 		}

// 		tokenStr := parts[1]

// 		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		slog.Info("this is token","token",token)

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 			return
// 		}

// 		// ✅ Token valid — proceed to the next handler
// 		next.ServeHTTP(w, r)
// 	})
// }

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userCtxKey ctxKey = "userClaims"

type UserClaims struct {
	ID string `json:"ID"`
	// you can also embed RegisteredClaims for standard fields
	jwt.RegisteredClaims
}

func Authorized(next http.Handler) http.Handler {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `expected "Bearer <token>"`, http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// now you have typed claims
		// claims.ID, claims.Subject, claims.ExpiresAt, claims.IssuedAt, etc.
		// Convert RegisteredClaims time pointers to time.Time if non-nil:
		var exp time.Time
		if claims.ExpiresAt != nil {
			exp = claims.ExpiresAt.Time
		}
		var iat time.Time
		if claims.IssuedAt != nil {
			iat = claims.IssuedAt.Time
		}

		user := map[string]interface{}{
			"ID":    claims.ID,
			"sub":   claims.Subject,
			"exp":   exp,
			"iat":   iat,
			"scope": claims.Audience, // example with other standard fields
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (map[string]interface{}, error) {
	v := r.Context().Value(userCtxKey)
	if v == nil {
		return nil, errors.New("no user in context")
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("invalid user in context")
}
