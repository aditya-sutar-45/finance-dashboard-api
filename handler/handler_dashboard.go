package handler

import (
	"fmt"
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/models"
	"github.com/aditya-sutar-45/finance-dashboard-api/token"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
	"github.com/google/uuid"
)

func (h *Handler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	userID, userRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	query := r.URL.Query()
	userIDParam := query.Get("user_id")

	var uID uuid.NullUUID
	if userRole == token.RoleViewer {
		uID = uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		}
	} else {
		if userIDParam != "" {
			parsedID, err := uuid.Parse(userIDParam)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid user_id")
				return
			}

			uID = uuid.NullUUID{
				UUID:  parsedID,
				Valid: true,
			}
		} else {
			uID = uuid.NullUUID{Valid: false}
		}
	}

	data, err := h.DB.GetDashboardSummary(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch summary:\n %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseSummaryToDashboardSummary(data))
}

func (h *Handler) GetCategoryAnalysis(w http.ResponseWriter, r *http.Request) {
	userID, userRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	query := r.URL.Query()
	userIDParam := query.Get("user_id")

	var uID uuid.NullUUID
	if userRole == token.RoleViewer {
		uID = uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		}
	} else {
		if userIDParam != "" {
			parsedID, err := uuid.Parse(userIDParam)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid user_id")
				return
			}

			uID = uuid.NullUUID{
				UUID:  parsedID,
				Valid: true,
			}
		} else {
			uID = uuid.NullUUID{Valid: false}
		}
	}

	data, err := h.DB.GetCategoryAnalysis(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch summary")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseCategoryAnalysisToCategoryAnalysisRows(data))
}

func (h *Handler) GetTrends(w http.ResponseWriter, r *http.Request) {
	userID, userRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	query := r.URL.Query()
	userIDParam := query.Get("user_id")

	var uID uuid.NullUUID
	if userRole == token.RoleViewer {
		uID = uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		}
	} else {
		if userIDParam != "" {
			parsedID, err := uuid.Parse(userIDParam)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid user_id")
				return
			}

			uID = uuid.NullUUID{
				UUID:  parsedID,
				Valid: true,
			}
		} else {
			uID = uuid.NullUUID{Valid: false}
		}
	}

	data, err := h.DB.GetTrends(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch summary")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseTrendsRowsToTrendsRows(data))
}

func (h *Handler) GetRecent(w http.ResponseWriter, r *http.Request) {
	userID, userRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	query := r.URL.Query()
	userIDParam := query.Get("user_id")

	var uID uuid.NullUUID
	if userRole == token.RoleViewer {
		uID = uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		}
	} else {
		if userIDParam != "" {
			parsedID, err := uuid.Parse(userIDParam)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid user_id")
				return
			}

			uID = uuid.NullUUID{
				UUID:  parsedID,
				Valid: true,
			}
		} else {
			uID = uuid.NullUUID{Valid: false}
		}
	}

	data, err := h.DB.GetRecent(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch summary")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordsToRecords(data))
}
