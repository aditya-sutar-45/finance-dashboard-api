package app

import (
	"net/http"

	"github.com/aditya-sutar-45/finance-dashboard-api/handler"
	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes(db *database.Queries) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := &handler.Handler{
		DB: db,
	}

	router.Route("/records", func(r chi.Router) {
		loadRecordRoutes(r, h)
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
