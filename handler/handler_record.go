package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromClaims(r)
	fmt.Println(userID)
}

func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetRecordByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateRecordByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteRecordByID(w http.ResponseWriter, r *http.Request) {
}
