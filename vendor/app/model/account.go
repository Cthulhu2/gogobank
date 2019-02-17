package model

import (
	"fmt"
	"sync"
)

type Account struct {
	ID      int64 `json: "id"`
	Balance int64 `json: "balance"`
}

type NotExistsError struct {
	ID int64
}

type NotEnoughMoneyError struct {
	ID   int64
	Summ int64
}

func (e *NotExistsError) Error() string {
	return fmt.Sprintf("'%d' account is not exists", e.ID)
}

func (e *NotEnoughMoneyError) Error() string {
	return fmt.Sprintf("'%d' account has not enough money (%d)", e.ID, e.Summ)
}

var accountsIdSeq int64 = 0
var accountsIdSeqMutex = &sync.Mutex{}

var accounts = make(map[int64]*Account)
var transactionMutex = &sync.Mutex{}

func AccountNew(balance int64) *Account {
	accountsIdSeqMutex.Lock()
	accountsIdSeq += 1
	id := accountsIdSeq
	accountsIdSeqMutex.Unlock()

	account := &Account{
		ID:      id,
		Balance: balance,
	}
	accounts[id] = account
	return account
}

func AccountGet(id int64) (acc *Account, err error) {
	acc, exists := accounts[id]
	if !exists {
		return nil, &NotExistsError{id}
	}
	return acc, nil
}

// Transfer money from account to another one
func Transaction(idFrom int64, idTo int64, summ int64) error {
	accFrom, exists := accounts[idFrom]
	if !exists {
		return &NotExistsError{idFrom} //
	}
	accTo, exists := accounts[idTo]
	if !exists {
		return &NotExistsError{idTo} //
	}

	if idFrom == idTo || summ == 0 {
		return nil // OK
	}

	transactionMutex.Lock()
	if accFrom.Balance <= summ {
		transactionMutex.Unlock()
		return &NotEnoughMoneyError{idFrom, summ}
	}
	accFrom.Balance -= summ
	accTo.Balance += summ
	transactionMutex.Unlock()
	return nil
}
