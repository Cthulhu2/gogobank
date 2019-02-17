package handler

import (
	"net/http"
	"strconv"

	"app/model"
)

func TransactionPOST(w http.ResponseWriter, r *http.Request) {
	balance, err := strconv.ParseInt(r.PostFormValue("balance"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Incorrect balance value"))
		return //
	}
	acc := model.AccountNew(balance)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(acc.ID)))
}
