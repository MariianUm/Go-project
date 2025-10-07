package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return fmt.Errorf("transaction amount cannot be zero")
	}
	tx.ID = len(transactions) + 1
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func main() {
	fmt.Println("Ledger service started")

	tx1 := Transaction{
		Amount:      1500.50,
		Category:    "Salary",
		Description: "Monthly salary",
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	tx2 := Transaction{
		Amount:      45.75,
		Category:    "Food",
		Description: "Groceries",
		Date:        time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
	}

	tx3 := Transaction{
		Amount:      120.00,
		Category:    "Entertainment",
		Description: "Cinema tickets",
	}

	if err := AddTransaction(tx1); err != nil {
		fmt.Printf("Error adding transaction: %v\n", err)
	}

	if err := AddTransaction(tx2); err != nil {
		fmt.Printf("Error adding transaction: %v\n", err)
	}

	if err := AddTransaction(tx3); err != nil {
		fmt.Printf("Error adding transaction: %v\n", err)
	}

	invalidTx := Transaction{
		Amount:   0,
		Category: "Test",
	}
	if err := AddTransaction(invalidTx); err != nil {
		fmt.Printf("Expected error for zero amount: %v\n", err)
	}

	fmt.Println("\nAll transactions:")
	allTx := ListTransactions()
	for _, tx := range allTx {
		fmt.Printf("ID: %d, Amount: $%.2f, Category: %s, Description: %s, Date: %s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description, tx.Date.Format("2006-01-02"))
	}
}
