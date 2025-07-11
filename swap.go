package main

import (
	"fmt"
	"sync"
	"time"
)

func getExchangeRate(from, to string) float64 {
	if from == "USD" && to == "NGN" {
		return 1500 // 1 USD = 1500 NGN
	}
	return 1 / 1500 // 1 NGN = 0.000666 USD
}

func (a *App) InitiateSwap(userID, fromCurrency, toCurrency string, amount float64) error {
	// 1. Get or create user-specific mutex
	mu, _ := a.mu.LoadOrStore(userID, &sync.Mutex{})
	userLock := mu.(*sync.Mutex)

	// 2. Lock with visible start/end markers
	a.logger.Printf("[%s] WAITING FOR LOCK (%s→%s)", userID, fromCurrency, toCurrency)
	userLock.Lock()
	defer func() {
		userLock.Unlock()
		a.logger.Printf("[%s] 🔓 LOCK RELEASED", userID)
	}()

	a.logger.Printf("[%s] ▶ START SWAP %.2f %s→%s", userID, amount, fromCurrency, toCurrency)

	// 3. Verify balance
	if a.balances[userID][fromCurrency] < amount {
		return fmt.Errorf("insufficient funds")
	}

	// 4. Calculate exchange (simplified)
	rate := getExchangeRate(fromCurrency, toCurrency)
	toAmount := amount * rate

	// 5. Simulate processing delay
	time.Sleep(300 * time.Millisecond)

	// 6. Update balances
	a.balances[userID][fromCurrency] -= amount
	a.balances[userID][toCurrency] += toAmount

	a.logger.Printf("[%s] ✔ SWAP COMPLETE: %.2f %s → %.2f %s",
		userID, amount, fromCurrency, toAmount, toCurrency)
	return nil
}
