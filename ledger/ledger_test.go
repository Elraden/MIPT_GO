package ledger

import (
	"testing"
	"time"
)

func TestTransactionValidate(t *testing.T) {
	validDate := time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				Amount:   1000,
				Category: "Food",
				Date:     validDate,
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			tx: Transaction{
				Amount:   0,
				Category: "Food",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			tx: Transaction{
				Amount:   -100,
				Category: "Food",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			tx: Transaction{
				Amount:   1000,
				Category: "",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "empty date",
			tx: Transaction{
				Amount:   1000,
				Category: "Food",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.tx.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}

func TestBudgetValidate(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name: "valid budget",
			budget: Budget{
				Category: "Food",
				Limit:    5000,
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: Budget{
				Category: "Food",
				Limit:    0,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			budget: Budget{
				Category: "Food",
				Limit:    -1000,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			budget: Budget{
				Category: "",
				Limit:    5000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.budget.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}

func TestAddTransactionBudgetExceeded(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	err := SetBudget(Budget{
		Category: "Food",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("failed to set budget: %v", err)
	}

	firstTransaction := Transaction{
		Amount:   3000,
		Category: "Food",
		Date:     time.Now(),
	}

	if err := AddTransaction(firstTransaction); err != nil {
		t.Fatalf("failed to add first transaction: %v", err)
	}

	if len(ListTransactions()) != 1 {
		t.Fatalf("expected 1 transaction")
	}

	secondTransaction := Transaction{
		Amount:   2500,
		Category: "Food",
		Date:     time.Now(),
	}

	err = AddTransaction(secondTransaction)

	if err == nil {
		t.Fatal("expected budget exceeded error")
	}

	if err.Error() != "budget exceeded" {
		t.Errorf("expected budget exceeded, got %v", err)
	}

	transactions := ListTransactions()

	if len(transactions) != 1 {
		t.Errorf(
			"expected rejected transaction not to be saved, got %d transactions",
			len(transactions),
		)
	}
}
