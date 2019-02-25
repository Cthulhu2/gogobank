package model

import (
	"fmt"
	"sync"
)

// swagger:model Account
type Account struct {
	ID      int64 `json: "id"`
	Balance int64 `json: "balance"`
}

var NotExistsCode = 1

type NotExistsError struct {
	ID int64
}

var NotEnoughMoneyCode = 2

type NotEnoughMoneyError struct {
	ID   int64
	Summ int64
}

func (e *NotExistsError) Error() string {
	return fmt.Sprintf("Account '%d' doesn't exists", e.ID)
}

func (e *NotEnoughMoneyError) Error() string {
	return fmt.Sprintf("Account '%d' has not enough money (%d)", e.ID, e.Summ)
}

var accountsIdSeq int64
var accountsIdSeqMutex = &sync.Mutex{}

var accounts map[int64]*Account
var transferMutex = &sync.Mutex{}

func Init() {
	accountsIdSeq = 0
	accounts = make(map[int64]*Account)
}

func AccountNew(balance int64) *Account {
	account := &Account{
		Balance: balance,
	}

	accountsIdSeqMutex.Lock()
	accountsIdSeq += 1
	account.ID = accountsIdSeq
	accounts[account.ID] = account
	accountsIdSeqMutex.Unlock()

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
func Transfer(idFrom int64, idTo int64, summ int64) error {
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

	transferMutex.Lock()
	if accFrom.Balance < summ {
		transferMutex.Unlock()
		return &NotEnoughMoneyError{idFrom, summ}
	}
	accFrom.Balance -= summ
	accTo.Balance += summ
	transferMutex.Unlock()
	return nil
}
