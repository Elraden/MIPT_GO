## Gateway

Запуск:

`cd gateway`  
`go run .`

Сервер запускается на:

`http://localhost:8080`

При запуске Gateway подключается к PostgreSQL и Redis.

При успешном подключении:

`Connected to PostgreSQL`  
`Connected to Redis`

## Ledger

Ledger содержит бизнес-логику приложения.

Ledger отвечает за:

- работу с PostgreSQL;
- хранение бюджетов и транзакций;
- проверку лимитов бюджета;
- формирование отчётов;
- кеширование отчётов в Redis.

Gateway использует пакет Ledger для работы с данными.

## PostgreSQL

Для хранения бюджетов и транзакций используется PostgreSQL.

Запуск PostgreSQL через Docker:

`docker run --name cashcraft-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp -p 5433:5432 -d postgres:17-alpine`

### Переменная подключения

Для Bash:

`export DATABASE_URL="postgres://postgres:postgres@localhost:5433/cashapp?sslmode=disable"`

Если `DATABASE_URL` не задана, используются переменные:

`DB_HOST`  
`DB_PORT`  
`DB_USER`  
`DB_PASS`  
`DB_NAME`

## Миграции Goose

Миграции находятся в:

`ledger/migrations`

### Применить миграции:

`goose -dir ./ledger/migrations postgres "$DATABASE_URL" up`

### Проверить состояние миграций:

`goose -dir ./ledger/migrations postgres "$DATABASE_URL" status`

### Откатить последнюю миграцию:

`goose -dir ./ledger/migrations postgres "$DATABASE_URL" down`

Миграции необходимо применить перед первым запуском Gateway.

## Redis

Redis используется для кеширования результатов отчёта.

Запуск Redis через Docker:

`docker run -p 6379:6379 --name cashcraft-redis -d redis:7-alpine`

### Переменные окружения:

`REDIS_ADDR=localhost:6379`

`REDIS_DB=0`

`REDIS_PASSWORD=`

По умолчанию используется:

`localhost:6379`

и база Redis:

`0`

## Ping

`GET /ping`

#### Пример:

`curl http://localhost:8080/ping`

## Transactions

### Получить список транзакций:

`GET /api/transactions`

#### Пример:

`curl http://localhost:8080/api/transactions`

### Создать транзакцию:

`POST /api/transactions`

#### Пример:

`curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":450,"category":"Food","description":"Lunch","date":"2026-11-02"}'`

Транзакции сохраняются в PostgreSQL и не теряются после перезапуска Gateway.

Если сумма транзакций превышает установленный бюджет категории, сервер возвращает:

`409 Conflict`

Пример ответа:

`{"error":"budget exceeded"}`

## Budgets

### Получить список бюджетов:

`GET /api/budgets`

#### Пример:

`curl http://localhost:8080/api/budgets`

### Создать или обновить бюджет:

`POST /api/budgets`

#### Пример:

`curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"Food","limit":500}'`

Если бюджет для указанной категории уже существует, его лимит обновляется.

Бюджеты сохраняются в PostgreSQL.

## Reports

### Получить отчёт по расходам:

`GET /api/reports/summary`

Параметры:

`from` — начальная дата.

`to` — конечная дата.

Формат даты:

`YYYY-MM-DD`

#### Пример:

`curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"`

При первом запросе данные получаются из PostgreSQL и сохраняются в Redis на 30 секунд.

В логах:

`report cache miss`

Если выполнить тот же запрос повторно в течение 30 секунд:

`curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"`

данные будут получены из Redis.

В логах:

`report cache hit`

## Пул соединений PostgreSQL

Для подключения к PostgreSQL используются настройки:

`MaxOpenConns = 10`

`MaxIdleConns = 5`

`ConnMaxLifetime = 30 минут`

## ТЕСТЫ

### Запуск тестов Ledger:

`cd ledger`  
`go test ./...`

### Запуск тестов Gateway:

`cd gateway`  
`go test ./...`

### Запуск с подробным выводом:

`go test ./... -v`

#### Покрытие Ledger:

`cd ledger`  
`go test . -coverprofile=cover-ledger.out`  
`go tool cover -func=cover-ledger.out`

#### Покрытие Gateway API:

`cd gateway`  
`go test ./internal/api -coverprofile=cover-gateway.out`  
`go tool cover -func=cover-gateway.out`

### Текущее покрытие:

#### Ledger: 63.3%

#### Gateway API: 78.3%
