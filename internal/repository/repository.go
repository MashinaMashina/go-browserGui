package repository

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
)

type Repository struct {
	writer http.ResponseWriter
	token *Token
}

func NewRepository(w http.ResponseWriter, t *Token) *Repository {
	return &Repository{
		writer: w,
		token:  t,
	}
}

func (r *Repository) Get(k string) (value string, exists bool) {
	value, exists = r.token.Data[k]

	return value, exists
}

func (r *Repository) Set(k, value string) {
	r.token.Data[k] = value

	r.Save()
}

func (r *Repository) Save() {
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, r.token)
	tokenString, err := tk.SignedString(GetTokenSecretKey())

	if err != nil {
		r.writer.WriteHeader(http.StatusForbidden)
		json.NewEncoder(r.writer).Encode(AuthError{Message: "Can not generate token"})
		return
	}

	r.writer.Header().Set("Authorization", "Bearer " + tokenString)
}

func SetRepository(ctx context.Context, repository *Repository) context.Context {
	return context.WithValue(ctx, "repository", repository)
}

func GetRepository(ctx context.Context) *Repository {
	repo := ctx.Value("repository")

	switch repo.(type) {
	case *Repository:
		return repo.(*Repository)
	default:
		return nil
	}
}