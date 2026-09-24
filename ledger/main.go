package main

import (
	"errors"
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

var transactions = []Transaction{}

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("Сумма не может быть равна 0")
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

func main() {
	fmt.Println("Ledger service started")

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

	tx4 := Transaction{
		Category:    "Food",
		Description: "Покупка еды в Магните",
		Date:        time.Now(),
	}

	err := AddTransaction(tx1)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	err = AddTransaction(tx2)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	err = AddTransaction(tx3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	err = AddTransaction(tx4)
	if err != nil {
		fmt.Println("Ошибка:", err)
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
