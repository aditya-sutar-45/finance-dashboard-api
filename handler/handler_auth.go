package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/models"
	"github.com/aditya-sutar-45/finance-dashboard-api/token"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request json: %v", err))
		return
	}

	user, err := h.DB.GetUser(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("user not found: %v", err))
		return
	}

	err = utils.CheckPassword(params.Password, user.PasswordHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "error checking password")
		return
	}

	// create token
	accessToken, accessClaims, err := h.TokenMaker.CreateToken(user.ID, user.Email, token.Role(user.Role), 2*time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating the token: %v", err))
		return
	}

	// refresh token
	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(user.ID, user.Email, token.Role(user.Role), 24*time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating the token: %v", err))
		return
	}

	session, err := h.DB.CreateSession(r.Context(), database.CreateSessionParams{
		ID:           refreshClaims.ID.String(),
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt: sql.NullTime{
			Time:  refreshClaims.ExpiresAt.Time,
			Valid: true,
		},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to create session: %v", err))
		return
	}

	u := models.DatabaseUserToUser(user)

	res := models.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		User:                  u,
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	err := h.DB.DeleteSession(r.Context(), claims.ID.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "session does not exist")
		return
	}

	type Res struct {
		Message string `json:"message"`
	}

	utils.RespondWithJSON(w, http.StatusOK, Res{
		Message: "logout success",
	})
}

func (h *Handler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RefreshToken string `json:"refresh_token"`
	}
	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request: %v", err))
		return
	}

	refreshClaims, err := h.TokenMaker.VerifyToken(params.RefreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "error verifying token")
		return
	}

	session, err := h.DB.GetSession(r.Context(), refreshClaims.ID.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "error getting the session")
		return
	}

	if session.IsRevoked {
		utils.RespondWithError(w, http.StatusUnauthorized, "session revoked")
		return
	}

	if session.UserEmail != refreshClaims.Email {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid session")
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.Role, 15*time.Minute)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	res := models.RenewAcessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
	}

	utils.RespondWithJSON(w, http.StatusOK, res)
}

func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	err := h.DB.RevokeSession(r.Context(), claims.RegisteredClaims.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "session not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, struct{}{})
}
