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
	router.Post("/", h.CreateRecord)
	router.Get("/", h.GetRecords)
	router.Get("/{id}", h.GetRecordByID)
	router.Put("/{id}", h.UpdateRecordByID)
	router.Delete("/{id}", h.DeleteRecordByID)
}

func loadAuthRoutes(router chi.Router, h *handler.Handler) {
	router.Post("/", h.CreateUser)
	router.Get("/", h.ListUsers)

	router.Post("/login", h.LoginUser)
	router.Get("/logout/{id}", h.LogoutUser)

	router.Route("/tokens", func(r chi.Router) {
		r.Post("/renew", h.RenewAccessToken)
		r.Post("/revoke/{id}", h.RevokeSession)
	})

	router.Delete("/{id}", h.DeleteUser)
}
