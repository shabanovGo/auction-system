# Auction System

## Технологический стек

- Go 1.22.3
- gRPC/Protocol Buffers
- PostgreSQL 14
- Docker & Docker Compose
- gRPC-Gateway для REST API
- Swagger/OpenAPI для документации API

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

1. Клонировать репозиторий:

```bash
git clone <repository-url>
cd auction-system
```

2. Установить зависимости для генерации proto файлов:

```bash
make install-proto-deps
```
3. Сгенерировать код из proto файлов:

```bash
make generate-proto
```

4. Запустить приложение:

```bash
make start
```

## Команды Make

```bash
make start       # Запуск всего приложения
make up          # Запуск контейнеров
make down        # Остановка контейнеров
make build       # Сборка приложения
make migrate     # Применение миграций
make seed        # Загрузка тестовых данных
make logs        # Просмотр логов всех сервисов
make app-logs    # Логи только приложения
make db-logs     # Логи только базы данных
make reset-db    # Пересоздание базы данных
make restart-app # Перезапуск приложения
```