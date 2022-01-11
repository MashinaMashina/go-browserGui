package repository

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
	"strings"
)

type AuthError struct {
	Message string `json:"message"`
}

func GetTokenSecretKey() []byte {
	return []byte("my$secret$key")
}

func GetTokenSecretKeyByToken(t *jwt.Token) (interface{}, error){
	return GetTokenSecretKey(), nil
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		token := &Token{}
		token.Data = make(map[string]string)

		if header != "" {
			header = strings.Replace(header, "Bearer ", "", 1)

			_, err := jwt.ParseWithClaims(header, token, GetTokenSecretKeyByToken)

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(AuthError{Message: err.Error()})
				return
			}
		}

		repoCtx := SetRepository(r.Context(), NewRepository(w, token))

		next.ServeHTTP(w, r.WithContext(repoCtx))
	})
}
