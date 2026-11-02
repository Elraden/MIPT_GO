package ledger

import (
	"database/sql"
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

func Reset() {
	if db == nil {
		return
	}

	_, _ = db.Exec(`
		TRUNCATE TABLE expenses, budgets
		RESTART IDENTITY
	`)
}

type Validatable interface {
	Validate() error
}

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}

	if db == nil {
		return errors.New("database is not initialized")
	}

	query := `
		INSERT INTO budgets(category, limit_amount)
		VALUES($1, $2)
		ON CONFLICT(category)
		DO UPDATE SET limit_amount = EXCLUDED.limit_amount
	`

	_, err := db.Exec(query, b.Category, b.Limit)
	if err != nil {
		return fmt.Errorf("failed to save budget: %w", err)
	}

	return nil
}

func AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	if db == nil {
		return errors.New("database is not initialized")
	}

	var limit float64

	err := db.QueryRow(`
		SELECT limit_amount
		FROM budgets
		WHERE category = $1
	`, tx.Category).Scan(&limit)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to load budget: %w", err)
	}

	if err == nil {
		var spent float64

		err = db.QueryRow(`
			SELECT COALESCE(SUM(amount), 0)
			FROM expenses
			WHERE category = $1
		`, tx.Category).Scan(&spent)

		if err != nil {
			return fmt.Errorf("failed to calculate expenses: %w", err)
		}

		if spent+tx.Amount > limit {
			return errors.New("budget exceeded")
		}
	}

	err = db.QueryRow(`
		INSERT INTO expenses(
			amount,
			category,
			description,
			date
		)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`,
		tx.Amount,
		tx.Category,
		tx.Description,
		tx.Date,
	).Scan(&tx.ID)

	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

func ListTransactions() ([]Transaction, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	rows, err := db.Query(`
		SELECT
			id,
			amount,
			category,
			description,
			date
		FROM expenses
		ORDER BY date DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to load transactions: %w", err)
	}
	defer rows.Close()

	result := make([]Transaction, 0)

	for rows.Next() {
		var tx Transaction

		if err := rows.Scan(
			&tx.ID,
			&tx.Amount,
			&tx.Category,
			&tx.Description,
			&tx.Date,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		result = append(result, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read transactions: %w", err)
	}

	return result, nil
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

func ListBudgets() ([]Budget, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	rows, err := db.Query(`
		SELECT category, limit_amount
		FROM budgets
		ORDER BY category
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to load budgets: %w", err)
	}
	defer rows.Close()

	result := make([]Budget, 0)

	for rows.Next() {
		var budget Budget

		if err := rows.Scan(
			&budget.Category,
			&budget.Limit,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}

		result = append(result, budget)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read budgets: %w", err)
	}

	return result, nil
}
