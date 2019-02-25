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

func AccountNewRoutine(balance int64, chnl chan int64, t *testing.T) {
	var acc model.Account

	acc = AssertNewAccount(balance, t)

	chnl <- acc.Balance
}

func TestAccountSync(t *testing.T) {
	model.Init()
	chnl := make(chan int64)
	var summ int64 = 0
	var expectedSumm int64 = 0
	var count int64 = 0

	for i := 0; i < 10000; i++ {
		expectedSumm += int64(i)
		go AccountNewRoutine(int64(i), chnl, t)
	}

	for {
		summ += <-chnl
		count += 1
		if count == 10000 {
			close(chnl)
			break
		}
	}
	if summ != expectedSumm {
		t.Errorf("Unexpected summ: %v want %v", summ, expectedSumm)
	}
	//
	for i := 1; i <= 10000; i++ {
		AssertGetAccount(int64(i), t)
	}
	res := HandleGetAccount(10001, t)
	if res.Code != model.NotExistsCode {
		t.Errorf("Unexpected account exist! %v", res.Acc)
	}
}

func TransferRoutine(fromId int64, toId int64, summ int64,
	chnl chan int, t *testing.T) {

	var trans handler.TransRes
	trans = HandleTransfer(fromId, toId, summ, t)
	chnl <- trans.Code
}

func TestTransferSync(t *testing.T) {
	model.Init()
	chnl := make(chan int)
	var acc1 model.Account
	var acc2 model.Account
	var errCode = 0
	var count = 0

	acc1 = AssertNewAccount(10001, t)
	acc2 = AssertNewAccount(10002, t)

	for i := 0; i < 10000; i++ {
		go TransferRoutine(acc1.ID, acc2.ID, int64(1), chnl, t)
		go TransferRoutine(acc2.ID, acc1.ID, int64(1), chnl, t)
	}

	for {
		errCode = <-chnl
		if errCode != 0 {
			t.Errorf("Unexpected err code: %v want %v", errCode, 0)
		}
		count += 1
		if count == 20000 {
			close(chnl)
			break
		}
	}

	acc1 = AssertGetAccount(acc1.ID, t)
	acc2 = AssertGetAccount(acc2.ID, t)
	if acc1.Balance != 10001 {
		t.Errorf("Unexpected account balance: %v want %v", acc1.Balance, 10001)
	}
	if acc2.Balance != 10002 {
		t.Errorf("Unexpected account balance: %v want %v", acc1.Balance, 10002)
	}
}
