package app

import (
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/handler"
	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
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

	return router
}

func loadRecordRoutes(router chi.Router, h *handler.Handler) {
	tokenMaker := h.TokenMaker

	router.Use(handler.GetAuthMiddlwareFunc(tokenMaker))

	router.Get("/", h.GetRecords)
	router.Get("/{id}", h.GetRecordByID)

	router.With(handler.RequireRole("analyst", "admin")).Post("/", h.CreateRecord)
	router.With(handler.RequireRole("analyst", "admin")).Put("/{id}", h.UpdateRecordByID)

	router.With(handler.RequireRole("admin")).Delete("/{id}", h.DeleteRecordByID)
}

func loadAuthRoutes(router chi.Router, h *handler.Handler) {
	// Public
	router.Post("/", h.CreateUser)
	router.Post("/login", h.LoginUser)

	//  Refresh token endpoint
	router.Post("/tokens/renew", h.RenewAccessToken)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(handler.GetAuthMiddlwareFunc(h.TokenMaker))

		r.Get("/logout", h.LogoutUser)

		r.Post("/tokens/revoke", h.RevokeSession)

		r.With(handler.RequireRole("admin")).Get("/", h.ListUsers)
		r.With(handler.RequireRole("admin")).Delete("/{id}", h.DeleteUser)
	})
}
