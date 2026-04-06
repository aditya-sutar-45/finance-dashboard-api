package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aditya-sutar-45/finance-dashboard-api/token"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
)

type authKey struct{}

func GetAuthMiddlwareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read handler
			// verifiy token
			userClaims, err := veriftyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("error verifying the token: %v", err))
				return
			}
			// pass the payload down
			ctx := context.WithValue(r.Context(), authKey{}, userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...token.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userClaims, ok := r.Context().Value(authKey{}).(*token.UserClaims)
			if !ok || userClaims == nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "missing auth context")
				return
			}

			userLevel, ok := token.RoleHierarchy[userClaims.Role]
			if !ok {
				utils.RespondWithError(w, http.StatusForbidden, "invalid role")
				return
			}

			for _, role := range roles {
				requiredLevel, ok := token.RoleHierarchy[role]
				if !ok {
					continue
				}

				if userLevel >= requiredLevel {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.RespondWithError(w, http.StatusForbidden, "forbidden")
		})
	}
}

func veriftyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("auth header is missing")
	}
	// Bearer <token>
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid auth header")
	}

	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
