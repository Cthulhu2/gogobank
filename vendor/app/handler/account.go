package handler

import (
	"net/http"
	"strconv"
	"strings"

	"app/model"
)

func AccountGET(w http.ResponseWriter, r *http.Request) {
	sId := strings.TrimPrefix(r.URL.Path, "/account/")
	id, err := strconv.ParseInt(sId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Incorrect ID"))
		return //
	}
	acc, err := model.AccountGet(id)
	if err != nil {
		switch err.(type) {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		case *model.NotExistsError:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(acc.Balance)))
	}
}

func AccountPOST(w http.ResponseWriter, r *http.Request) {
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
