package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/0xzcarlet/sregep/backend/internal/domain"
)

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTransaction(w, r)
	case http.MethodGet:
		h.listTransactions(w, r)
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input domain.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.finance.CreateTransaction(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	result, err := h.finance.ListTransactions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	userID := r.URL.Query().Get("user_id")
	result, err := h.finance.Summary(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result})
}
