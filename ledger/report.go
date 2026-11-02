package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportSummaryItem struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func GetReportSummary(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]ReportSummaryItem, error) {

	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	key := fmt.Sprintf(
		"report:summary:%s:%s",
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)

	if redisClient != nil {
		cached, err := redisClient.Get(ctx, key).Bytes()

		if err == nil {
			var result []ReportSummaryItem

			if err := json.Unmarshal(cached, &result); err == nil {
				log.Println("report cache hit")
				return result, nil
			}
		}

		if err != nil && !errors.Is(err, redis.Nil) {
			log.Printf("Redis read error: %v", err)
		}

		log.Println("report cache miss")
	}

	rows, err := db.QueryContext(ctx, `
		SELECT
			category,
			COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE date BETWEEN $1 AND $2
		GROUP BY category
		ORDER BY category
	`, from, to)

	if err != nil {
		return nil, fmt.Errorf("failed to build report: %w", err)
	}
	defer rows.Close()

	result := make([]ReportSummaryItem, 0)

	for rows.Next() {
		var item ReportSummaryItem

		if err := rows.Scan(
			&item.Category,
			&item.Total,
		); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read report: %w", err)
	}

	if redisClient != nil {
		data, err := json.Marshal(result)
		if err == nil {
			if err := redisClient.Set(
				ctx,
				key,
				data,
				30*time.Second,
			).Err(); err != nil {
				log.Printf("Redis write error: %v", err)
			}
		}
	}

	return result, nil
}
