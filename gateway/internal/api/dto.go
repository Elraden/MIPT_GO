package api

import (
	"ledger"
	"time"
)

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionResponse struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

func CreateTransaction(req CreateTransactionRequest) (ledger.Transaction, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ledger.Transaction{}, err
	}
	tx := ledger.Transaction{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        date,
	}
	return tx, nil
}

func ToTransactionResponse(tx ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date.Format("2006-01-02"),
	}
}

func CreateBudget(req CreateBudgetRequest) ledger.Budget {
	return ledger.Budget{
		Category: req.Category,
		Limit:    req.Limit,
	}
}

func ToBudgetResponse(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit:    b.Limit,
		Period:   b.Period,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
