package handler

import (
	"database/sql"
	"encoding/json"
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
	"github.com/google/uuid"
)

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID   string    `json:"user_id"`
		Amount   string    `json:"amount"`
		Type     string    `json:"type"`
		Category string    `json:"category"`
		Note     string    `json:"note"`
		Date     time.Time `json:"date"`
	}

	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	uID, err := uuid.Parse(params.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id:  %v", err))
		return
	}
	// check if user exists in the database
	_, err = h.DB.GetUserByID(r.Context(), uID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	// Amount validation
	if params.Amount == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "amount is required")
		return
	}

	amountFloat, err := strconv.ParseFloat(params.Amount, 64)
	if err != nil || amountFloat <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "amount must be a valid positive number")
		return
	}

	// Type validation
	if params.Type != "income" && params.Type != "expense" {
		utils.RespondWithError(w, http.StatusBadRequest, "type must be 'income' or 'expense'")
		return
	}

	// Category validation
	if strings.TrimSpace(params.Category) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "category is required")
		return
	}

	// Date validation
	if params.Date.IsZero() {
		utils.RespondWithError(w, http.StatusBadRequest, "date is required")
		return
	}

	if params.Date.After(time.Now()) {
		utils.RespondWithError(w, http.StatusBadRequest, "date cannot be in the future")
		return
	}

	note := strings.TrimSpace(params.Note)

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
			String: note,
			Valid:  note != "",
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
	typeParam := strings.TrimSpace(query.Get("type"))
	categoryParam := strings.TrimSpace(query.Get("category"))
	startDateParam := query.Get("start_date")
	endDateParam := query.Get("end_date")

	var startDate time.Time
	var endDate time.Time

	if startDateParam != "" {
		startDate, err = time.Parse("2006-01-02", startDateParam)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid start_date format (YYYY-MM-DD)")
			return
		}
	}

	if endDateParam != "" {
		endDate, err = time.Parse("2006-01-02", endDateParam)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid end_date format (YYYY-MM-DD)")
			return
		}
	}

	switch uRole {
	case "analyst", "admin":
		h.getAllRecords(w, r, query, typeParam, categoryParam, startDate, endDate)
	case "viewer":
		h.getUserRecords(w, r, uID, typeParam, categoryParam, startDate, endDate)
	}
}

func (h *Handler) getUserRecords(w http.ResponseWriter, r *http.Request, userID uuid.UUID, typeParam string, categoryParam string, startDate time.Time, endDate time.Time) {
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
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "could not fetch records from the database")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordsToRecords(records))
}

func (h *Handler) getAllRecords(w http.ResponseWriter, r *http.Request, query url.Values, typeParam string, categoryParam string, startDate time.Time, endDate time.Time) {
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
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("error fetching records from database: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseRecordsToRecords(records))
}

func (h *Handler) GetRecordByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateRecordByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteRecordByID(w http.ResponseWriter, r *http.Request) {
}
