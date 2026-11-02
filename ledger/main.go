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

func (tx Transaction) Validate() error {
	if tx.Amount <= 0 {
		return errors.New("amount must be > 0")
	}

	if tx.Category == "" {
		return errors.New("category cannot be empty")
	}

	if tx.Date.IsZero() {
		return errors.New("date cannot be empty")
	}
	return nil
}

type Budget struct {
	Category string
	Limit    float64
	Period   string
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("budget limit must be > 0")
	}

	if b.Category == "" {
		return errors.New("category cannot be empty")
	}

	return nil
}

type Validatable interface {
	Validate() error
}

var budgets = make(map[string]Budget)

func SetBudget(b Budget) error {
	err := b.Validate()
	if err != nil {
		return err
	}

	budgets[b.Category] = b
	return nil
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	err := tx.Validate()
	if err != nil {
		return err
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
		if err := SetBudget(budget); err != nil {
			return err
		}
	}

	return nil
}

func CheckValid(v Validatable) error {
	return v.Validate()
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

	if err := SetBudget(educationBudget); err != nil {
		fmt.Println("Error:", err)
		return
	}

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

	err = CheckValid(tx1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Transaction is valid")
	}

	err = CheckValid(educationBudget)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Budget is valid")
	}

}
