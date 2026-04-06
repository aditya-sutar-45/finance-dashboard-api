package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/models"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
	"github.com/aditya-sutar-45/finance-dashboard-api/validators"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	var params models.CreateRecordParameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}
	if err := validators.ValidateCreateRecord(&params, h.DB, r.Context()); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("validation error: %v", err))
		return
	}

	uID, err := uuid.Parse(params.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}
	// check if user exists in the database
	_, err = h.DB.GetUserByID(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	createdByID, _, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting user id: %v", err))
		return
	}

	record, err := h.DB.CreateRecord(r.Context(), database.CreateRecordParams{
		ID:     uuid.New(),
		UserID: uID,
		CreatedBy: uuid.NullUUID{
			UUID:  createdByID,
			Valid: true,
		},
		Category: params.Category,
		Amount:   params.Amount,
		Type:     params.Type,
		Note: sql.NullString{
			String: params.Note,
			Valid:  params.Note != "",
		},
		Date: params.Date,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error creating record: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseRecordToRecord(record))
}

func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) {
	uID, uRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting user id: %v", err))
		return
	}

	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")
	typeParam := strings.TrimSpace(query.Get("type"))
	categoryParam := strings.TrimSpace(query.Get("category"))
	startDateParam := query.Get("start_date")
	endDateParam := query.Get("end_date")

	var startDate time.Time
	var endDate time.Time

	if err := validators.ValidateGetRecords(typeParam, categoryParam, startDate, endDate, startDateParam, endDateParam); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error validating query params: %v", err))
		return
	}

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	switch uRole {
	case "analyst", "admin":
		h.getAllRecords(w, r, int32(page), int32(limit), int32(offset), query, typeParam, categoryParam, startDate, endDate)
	case "viewer":
		h.getUserRecords(w, r, int32(page), int32(limit), int32(offset), uID, typeParam, categoryParam, startDate, endDate)
	}
}

func (h *Handler) getUserRecords(w http.ResponseWriter, r *http.Request, page int32, limit int32, offset int32, userID uuid.UUID, typeParam string, categoryParam string, startDate time.Time, endDate time.Time) {
	records, err := h.DB.ListRecords(r.Context(), database.ListRecordsParams{
		UserID: userID,
		Type: sql.NullString{
			String: typeParam,
			Valid:  typeParam != "",
		},
		Category: sql.NullString{
			String: categoryParam,
			Valid:  categoryParam != "",
		},
		StartDate: sql.NullTime{
			Time:  startDate,
			Valid: !startDate.IsZero(),
		},
		EndDate: sql.NullTime{
			Time:  endDate,
			Valid: !endDate.IsZero(),
		},
		PageOffset: offset,
		PageLimit:  limit,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "could not fetch records from the database")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.GetRecordsResponse{
		Page:    page,
		Limit:   limit,
		Records: models.DatabaseRecordsToRecords(records),
	})
}

func (h *Handler) getAllRecords(w http.ResponseWriter, r *http.Request, page int32, limit int32, offset int32, query url.Values, typeParam string, categoryParam string, startDate time.Time, endDate time.Time) {
	var filterUserID uuid.UUID
	filterUserIDString := strings.TrimSpace(query.Get("user_id"))

	if filterUserIDString != "" {
		// check uuid
		log.Println(filterUserIDString)

		parsedID, err := uuid.Parse(filterUserIDString)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid user id format")
			return
		}

		filterUserID = parsedID

		// check if the user exists
		_, err = h.DB.GetUserByID(r.Context(), filterUserID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}
	}

	records, err := h.DB.ListAllRecords(r.Context(), database.ListAllRecordsParams{
		UserID: uuid.NullUUID{
			UUID:  filterUserID,
			Valid: filterUserIDString != "",
		},
		Type: sql.NullString{
			String: typeParam,
			Valid:  typeParam != "",
		},
		Category: sql.NullString{
			String: categoryParam,
			Valid:  categoryParam != "",
		},
		StartDate: sql.NullTime{
			Time:  startDate,
			Valid: !startDate.IsZero(),
		},
		EndDate: sql.NullTime{
			Time:  endDate,
			Valid: !endDate.IsZero(),
		},
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("error fetching records from database: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.GetRecordsResponse{
		Page:    page,
		Limit:   limit,
		Records: models.DatabaseRecordsToRecords(records),
	})
}

func (h *Handler) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	recordID, err := uuid.Parse(idString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid record id")
		return
	}

	userID, userRole, err := getUserIDFromClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "no claims found")
		return
	}

	switch userRole {
	case "admin", "analyst":
		record, err := h.DB.GetRecordByID(r.Context(), recordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusNotFound, "record not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch record")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordToRecord(record))
	case "viewer":
		record, err := h.DB.GetViewerRecordByID(r.Context(), database.GetViewerRecordByIDParams{
			ID:     recordID,
			UserID: userID,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusNotFound, "record not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch record")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordToRecord(record))
	default:
		utils.RespondWithError(w, http.StatusForbidden, "invalid role")
		return
	}
}

func (h *Handler) UpdateRecordByID(w http.ResponseWriter, r *http.Request) {
	var params models.UpdateRecordParameters

	idString := chi.URLParam(r, "id")
	recordID, err := uuid.Parse(idString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid record id")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validators.ValidatgeUpdateRecord(params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error validating body: %v", err))
		return
	}

	record, err := h.DB.PatchRecordByID(r.Context(), database.PatchRecordByIDParams{
		ID: recordID,
		Amount: sql.NullString{
			String: utils.GetString(params.Amount),
			Valid:  params.Amount != nil,
		},
		Type: sql.NullString{
			String: utils.GetString(params.Type),
			Valid:  params.Type != nil,
		},
		Category: sql.NullString{
			String: utils.GetString(params.Category),
			Valid:  params.Category != nil,
		},
		Note: sql.NullString{
			String: utils.GetString(params.Note),
			Valid:  params.Note != nil,
		},
		Date: sql.NullTime{
			Time:  utils.GetTime(params.Date),
			Valid: params.Date != nil,
		},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to update record")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordToRecord(record))
}

func (h *Handler) DeleteRecordByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	recordID, err := uuid.Parse(idString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid record id")
		return
	}

	err = h.DB.DeleteRecordByID(r.Context(), recordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "record not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting record")
		return
	}

	type res struct {
		Message string `json:"message"`
	}

	utils.RespondWithJSON(w, http.StatusOK, res{
		Message: "successfully deleted",
	})
}
