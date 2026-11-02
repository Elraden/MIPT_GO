package main

import (
	"fmt"
	"gateway/internal/api"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	handler := api.LoggingMiddleware(mux)

	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/api/transactions", api.TransactionHandler)
	mux.HandleFunc("/api/budgets", api.BudgetsHandler)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "pong")
}
