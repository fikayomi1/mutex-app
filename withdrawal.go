package main

import (
	"fmt"
	"sync"
	"time"
)

// For withdrawals (simple float balances)
var currentWithdrawalBalances = map[string]float64{
	"001": 500.0,
	"002": 150.0,
	"003": 400.0,
}

func getBalance(businessID string) float64 {
	return currentWithdrawalBalances[businessID]
}

func updateBalance(businessID string, amount float64) {
	currentWithdrawalBalances[businessID] = amount
}

func (a *App) InitiateWithdrawal(businessID string, amount float64, transactionType string) (bool, error) {
	mu, _ := a.mu.LoadOrStore(businessID, &sync.Mutex{})
	userLock := mu.(*sync.Mutex)

	userLock.Lock()

	fmt.Printf("User %s starting withdrawal of $%.2f, (%s)\n", businessID, amount, transactionType)

	currentBalance := getBalance(businessID)
	if currentBalance < amount {
		return false, fmt.Errorf("insufficient funds for %s", businessID)
	}

	time.Sleep(500 * time.Millisecond)

	updateBalance(businessID, currentBalance-amount)

	fmt.Printf("User %s successfully withdrew $%.2f. \nNew balance: $%.2f, (%s) \n",
		businessID, amount, getBalance(businessID), transactionType)
	return true, nil
}
