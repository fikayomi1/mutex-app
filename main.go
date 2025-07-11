package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type App struct {
	mu       sync.Map // userID -> *sync.Mutex
	balances map[string]map[string]float64
	logger   *log.Logger
}

func main() {
	// Initialize with test data
	account := App{
		balances: map[string]map[string]float64{
			"user1": {"USD": 1000, "NGN": 0},
			"user2": {"USD": 500, "NGN": 0},
		},
		logger: log.New(os.Stdout, "", log.LstdFlags), // Initialize logger
	}

	// Print initial balances
	printBalances := func() {
		fmt.Println("\nCurrent Balances:")
		for user, balances := range account.balances {
			fmt.Printf("  %s: USD=%.2f, NGN=%.2f\n", user, balances["USD"], balances["NGN"])
		}
	}
	printBalances()

	// Test scenario
	var wg sync.WaitGroup

	// User1 swaps (should execute sequentially)
	wg.Add(4) // Corrected to match number of goroutines
	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond)
		err := account.InitiateSwap("user2", "NGN", "USD", 75000)
		if err != nil {
			account.logger.Printf("Error in swap: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := account.InitiateSwap("user1", "USD", "NGN", 100)
		if err != nil {
			account.logger.Printf("Error in swap: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := account.InitiateSwap("user2", "USD", "NGN", 50)
		if err != nil {
			account.logger.Printf("Error in swap: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond) // Stagger slightly
		err := account.InitiateSwap("user1", "USD", "NGN", 200)
		if err != nil {
			account.logger.Printf("Error in swap: %v", err)
		}
	}()

	wg.Wait()
	printBalances()
}
