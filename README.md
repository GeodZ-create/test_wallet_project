# Wallet API

REST API для управления балансом кошелька.

Поддерживаются операции:
- пополнение кошелька
- списание средств
- получение текущего баланса

Стек:
- Go
- PostgreSQL
- GORM - для работы с PostgreSQL
- goose - для миграций
- Docker Compose

## Конфигурация

Проект читает переменные окружения из `config.env`.

Пример файла лежит в config.env.example.

Основные переменные:
- `PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `DATABASE_URL`

Для локального запуска создать `config.env` на основе `config.env.example`.

## Локальный запуск

1. Поднять PostgreSQL локально или через Docker.
2. Создать `config.env`.
3. Прогнать миграции.
4. Запустить приложение.

Пример команд:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migration postgres "host=localhost user=postgres password=change_me dbname=wallet_db port=5432 sslmode=disable" up
go run main.go
```

## Запуск через Docker Compose

1. Создать локальный `config.env`.
2. Поднять БД:

```bash
docker compose up -d db
```

3. Прогнать миграции:

```bash
docker compose run --rm migrate up
```

4. Поднять приложение:

```bash
docker compose up -d app
```

PostgreSQL проброшен на локальный порт `5433`, приложение доступно на `9091`.

## API

### POST `/api/v1/wallet`

Изменение баланса кошелька.

Пример запроса:

```json
{
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "deposit",
  "amount": 100
}
```

Поддерживаемые `operationType`:
- `deposit`
- `withdraw`

Пример успешного ответа:

```json
{
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "deposit",
  "amount": 100,
  "balance": 150
}
```

### GET `/api/v1/wallets/{WALLET_UUID}`

Получение текущего баланса кошелька.

Пример ответа:

```json
{
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "balance": 150
}
```

### Ошибки

Пример формата ошибки:

```json
{
  "error": "wallet not found"
}
```

Основные коды ответа:
- `400` — невалидный запрос
- `404` — кошелек не найден
- `409` — недостаточно средств
- `500` — внутренняя ошибка сервера

## Тесты

Запуск всех тестов:

```bash
go test ./...
```

Отдельно:

```bash
go test ./repository
go test ./httpHandlers
```

Что покрыто:
- интеграционные тесты `repository`
- unit-тесты `httpHandlers`
- конкурентный тест на параллельные пополнения одного кошелька

## Нагрузочная проверка

Для ручной проверки нагрузки использовался `hey`.

Сценарий:
- все запросы отправлялись на один и тот же `walletId`
- использовался `POST /api/v1/wallet`
- операция: `deposit`
- сумма: `1`

Пример:

```bash
hey -n 20000 -c 250 -m POST -H "Content-Type: application/json" -D body.json http://localhost:9091/api/v1/wallet
```

`body.json`:

```json
{
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "deposit",
  "amount": 1
}
```
Результат одного из прогонов:

- 20000 / 20000 ответов с кодом 200
- 0 ответов 5xx
- 4221 RPS
- Средняя продолжительность: 55.9 ms

После нагрузочного прогона дополнительно проверялся итоговый баланс кошелька.
Баланс соответствовал числу успешно обработанных запросов.
