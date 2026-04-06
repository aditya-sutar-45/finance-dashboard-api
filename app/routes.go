package app

import (
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/handler"
	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/aditya-sutar-45/finance-dashboard-api/token"
	"github.com/aditya-sutar-45/finance-dashboard-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes(db *database.Queries, secretKey string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, struct{}{})
	})
	router.Get("/err", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithError(w, 400, "Something went wrong!")
	})

	h := handler.NewHandler(db, secretKey)

	router.Route("/records", func(r chi.Router) {
		loadRecordRoutes(r, h)
	})

	router.Route("/users", func(r chi.Router) {
		loadAuthRoutes(r, h)
	})

	router.Route("/dashboard", func(r chi.Router) {
		loadDashboardRoutes(r, h)
	})

	return router
}

func loadRecordRoutes(router chi.Router, h *handler.Handler) {
	tokenMaker := h.TokenMaker

	router.Use(handler.GetAuthMiddlwareFunc(tokenMaker))

	router.With(handler.RequireRole(token.RoleAnalyst)).Get("/{id}", h.GetRecordByID)
	router.With(handler.RequireRole(token.RoleAnalyst)).Get("/", h.GetRecords)

	router.With(handler.RequireRole(token.RoleAdmin)).Post("/", h.CreateRecord)
	router.With(handler.RequireRole(token.RoleAdmin)).Patch("/{id}", h.UpdateRecordByID)
	router.With(handler.RequireRole(token.RoleAdmin)).Delete("/{id}", h.DeleteRecordByID)
	router.With(handler.RequireRole(token.RoleAdmin)).Get("/deleted", h.GetDeletedRecords)
	router.With(handler.RequireRole(token.RoleAdmin)).Delete("/{id}/h", h.HardDeleteRecordByID)
}

func loadAuthRoutes(router chi.Router, h *handler.Handler) {
	// Public
	router.Post("/login", h.LoginUser)

	//  Refresh token endpoint
	router.Post("/tokens/renew", h.RenewAccessToken)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(handler.GetAuthMiddlwareFunc(h.TokenMaker))
		r.With(handler.RequireRole(token.RoleAdmin)).Get("/deleted", h.GetDeletedUsers)

		r.Get("/logout", h.LogoutUser)

		r.Post("/tokens/revoke", h.RevokeSession)

		r.With(handler.RequireRole(token.RoleAdmin)).Post("/", h.CreateUser)
		r.With(handler.RequireRole(token.RoleAdmin)).Get("/", h.ListUsers)
		r.With(handler.RequireRole(token.RoleAdmin)).Delete("/{id}", h.DeleteUser)
		r.With(handler.RequireRole(token.RoleAdmin)).Delete("/{id}/h", h.HardDeleteUser)
	})
}

func loadDashboardRoutes(router chi.Router, h *handler.Handler) {
	tokenMaker := h.TokenMaker
	router.Use(handler.GetAuthMiddlwareFunc(tokenMaker))

	router.With(handler.RequireRole(token.RoleViewer)).Get("/summary", h.GetDashboardSummary)
	router.With(handler.RequireRole(token.RoleViewer)).Get("/categories", h.GetCategoryAnalysis)
	router.With(handler.RequireRole(token.RoleViewer)).Get("/trends", h.GetTrends)
	router.With(handler.RequireRole(token.RoleViewer)).Get("/recent", h.GetRecent)
}
