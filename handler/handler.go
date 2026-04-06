// Package handler
package handler

import (
	"fmt"
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/token"
	"github.com/google/uuid"
)

type Handler struct {
	DB         *database.Queries
	TokenMaker *token.JWTMaker
}

func NewHandler(db *database.Queries, secretKey string) *Handler {
	return &Handler{
		DB:         db,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func getUserIDFromClaims(r *http.Request) (uuid.UUID, token.Role, error) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	if claims == nil {
		return uuid.Nil, "", fmt.Errorf("claims not found in context")
	}

	return claims.ID, claims.Role, nil
}
