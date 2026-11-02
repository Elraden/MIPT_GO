package api

import (
	"encoding/json"
	"ledger"
	"net/http"
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

		transactions := ledger.ListTransactions()
		if len(transactions) == 0 {
			writeJSONError(w, http.StatusInternalServerError, "transaction was not saved")
			return
		}
		tx = transactions[len(transactions)-1]
		response := ToTransactionResponse(tx)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	case http.MethodGet:
		transactions := ledger.ListTransactions()
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
		budgets := ledger.ListBudgets()
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
