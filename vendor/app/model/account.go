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
	return fmt.Sprintf("Account '%d' doesn't exist", e.ID)
}

func (e *NotEnoughMoneyError) Error() string {
	return fmt.Sprintf("Account '%d' has not enough money (%d)", e.ID, e.Summ)
}

var accountsIdSeq int64
var accountsMutex = &sync.Mutex{}
var accounts map[int64]*Account

func Init() {
	accountsIdSeq = 0
	accounts = make(map[int64]*Account)
}

func AccountNew(balance int64) *Account {
	account := &Account{
		Balance: balance,
	}

	accountsMutex.Lock()
	accountsIdSeq += 1
	account.ID = accountsIdSeq
	accounts[account.ID] = account
	accountsMutex.Unlock()

	return account
}

func AccountGet(id int64) (acc *Account, err error) {
	accountsMutex.Lock()
	acc, exists := accounts[id]
	accountsMutex.Unlock()

	if !exists {
		return nil, &NotExistsError{id}
	}
	return acc, nil
}

// Transfer money from an account to another one
func Transfer(idFrom int64, idTo int64, summ int64) error {
	accountsMutex.Lock()

	accFrom, existsFrom := accounts[idFrom]
	accTo, existsTo := accounts[idTo]
	var err error = nil

	if !existsFrom {
		err = &NotExistsError{idFrom}
	} else if !existsTo {
		err = &NotExistsError{idTo}
	} else if accFrom.Balance < summ {
		err = &NotEnoughMoneyError{idFrom, summ}
	} else {
		accFrom.Balance -= summ
		accTo.Balance += summ
	}

	accountsMutex.Unlock()

	return err
}
