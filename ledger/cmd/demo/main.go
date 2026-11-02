package main

import (
	"bufio"
	"fmt"
	"ledger"
	"os"
	"time"
)

func main() {
	fmt.Println("Ledger service started")

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("failed to open budgets file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	err = ledger.LoadBudgets(reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	budgetList := ledger.ListBudgets()
	for _, budget := range budgetList {
		fmt.Printf("%s | %.2f | %s\n",
			budget.Category,
			budget.Limit,
			budget.Period,
		)
	}

	tx1 := ledger.Transaction{
		Amount:      1000,
		Category:    "Food",
		Description: "Покупка еды в Пятерочке",
		Date:        time.Now(),
	}

	tx2 := ledger.Transaction{
		Amount:      80,
		Category:    "Transport",
		Description: "Оплата поездки в автобусе",
		Date:        time.Now(),
	}

	tx3 := ledger.Transaction{
		Amount:      50000,
		Category:    "Education",
		Description: "Оплата семестра обучения",
		Date:        time.Now(),
	}

	educationBudget := ledger.Budget{
		Category: "Education",
		Limit:    40000,
	}

	if err := ledger.SetBudget(educationBudget); err != nil {
		fmt.Println("Error:", err)
		return
	}

	tx4 := ledger.Transaction{
		Category:    "Food",
		Description: "Покупка еды в Магните",
		Date:        time.Now(),
	}

	err = ledger.AddTransaction(tx1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = ledger.AddTransaction(tx2)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = ledger.AddTransaction(tx3)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = ledger.AddTransaction(tx4)
	if err != nil {
		fmt.Println("Error:", err)
	}

	transactionList := ledger.ListTransactions()

	for _, tx := range transactionList {
		fmt.Printf("%d | %.2f | %s | %s | %v\n",
			tx.ID,
			tx.Amount,
			tx.Category,
			tx.Description,
			tx.Date.Format("2006-01-02"))
	}

	err = ledger.CheckValid(tx1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Transaction is valid")
	}

	err = ledger.CheckValid(educationBudget)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Budget is valid")
	}

}
