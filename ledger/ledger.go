package ledger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func ListBudgets() []Budget {
	result := make([]Budget, 0, len(budgets))

	for _, budget := range budgets {
		result = append(result, budget)
	}

	return result
}
