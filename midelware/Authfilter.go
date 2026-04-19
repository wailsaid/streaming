package midelware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

type contextKey string

const UserContextKey contextKey = "user"

// JwtFilter is a standard http middleware that validates Bearer JWT tokens.
func JwtFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log path
		log.Println("JwtFilter path", r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		stringToken := strings.TrimPrefix(authHeader, "Bearer ")

		// If no header, check the cookie
		if stringToken == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				stringToken = cookie.Value
			}
		}

		if stringToken == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil || !token.Valid {
			log.Println("JwtFilter token invalid")
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		log.Println("JwtFilter token valid")

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Check expiry
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				http.Error(w, `{"error":"token expired"}`, http.StatusUnauthorized)
				return
			}
		}
		log.Println("JwtFilter claims", claims)

		// Load user from DB
		var user models.User
		res := database.DB.First(&user, "email = ?", claims["sub"])
		if res.Error != nil || user.ID == 0 {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		log.Println("JwtFilter user", user)

		// Store user in context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
