package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"app/model"
)

// GetAccRes
// Get Account response model.
// swagger:response getAccRes
type swagGetAccRes struct {
	// in: body
	Body GetAccRes
}

type GetAccRes struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg,omitempty"`
	Acc  model.Account `json:"acc,omitempty"`
}

// swagger:operation GET /account/{id} account noReq
// ---
// summary: The account with balance by ID.
// description: If the account doesn't exists NotExistsCode will be returned.
// parameters:
// - name: id
//   in: path
//   description: The account ID.
//   type: integer
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/getAccRes"
func AccountGET(w http.ResponseWriter, r *http.Request) {
	var res GetAccRes
	w.WriteHeader(http.StatusOK)

	sId := strings.TrimPrefix(r.URL.Path, "/v1/account/")
	id, err := strconv.ParseInt(sId, 10, 64)
	if err != nil {
		res = GetAccRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error()}
	} else {
		acc, err := model.AccountGet(id)
		if err != nil {
			switch err.(type) {
			default:
				res = GetAccRes{
					Code: http.StatusInternalServerError,
					Msg:  err.Error()}

			case *model.NotExistsError:
				res = GetAccRes{
					Code: model.NotExistsCode,
					Msg:  err.Error()}
			}
		} else {
			res = GetAccRes{Code: 0, Acc: *acc}
		}
	}

	json.NewEncoder(w).Encode(res)
}

// newAccReq
// Create Account request model.
// swagger:parameters newAccReq
type swagNewAccReq struct {
	// in: body
	Body NewAccReq
}

type NewAccReq struct {
	Balance int64 `json:"balance"`
}

// newAccRes
// Create Account response model.
// swagger:response newAccRes
type swagNewAccRes struct {
	// in: body
	Body NewAccRes
}

type NewAccRes struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg,omitempty"`
	Acc  model.Account `json:"acc,omitempty"`
}

// swagger:operation POST /account/ account newAccReq
// ---
// summary: Create an account with balance.
// responses:
//   "200":
//     "$ref": "#/responses/newAccRes"
func AccountPOST(w http.ResponseWriter, r *http.Request) {
	var req NewAccReq
	var res NewAccRes

	w.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &req)
	if err != nil {
		res = NewAccRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error()}
	} else {
		acc := model.AccountNew(req.Balance)
		res = NewAccRes{Code: 0, Acc: *acc}
	}

	json.NewEncoder(w).Encode(res)
}

// transReq
// Transfer request model.
// swagger:parameters transReq
type swagTransReq struct {
	// in: body
	Body TransReq
}

type TransReq struct {
	// The sender account ID.
	IdFrom int64 `json:"idFrom"`
	// The recipient account ID.
	IdTo int64 `json:"idTo"`
	// The transfer summ.
	Summ int64 `json:"summ"`
}

// transRes
// Transfer response model.
// swagger:response transRes
type swagTransRes struct {
	// in: body
	Body TransRes
}

type TransRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// swagger:operation POST /account/trans/ account transReq
// ---
// summary: Transfer money from an account to another one.
// description: If 'from/to' account doesn't exists NotExistsCode will be
//   returned. If 'from'-account has not enough money NotEnoughMoneyCode will
//   be returned.
// responses:
//   "200":
//     "$ref": "#/responses/transRes"
func TransferPOST(w http.ResponseWriter, r *http.Request) {
	var req TransReq
	var res TransRes

	w.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &req)
	if err != nil {
		res = TransRes{
			Code: http.StatusBadRequest,
			Msg:  err.Error()}
	} else {
		err = model.Transfer(req.IdFrom, req.IdTo, req.Summ)
		if err != nil {
			switch err.(type) {
			default:
				res = TransRes{
					Code: http.StatusInternalServerError,
					Msg:  err.Error()}

			case *model.NotEnoughMoneyError:
				res = TransRes{
					Code: model.NotEnoughMoneyCode,
					Msg:  err.Error()}

			case *model.NotExistsError:
				res = TransRes{
					Code: model.NotExistsCode,
					Msg:  err.Error()}
			}
		} else {
			res = TransRes{Code: 0}
		}
	}

	json.NewEncoder(w).Encode(res)
}
