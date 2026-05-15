# go-todoapp

REST API приложение на Go для управления задачами (todo). Поддерживает пользователей, задачи и статистику.

## Стек

- **Go 1.25** — stdlib `net/http` для HTTP-сервера
- **PostgreSQL 18** — основное хранилище
- **pgx/v5** — драйвер и пул соединений PostgreSQL
- **golang-migrate** — управление миграциями БД
- **uber-go/zap** — структурированное логирование
- **go-playground/validator** — валидация входящих запросов
- **Docker / docker-compose** — окружение и деплой

## Архитектура

Проект организован по фичам (vertical slice):

```
internal/
├── core/                        # переиспользуемые компоненты
│   ├── config/                  # конфигурация приложения (env)
│   ├── domain/                  # базовые типы (Nullable, Uninitialized)
│   ├── errors/                  # общие sentinel-ошибки
│   ├── logger/                  # zap-обёртка + context propagation
│   ├── repository/postgres/pool/ # pgxpool
│   └── transport/http/          # сервер, роутер, middleware, request/response helpers
└── features/
    ├── users/                   # repository → service → transport/http
    ├── tasks/                   # repository → service → transport/http
    └── statistics/              # repository → service → transport/http
```

Каждая фича полностью инкапсулирована: репозиторий, сервис и HTTP-хэндлер общаются через интерфейсы.

## API

Все маршруты монтируются под префикс `/api/v1`.

### Users

| Метод  | Путь           | Описание          |
|--------|----------------|-------------------|
| POST   | /users         | Создать пользователя |
| GET    | /users         | Список пользователей |
| GET    | /users/{id}    | Получить пользователя |
| PATCH  | /users/{id}    | Обновить пользователя |
| DELETE | /users/{id}    | Удалить пользователя |

**Query params для GET /users:** `limit`, `offset`

**Тело POST /users:**
```json
{
  "full_name": "Ivan Ivanov",
  "phone_number": "+79001234567"
}
```

**Тело PATCH /users/{id}** (все поля опциональны, поддерживается null для `phone_number`):
```json
{
  "full_name": "Ivan Petrov",
  "phone_number": null
}
```

---

### Tasks

| Метод  | Путь           | Описание   |
|--------|----------------|------------|
| POST   | /tasks         | Создать задачу |
| GET    | /tasks         | Список задач |
| GET    | /tasks/{id}    | Получить задачу |
| PATCH  | /tasks/{id}    | Обновить задачу |
| DELETE | /tasks/{id}    | Удалить задачу |

**Query params для GET /tasks:** `user_id`, `limit`, `offset`

**Тело POST /tasks:**
```json
{
  "title": "Сходить в магазин",
  "description": "Купить молоко и хлеб",
  "author_user_id": 1
}
```

**Тело PATCH /tasks/{id}:**
```json
{
  "title": "Новое название",
  "completed": true
}
```

---

### Statistics

| Метод | Путь        | Описание         |
|-------|-------------|------------------|
| GET   | /statistics | Получить статистику |

**Query params:** `user_id`, `from` (дата, включительно), `to` (дата, включительно)

**Ответ:**
```json
{
  "tasks_created": 10,
  "tasks_completed": 7,
  "tasks_completed_rate": 0.7,
  "tasks_average_completion_time": "2h30m0s"
}
```

---

## Быстрый старт

### Требования

- Docker и docker-compose
- Go 1.25+ (для локального запуска)
- `make`

### 1. Настройка окружения

```bash
cp .env.example .env
```

Заполните `.env`:

```env
HTTP_ADDR=:5050
HTTP_SHUTDOWN_TIMEOUT=30s

POSTGRES_USER=todoapp
POSTGRES_PASSWORD=secret
POSTGRES_DB=todoapp
POSTGRES_TIMEOUT=10s

LOGGER_LEVEL=DEBUG

TIME_ZONE=Europe/Moscow
```

### 2. Запуск PostgreSQL

```bash
make env-up
```

### 3. Применение миграций

```bash
make migrate-up
```

### 4. Запуск приложения

**Локально (без Docker):**
```bash
make todoapp-run
```

**В Docker:**
```bash
make todoapp-deploy
```

Приложение будет доступно на `http://localhost:5050`.

---

## Makefile — все команды

### Окружение

| Команда               | Описание                                |
|-----------------------|-----------------------------------------|
| `make env-up`         | Запустить PostgreSQL                    |
| `make env-down`       | Остановить PostgreSQL                   |
| `make env-cleanup`    | Остановить и удалить данные БД          |
| `make env-port-forward` | Пробросить порт 5432 на localhost     |
| `make env-port-close` | Закрыть проброс порта                   |

### Миграции

| Команда                        | Описание                             |
|--------------------------------|--------------------------------------|
| `make migrate-up`              | Применить все миграции               |
| `make migrate-down`            | Откатить последнюю миграцию          |
| `make migrate-create seq=name` | Создать новый файл миграции          |

### Приложение

| Команда                | Описание                                     |
|------------------------|----------------------------------------------|
| `make todoapp-run`     | Запустить локально (Go + внешний PostgreSQL) |
| `make todoapp-deploy`  | Сбилдить и запустить в Docker                |
| `make todoapp-undeploy`| Остановить Docker-контейнер приложения       |
| `make ps`              | Показать запущенные контейнеры               |
| `make logs-cleanup`    | Удалить файлы логов                          |

---

## Переменные окружения

| Переменная              | По умолчанию    | Описание                          |
|-------------------------|-----------------|-----------------------------------|
| `HTTP_ADDR`             | `:5050`         | Адрес HTTP-сервера                |
| `HTTP_SHUTDOWN_TIMEOUT` | `30s`           | Таймаут graceful shutdown         |
| `POSTGRES_HOST`         | —               | Хост PostgreSQL                   |
| `POSTGRES_USER`         | —               | Пользователь PostgreSQL           |
| `POSTGRES_PASSWORD`     | —               | Пароль PostgreSQL                 |
| `POSTGRES_DB`           | —               | Имя базы данных                   |
| `POSTGRES_TIMEOUT`      | `10s`           | Таймаут подключения               |
| `LOGGER_LEVEL`          | `DEBUG`         | Уровень логирования (DEBUG/INFO/WARN/ERROR) |
| `LOGGER_FOLDER`         | —               | Папка для файлов логов            |
| `TIME_ZONE`             | `UTC`           | Временная зона приложения         |

---

## Middleware

HTTP-сервер использует цепочку middleware:

- **RequestID** — генерирует/пробрасывает заголовок `X-Request-ID`
- **Logger** — добавляет request_id и URL в контекст логгера
- **Trace** — логирует входящий запрос и итоговый статус/latency
- **Panic** — перехватывает панику и возвращает 500

---

## Схема БД

```sql
CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id           SERIAL PRIMARY KEY,
    version      BIGINT NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL,  -- 3–100 символов
    phone_number VARCHAR(15)             -- формат: +[цифры], 10–15 символов
);

CREATE TABLE todoapp.tasks (
    id             SERIAL PRIMARY KEY,
    version        BIGINT NOT NULL DEFAULT 1,
    title          VARCHAR(100) NOT NULL,   -- 1–100 символов
    description    VARCHAR(1000),           -- 1–1000 символов, nullable
    completed      BOOLEAN NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    completed_at   TIMESTAMPTZ,             -- обязателен если completed=true
    author_user_id INTEGER NOT NULL REFERENCES todoapp.users(id)
);
```