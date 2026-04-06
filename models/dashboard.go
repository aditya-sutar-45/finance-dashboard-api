package models

import (
	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
)

type DashboardSummaryRow struct {
	TotalIncome  string `json:"total_income"`
	TotalExpense string `json:"total_expense"`
}

func DatabaseSummaryToDashboardSummary(dbSummary database.GetDashboardSummaryRow) DashboardSummaryRow {
	return DashboardSummaryRow{
		TotalIncome:  dbSummary.TotalIncome,
		TotalExpense: dbSummary.TotalExpense,
	}
}

type CategoryAnalysisRow struct {
	Category string
	Total    int64
}

func DatabaseCategoryAnalysisToCategoryAnalysis(dbCategory database.GetCategoryAnalysisRow) CategoryAnalysisRow {
	return CategoryAnalysisRow{
		Category: dbCategory.Category,
		Total:    dbCategory.Total,
	}
}

func DatabaseCategoryAnalysisToCategoryAnalysisRows(dbCategory []database.GetCategoryAnalysisRow) []CategoryAnalysisRow {
	result := []CategoryAnalysisRow{}
	for _, d := range dbCategory {
		result = append(result, DatabaseCategoryAnalysisToCategoryAnalysis(d))
	}

	return result
}

type TrendsRow struct {
	Month   string `json:"month"`
	Income  int64  `json:"income"`
	Expense int64  `json:"expense"`
}

func DatabaseTrendsRowToTrendsRow(dbRow database.GetTrendsRow) TrendsRow {
	return TrendsRow{
		Month:   dbRow.Month,
		Income:  dbRow.Income,
		Expense: dbRow.Expense,
	}
}

func DatabaseTrendsRowsToTrendsRows(dbRows []database.GetTrendsRow) []TrendsRow {
	result := []TrendsRow{}
	for _, r := range dbRows {
		result = append(result, DatabaseTrendsRowToTrendsRow(r))
	}

	return result
}
