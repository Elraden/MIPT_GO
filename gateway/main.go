package main

import (
	"context"
	"fmt"
	"gateway/internal/api"
	"ledger"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	if err := ledger.InitDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer ledger.CloseDB()

	log.Println("Connected to PostgreSQL")

	if err := ledger.InitRedis(ctx); err != nil {
		log.Fatal("Redis initialization failed:", err)
	}
	defer ledger.CloseRedis()

	log.Println("Connected to Redis")

	mux := http.NewServeMux()

	handler := api.LoggingMiddleware(mux)

	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/api/transactions", api.TransactionHandler)
	mux.HandleFunc("/api/budgets", api.BudgetsHandler)
	mux.HandleFunc("/api/reports/summary", api.ReportSummaryHandler)
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
