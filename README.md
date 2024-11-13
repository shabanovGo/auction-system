# Auction System

## Технологический стек

- Go 1.22.3
- gRPC/Protocol Buffers
- PostgreSQL 14
- Docker & Docker Compose
- gRPC-Gateway для REST API

## Архитектура

Проект построен с использованием чистой архитектуры (Clean Architecture):

### Слои приложения

- **API** (`api/proto/`) - Определения протокола Protocol Buffers
- **Interfaces** (`internal/interfaces/`) - gRPC и REST обработчики
- **Application** (`internal/application/`) - Бизнес-логика и use cases
- **Domain** (`internal/domain/`) - Бизнес-сущности и интерфейсы
- **Infrastructure** (`internal/infrastructure/`) - Реализация репозиториев

### Основные сущности

- Users - Пользователи системы
- Lots - Лоты для аукционов
- Auctions - Аукционы
- Bids - Ставки пользователей

## Установка

### Предварительные требования

- Docker и Docker Compose
- Go 1.22.3+
- Make
- Protocol Buffers compiler

### Шаги установки

1. Установить зависимости для генерации proto файлов:

```bash
make install-proto-deps
```
2. Сгенерировать код из proto файлов:

```bash
make gen-proto
```

3. Запустить приложение:

```bash
make start
```

## Команды Make

```bash
make start        # Запуск всего приложения
make up           # Запуск контейнеров
make down         # Остановка контейнеров
make build        # Сборка приложения
make migrate      # Применение миграций
make seed         # Загрузка тестовых данных
make logs         # Просмотр логов всех сервисов
make app-logs     # Логи только приложения
make test         # Запуск тестов
make test-verbose # Запуск тестов с подробным выводом
make db-logs      # Логи только базы данных
make reset-db     # Пересоздание базы данных
make restart-app  # Перезапуск приложения
```

API Endpoints
User Service
1. Создание пользователя

POST /api/v1/users

Тело запроса:
```bash
{
    "username": "ivan_petrov",
    "email": "ivan@example.com"
}
```

2. Получение пользователя

GET /api/v1/users/{id}

3. Обновление пользователя

PUT /api/v1/users/{id}

Тело запроса:
```bash
{
    "username": "new_username",
    "email": "new_email@example.com"
}
```

4. Удаление пользователя

DELETE /api/v1/users/{id}

5. Список пользователей

GET /api/v1/users?page_size=10&page_number=1

6. Обновление баланса

POST /api/v1/users/{user_id}/balance

Тело запроса:
```bash
{
    "amount": 1000.00
}
```

Lot Service
1. Создание лота

POST /api/v1/lots

Тело запроса:
```bash
{
    "title": "Антикварные часы",
    "description": "Швейцарские часы 19 века",
    "start_price": 5000.00,
    "creator_id": 1
}
```

2. Получение лота

GET /api/v1/lots/{id}

3. Обновление лота

PUT /api/v1/lots/{id}

Тело запроса:
```bash
{
    "title": "Обновленное название",
    "description": "Обновленное описание",
    "start_price": 6000.00
}
```
4. Удаление лота

DELETE /api/v1/lots/{id}

5. Список лотов

GET /api/v1/lots?page_size=10&page_number=1

Auction Service
1. Создание аукциона

POST /api/v1/auctions

Тело запроса:
```bash
{
    "lot_id": 1,
    "start_price": 5000.00,
    "min_step": 100.00,
    "start_time": "2024-03-20T10:00:00Z",
    "end_time": "2024-03-25T10:00:00Z"
}
```

2. Получение аукциона

GET /api/v1/auctions/{id}

3. Обновление аукциона
PUT /api/v1/auctions/{id}

Тело запроса:
```bash
{
    "start_price": 5500.00,
    "min_step": 150.00,
    "start_time": "2024-03-21T10:00:00Z",
    "end_time": "2024-03-26T10:00:00Z",
    "status": "ACTIVE"
}
```

4. Удаление аукциона

GET /api/v1/auctions?page_size=10&page_number=1&status=ACTIVE

Bid Service
1. Размещение ставки

POST /api/v1/bids

Тело запроса:
```bash
{
    "auction_id": 1,
    "user_id": 1,
    "amount": 5100.00
}
```
2. Получение ставки

GET /api/v1/bids/{id}

3. Список ставок

GET /api/v1/bids?auction_id=1&page_size=10&page_number=1
