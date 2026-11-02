## Gateway

Запуск:

`cd gateway`
`go run .`

Сервер запускается на:

`http://localhost:8080`

## Ledger

Запуск:

`cd ledger`
`go run .`

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
  -d '{"amount":450,"category":"Food","description":"Lunch","date":"2026-11-02"}'
`

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
