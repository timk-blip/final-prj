package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	pass := os.Getenv("TODO_PASSWORD")
	var input struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if input.Password != pass {
		writeError(w, fmt.Errorf("password mismatch"), http.StatusUnauthorized)
		return
	}
	sum := sha256.Sum256([]byte(input.Password))
	hash := hex.EncodeToString(sum[:])
	claims := jwt.MapClaims{"password_hash": hash}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(pass))
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(pass), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		sum := sha256.Sum256([]byte(pass))
		if claims["password_hash"] != hex.EncodeToString(sum[:]) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
