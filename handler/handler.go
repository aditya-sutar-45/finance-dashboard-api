// Package handler
package handler

import (
	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/token"
)

type Handler struct {
	DB         *database.Queries
	tokenMaker *token.JWTMaker
}

func NewHandler(db *database.Queries, secretKey string) *Handler {
	return &Handler{
		DB:         db,
		tokenMaker: token.NewJWTMaker(secretKey),
	}
}
