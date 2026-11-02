package api

import (
	"encoding/json"
	"ledger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/transactions", TransactionHandler)
	mux.HandleFunc("/api/budgets", BudgetsHandler)

	return mux
}

func performRequest(
	handler http.Handler,
	method string,
	path string,
	body string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	return recorder
}

func assertJSONContentType(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	expected := "application/json; charset=utf-8"
	actual := recorder.Header().Get("Content-Type")

	if actual != expected {
		t.Errorf(
			"expected Content-Type %q, got %q",
			expected,
			actual,
		)
	}
}

func TestBudgetsHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	router := newTestRouter()

	t.Run("post valid budget", func(t *testing.T) {
		body := `{"category":"Food","limit":500}`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/budgets",
			body,
		)

		if recorder.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var response BudgetResponse

		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Category != "Food" {
			t.Errorf("expected category Food, got %s", response.Category)
		}

		if response.Limit != 500 {
			t.Errorf("expected limit 500, got %.2f", response.Limit)
		}
	})

	t.Run("get budgets", func(t *testing.T) {
		recorder := performRequest(
			router,
			http.MethodGet,
			"/api/budgets",
			"",
		)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var responses []BudgetResponse

		if err := json.NewDecoder(recorder.Body).Decode(&responses); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(responses) != 1 {
			t.Fatalf("expected 1 budget, got %d", len(responses))
		}

		if responses[0].Category != "Food" {
			t.Errorf(
				"expected category Food, got %s",
				responses[0].Category,
			)
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		body := `{"category":"Transport","limit":0}`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/budgets",
			body,
		)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var response ErrorResponse

		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode error response: %v", err)
		}

		if response.Error == "" {
			t.Error("expected error message")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		body := `{"category":`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/budgets",
			body,
		)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)
	})
}

func TestTransactionsHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	router := newTestRouter()

	t.Run("create budget", func(t *testing.T) {
		body := `{"category":"Food","limit":500}`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/budgets",
			body,
		)

		if recorder.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				recorder.Code,
			)
		}
	})

	t.Run("ok", func(t *testing.T) {
		body := `{
			"amount":450,
			"category":"Food",
			"description":"Lunch",
			"date":"2026-11-02"
		}`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/transactions",
			body,
		)

		if recorder.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var response TransactionResponse

		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("expected ID 1, got %d", response.ID)
		}

		if response.Amount != 450 {
			t.Errorf("expected amount 450, got %.2f", response.Amount)
		}

		if response.Category != "Food" {
			t.Errorf("expected category Food, got %s", response.Category)
		}
	})

	t.Run("get transactions", func(t *testing.T) {
		recorder := performRequest(
			router,
			http.MethodGet,
			"/api/transactions",
			"",
		)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var responses []TransactionResponse

		if err := json.NewDecoder(recorder.Body).Decode(&responses); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(responses) != 1 {
			t.Fatalf(
				"expected 1 transaction, got %d",
				len(responses),
			)
		}

		if responses[0].Amount != 450 {
			t.Errorf(
				"expected amount 450, got %.2f",
				responses[0].Amount,
			)
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		body := `{
			"amount":100,
			"category":"Food",
			"description":"Coffee",
			"date":"2026-11-02"
		}`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/transactions",
			body,
		)

		if recorder.Code != http.StatusConflict {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusConflict,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var response ErrorResponse

		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode error response: %v", err)
		}

		if response.Error != "budget exceeded" {
			t.Errorf(
				"expected budget exceeded, got %q",
				response.Error,
			)
		}

		if len(ledger.ListTransactions()) != 1 {
			t.Error("rejected transaction was saved")
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		body := `{"amount":`

		recorder := performRequest(
			router,
			http.MethodPost,
			"/api/transactions",
			body,
		)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		assertJSONContentType(t, recorder)

		var response ErrorResponse

		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode error response: %v", err)
		}

		if response.Error == "" {
			t.Error("expected error message")
		}
	})
}
