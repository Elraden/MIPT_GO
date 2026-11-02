package api

import (
	"encoding/json"
	"ledger"
	"net/http"
	"time"
)

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req CreateTransactionRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		tx, err := CreateTransaction(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = tx.Validate()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = ledger.AddTransaction(tx)
		if err != nil {
			if err.Error() == "budget exceeded" {
				writeJSONError(w, http.StatusConflict, err.Error())
				return
			}
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		transactions, err := ledger.ListTransactions()
		if err != nil {
			writeJSONError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}
		if len(transactions) == 0 {
			writeJSONError(w, http.StatusInternalServerError, "transaction was not saved")
			return
		}
		tx = transactions[0]
		response := ToTransactionResponse(tx)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	case http.MethodGet:
		transactions, err := ledger.ListTransactions()
		if err != nil {
			writeJSONError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}
		responses := make([]TransactionResponse, 0)
		for _, tx := range transactions {
			response := ToTransactionResponse(tx)
			responses = append(responses, response)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func BudgetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req CreateBudgetRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		budget := CreateBudget(req)

		err = budget.Validate()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = ledger.SetBudget(budget)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := ToBudgetResponse(budget)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	case http.MethodGet:
		budgets, err := ledger.ListBudgets()
		if err != nil {
			writeJSONError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}
		responses := make([]BudgetResponse, 0)
		for _, budget := range budgets {
			response := ToBudgetResponse(budget)
			responses = append(responses, response)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(
		ErrorResponse{Error: message},
	)
}

func ReportSummaryHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		writeJSONError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	fromString := r.URL.Query().Get("from")
	toString := r.URL.Query().Get("to")

	from, err := time.Parse("2006-01-02", fromString)
	if err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid from date",
		)
		return
	}

	to, err := time.Parse("2006-01-02", toString)
	if err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid to date",
		)
		return
	}

	if from.After(to) {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"from date cannot be after to date",
		)
		return
	}

	result, err := ledger.GetReportSummary(
		r.Context(),
		from,
		to,
	)

	if err != nil {
		writeJSONError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}
