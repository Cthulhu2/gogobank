package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/handler"
	"app/model"
)

func TestIndexHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler.IndexGET)
	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: %v want %v", status, http.StatusOK)
	}

	expected := `Hello World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("Unexpected body: %v want %v", actual, expected)
	}
}

func AssertNewAccount(balance int64, t *testing.T) model.Account {
	reqBody, _ := json.Marshal(handler.NewAccReq{Balance: balance})
	req, _ := http.NewRequest("POST", "/v1/account/",
		bytes.NewReader(reqBody))

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler.AccountPOST)
	hf.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: %v want %v", status, http.StatusOK)
	}

	var res handler.NewAccRes
	body, _ := ioutil.ReadAll(recorder.Body)
	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Errorf("Unexpected json response: %s", body)
	}
	if res.Code != 0 {
		t.Errorf("Unexpected err code: %v, want %v.", res.Code, 0)
	}
	if res.Acc.ID <= 0 {
		t.Errorf("Unexpected ID: %v, want int greater 0.", res.Acc.ID)
	}
	if res.Acc.Balance != balance {
		t.Errorf("Unexpected Balance: %v, want %v.", res.Acc.Balance, balance)
	}
	return res.Acc
}

func HandleGetAccount(id int64, t *testing.T) handler.GetAccRes {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/account/%d", id), nil)

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler.AccountGET)
	hf.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: %v want %v", status, http.StatusOK)
	}

	var res handler.GetAccRes
	body, _ := ioutil.ReadAll(recorder.Body)
	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Errorf("Unexpected json response: %s", body)
	}
	return res
}

func AssertGetAccount(id int64, t *testing.T) model.Account {
	var res = HandleGetAccount(id, t)
	if res.Code != 0 {
		t.Errorf("Unexpected err code: got %v, want %v.", res.Code, 0)
	}
	if res.Acc.ID != id {
		t.Errorf("Unexpected ID: got %v, want %v.", res.Acc.ID, id)
	}
	return res.Acc
}

func TestAccount(t *testing.T) {
	model.Init()
	var acc1 model.Account
	var acc2 model.Account

	acc1 = AssertNewAccount(200, t)
	acc2 = AssertNewAccount(300, t)
	if acc1.ID != 1 {
		t.Errorf("Unexpected ID: %v, want %v.", acc1.ID, 1)
	}
	if acc2.ID != 2 {
		t.Errorf("Unexpected ID: %v, want %v.", acc2.ID, 2)
	}

	acc1 = AssertGetAccount(1, t)
	acc2 = AssertGetAccount(2, t)
	if acc1.Balance != 200 {
		t.Errorf("Unexpected Balance: %v, want %v.", acc1.Balance, 200)
	}
	if acc2.Balance != 300 {
		t.Errorf("Unexpected Balance: %v, want %v.", acc2.Balance, 300)
	}
}

func TestAccountFail(t *testing.T) {
	model.Init()
	var accRes3 handler.GetAccRes

	accRes3 = HandleGetAccount(3, t)
	if accRes3.Code != model.NotExistsCode {
		t.Errorf("Unexpected err code: %v, want %v.",
			accRes3.Code, model.NotExistsCode)
	}
}

func HandleTransfer(
	idFrom int64, idTo int64, summ int64,
	t *testing.T) handler.TransRes {

	reqBody, _ := json.Marshal(handler.TransReq{
		IdFrom: idFrom, IdTo: idTo, Summ: summ})

	req, _ := http.NewRequest("POST", "/v1/account/trans/",
		bytes.NewReader(reqBody))

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler.TransferPOST)
	hf.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: %v want %v", status, http.StatusOK)
	}

	var res handler.TransRes
	body, _ := ioutil.ReadAll(recorder.Body)
	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Errorf("Unexpected json response: %s", body)
	}
	return res
}

func TestTransfer(t *testing.T) {
	model.Init()
	var acc1 model.Account
	var acc2 model.Account
	var trans handler.TransRes

	AssertNewAccount(200, t)
	AssertNewAccount(300, t)

	trans = HandleTransfer(1, 2, 200, t)
	if trans.Code != 0 {
		t.Errorf("Unexpected err code: %v want %v", trans.Code, 0)
	}

	acc1 = AssertGetAccount(1, t)
	acc2 = AssertGetAccount(2, t)
	if acc1.Balance != 0 {
		t.Errorf("Unexpected Balance: %v want %v", acc1.Balance, 0)
	}
	if acc2.Balance != 500 {
		t.Errorf("Unexpected Balance: %v want %v", acc2.Balance, 500)
	}
}

func TestTransferFail(t *testing.T) {
	model.Init()
	var acc1 model.Account
	var acc2 model.Account
	var trans handler.TransRes

	acc1 = AssertNewAccount(200, t)
	acc2 = AssertNewAccount(300, t)

	trans = HandleTransfer(1, 2, 500, t)
	if trans.Code != model.NotEnoughMoneyCode {
		t.Errorf("Unexpected err code: %v want %v",
			trans.Code, model.NotEnoughMoneyCode)
	}

	acc1 = AssertGetAccount(1, t)
	acc2 = AssertGetAccount(2, t)
	if acc1.Balance != 200 {
		t.Errorf("Unexpected Balance: %v want %v", acc1.Balance, 200)
	}
	if acc2.Balance != 300 {
		t.Errorf("Unexpected Balance: %v want %v", acc2.Balance, 300)
	}

	trans = HandleTransfer(10, 20, 500, t)
	if trans.Code != model.NotExistsCode {
		t.Errorf("Unexpected err code: %v want %v",
			trans.Code, model.NotExistsCode)
	}
}
