package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SerzhLimon/GopherMartSolo/pkg/helpers"
	"github.com/SerzhLimon/GopherMartSolo/pkg/user"
)

func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	row := h.Storage.DBStorage.InsertWithReturning(queryInsertWithReturning, userID)
	userBalance, err := helpers.GetUserBalance(row)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	balance := roundTo4Decimals(userBalance.Current - req.Sum)
	if userBalance.Current <= 0 || balance < 0 {
		http.Error(w, "Internal server error", http.StatusPaymentRequired)
		return
	}

	userBalance.Withdrawn += req.Sum
	userBalance.Current = balance

	
	err = h.Storage.DBStorage.Insert(queryUpdateUserBalance, userID, userBalance.Current, userBalance.Withdrawn, time.Now())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	err = h.Storage.DBStorage.Insert(queryInsertBalanceOperations, userID, req.Sum, req.Order, time.Now())
	if err != nil {
		http.Error(w, "Failed to save withdrawal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func roundTo4Decimals(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}
