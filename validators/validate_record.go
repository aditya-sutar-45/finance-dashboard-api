// Package validators
package validators

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/models"
)

func ValidateCreateRecord(params *models.CreateRecordParameters, DB *database.Queries, ctx context.Context) error {
	if params.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Amount validation
	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}

	amountFloat, err := strconv.ParseFloat(params.Amount, 64)
	if err != nil || amountFloat <= 0 {
		return fmt.Errorf("amount must be a valid positive number")
	}

	// Type validation
	if params.Type != "income" && params.Type != "expense" {
		return fmt.Errorf("type must be 'income' or 'expense'")
	}

	// Category validation
	if strings.TrimSpace(params.Category) == "" {
		return fmt.Errorf("category is required")
	}

	// Date validation
	if params.Date.IsZero() {
		return fmt.Errorf("date is required")
	}

	if params.Date.After(time.Now()) {
		return fmt.Errorf("date cannot be in the future")
	}

	params.Note = strings.TrimSpace(params.Note)

	return nil
}

func ValidateGetRecords(
	typeParam string,
	categoryParam string,
	startDate time.Time,
	endDate time.Time,
	startDateRaw string,
	endDateRaw string,
) error {
	// Type validation
	if typeParam != "" && typeParam != "income" && typeParam != "expense" {
		return fmt.Errorf("type must be 'income' or 'expense'")
	}

	// Date validation

	// If only one date is provided then allowed (depends on your design)
	// but if both are present then validate range
	if startDateRaw != "" && endDateRaw != "" {
		if startDate.After(endDate) {
			return fmt.Errorf("start_date cannot be after end_date")
		}
	}

	// Prevent future filtering
	if !startDate.IsZero() && startDate.After(time.Now()) {
		return fmt.Errorf("start_date cannot be in the future")
	}

	if !endDate.IsZero() && endDate.After(time.Now()) {
		return fmt.Errorf("end_date cannot be in the future")
	}

	return nil
}

func ValidatgeUpdateRecord(params models.UpdateRecordParameters) error {
	if params.Amount == nil &&
		params.Type == nil &&
		params.Category == nil &&
		params.Note == nil &&
		params.Date == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Amount validation
	if params.Amount != nil {
		amountStr := strings.TrimSpace(*params.Amount)
		if amountStr == "" {
			return fmt.Errorf("amount cannot be empty")
		}

		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amountFloat <= 0 {
			return fmt.Errorf("amount must be a valid positive number")
		}
	}

	// Type validation
	if params.Type != nil {
		t := strings.TrimSpace(*params.Type)
		if t != "income" && t != "expense" {
			return fmt.Errorf("type must be 'income' or 'expense'")
		}
	}

	// Category validation
	if params.Category != nil {
		c := strings.TrimSpace(*params.Category)
		if c == "" {
			return fmt.Errorf("category cannot be empty")
		}
	}

	// Note normalization
	if params.Note != nil {
		trimmed := strings.TrimSpace(*params.Note)
		*params.Note = trimmed
	}

	// Date validation
	if params.Date != nil {
		if params.Date.IsZero() {
			return fmt.Errorf("date cannot be zero")
		}
		if params.Date.After(time.Now()) {
			return fmt.Errorf("date cannot be in the future")
		}
	}

	return nil
}
