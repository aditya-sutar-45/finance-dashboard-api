package models

import (
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/google/uuid"
)

type Record struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedBy uuid.UUID `json:"created_by"`
	Amount    string    `json:"amount"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetRecordsResponse struct {
	Page    int32    `json:"page"`
	Limit   int32    `json:"limit"`
	Records []Record `json:"records"`
}

type CreateRecordParameters struct {
	UserID   string    `json:"user_id"`
	Amount   string    `json:"amount"`
	Type     string    `json:"type"`
	Category string    `json:"category"`
	Note     string    `json:"note"`
	Date     time.Time `json:"date"`
}

type UpdateRecordParameters struct {
	Amount   *string    `json:"amount"`
	Type     *string    `json:"type"`
	Category *string    `json:"category"`
	Note     *string    `json:"note"`
	Date     *time.Time `json:"date"`
}

func DatabaseRecordToRecord(dbRecord database.Record) Record {
	var createdAt time.Time
	var updatedAt time.Time
	var createdBy uuid.UUID
	note := ""

	if dbRecord.CreatedBy.Valid {
		createdBy = dbRecord.CreatedBy.UUID
	}

	if dbRecord.UpdatedAt.Valid {
		updatedAt = dbRecord.UpdatedAt.Time
	}

	if dbRecord.CreatedAt.Valid {
		createdAt = dbRecord.CreatedAt.Time
	}

	if dbRecord.Note.Valid {
		note = dbRecord.Note.String
	}

	return Record{
		ID:        dbRecord.ID,
		UserID:    dbRecord.UserID,
		CreatedBy: createdBy,
		Amount:    dbRecord.Amount,
		Type:      dbRecord.Type,
		Category:  dbRecord.Category,
		Note:      note,
		Date:      dbRecord.Date,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func DatabaseRecordsToRecords(dbRecords []database.Record) []Record {
	records := []Record{}
	for _, r := range dbRecords {
		records = append(records, DatabaseRecordToRecord(r))
	}

	return records
}
