
# Medods - Сервис аутентификации

## Описание ТЗ
Реализовать сервис, который будет получать по API ФИО, из открытых API обогащать
ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в
БД. По запросу выдавать инфу о найденных людях.

## Архитектура

Проект построен с использованием Clean Architecture, что обеспечивает высокую гибкость и масштабируемость.

- **Model Layer** - содержит бизнес-модели
- **Service Layer** - реализует бизнес-логику
- **Repository Layer** - отвечает за работу с данными
- **Handler Layers** - обрабатывают HTTP запросы


## Технологии

- Go 1.21+
- PostgreSQL 14+
- Gorilla Mux Framework
- Slog (для логирования)
- Docker & Docker Compose
- Swagger (API документация)

## Требования

- Go 1.21 или выше
- Docker и Docker Compose
- PostgreSQL 14 или выше (если запуск без Docker)

## Установка и запуск

### Через Docker

1. Клонируйте репозиторий
2. Создайте `.env` файл на основе `.env.example` (он описан ниже)
3. Запустите через Docker Compose:
   ```bash
   docker-compose up -d


## Конфигурация ENV

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
SERVER_PORT=8085
```



## Rest методы
1. GET /persons/
2. POST /persons/
3. PUT /persons/{id}/
4. DELETE /persons/{id}/

## Swagger
Методы детально описаны в swagger и доступны по маршуту:

	http://localhost:8085/swagger/index.html