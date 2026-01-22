package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SerzhLimon/GopherMartSolo/pkg/helpers"
	"github.com/SerzhLimon/GopherMartSolo/pkg/user"
)

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	var userBalance user.Balance
	
	row := h.Storage.DBStorage.InsertWithReturning(queryGetBalance, userID)

	userBalance, err := helpers.GetUserBalance(row)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userBalance)
}
