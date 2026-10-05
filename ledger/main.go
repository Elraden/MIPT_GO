package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

type Budget struct {
	Category string
	Limit    float64
	Period   string
}

var budgets = make(map[string]Budget)

func SetBudget(b Budget) {
	budgets[b.Category] = b
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("amount cannot be zero")
	}

	budget, exists := budgets[tx.Category]
	if exists {
		total := 0.0
		for _, transaction := range transactions {
			if tx.Category == transaction.Category {
				total += transaction.Amount
			}
		}
		fullAmount := tx.Amount + total
		if fullAmount > budget.Limit {
			return errors.New("budget exceeded")
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)

	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func LoadBudgets(r io.Reader) error {
	var loadedBudgets []Budget

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&loadedBudgets)
	if err != nil {
		return fmt.Errorf("failed to load budgets: %w", err)
	}

	for _, budget := range loadedBudgets {
		SetBudget(budget)
	}

	return nil
}

func main() {
	fmt.Println("Ledger service started")

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("failed to open budgets file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	err = LoadBudgets(reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, budget := range budgets {
		fmt.Printf("%s | %.2f | %s\n",
			budget.Category,
			budget.Limit,
			budget.Period,
		)
	}

	tx1 := Transaction{
		Amount:      1000,
		Category:    "Food",
		Description: "Покупка еды в Пятерочке",
		Date:        time.Now(),
	}

	tx2 := Transaction{
		Amount:      80,
		Category:    "Transport",
		Description: "Оплата поездки в автобусе",
		Date:        time.Now(),
	}

	tx3 := Transaction{
		Amount:      50000,
		Category:    "Education",
		Description: "Оплата семестра обучения",
		Date:        time.Now(),
	}

	educationBudget := Budget{
		Category: "Education",
		Limit:    40000,
	}

	SetBudget(educationBudget)

	tx4 := Transaction{
		Category:    "Food",
		Description: "Покупка еды в Магните",
		Date:        time.Now(),
	}

	err = AddTransaction(tx1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = AddTransaction(tx2)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = AddTransaction(tx3)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = AddTransaction(tx4)
	if err != nil {
		fmt.Println("Error:", err)
	}

	result := ListTransactions()

	for _, tx := range result {
		fmt.Printf("%d | %.2f | %s | %s | %v\n",
			tx.ID,
			tx.Amount,
			tx.Category,
			tx.Description,
			tx.Date.Format("2006-01-02"))
	}

}
