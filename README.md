# Publishier

Сервис для публикации текстовых постов в Telegram-канал.  
Минимальная версия с чистой архитектурой: регистрация → привязка канала → создание поста → публикация.

---

## Требования

| Инструмент | Версия |
|------------|--------|
| Docker + Docker Compose | для запуска всего стека одной командой |
| Go | 1.25+ (только для локальной разработки) |
| Node.js | 22+ (только для локальной разработки фронтенда) |
| curl / httpie | для тестирования API |

---

## Быстрый старт (Docker — рекомендуется)

Запускает PostgreSQL, бэкенд и фронтенд одной командой.

### 1. Клонировать и настроить `.env`

```bash
git clone <url-репозитория>
cd Publishier
cp .env.example .env
```

Задайте в `.env` как минимум `DB_PASSWORD`, `JWT_SECRET` и `AES_SECRET` (16/24/32 символа).

### 2. Запустить всё

```bash
docker compose up --build -d
```

| Сервис | URL |
|--------|-----|
| Фронтенд | http://localhost:5173 |
| API (напрямую) | http://localhost:8080 |
| PostgreSQL | localhost:5432 |

Проверить статус:

```bash
docker compose ps
docker compose logs -f backend
```

Остановить:

```bash
docker compose down
```

Полный сброс (включая данные БД):

```bash
docker compose down -v
```

> Фронтенд в Docker проксирует `/api` на бэкенд через nginx — отдельно запускать `npm run dev` не нужно.

---

## Локальная разработка (без Docker для приложений)

### 1. Настроить `.env`

```bash
cp .env.example .env
```

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=publishier
DB_PASSWORD=your-strong-password-here
DB_NAME=publishier

JWT_SECRET=my-super-secret-jwt-key-change-me
AES_SECRET=1234567890123456
```

> **Важно:** `AES_SECRET` должен быть **ровно 16, 24 или 32 символа**.

### 2. Поднять только PostgreSQL

```bash
docker compose up -d postgres
```

#### Альтернатива — локальный PostgreSQL

```bash
sudo -u postgres psql
```

```sql
CREATE USER publishier WITH PASSWORD 'your-strong-password-here';
CREATE DATABASE publishier OWNER publishier;
GRANT ALL PRIVILEGES ON DATABASE publishier TO publishier;
\q
```

### 3. Запустить бэкенд

```bash
set -a && source .env && set +a
cd backend && go run ./cmd
```

При успешном запуске:

```
server started on port 8080
```

### 4. Запустить фронтенд (отдельный терминал)

```bash
cd frontend
npm install
npm run dev
```

Откройте [http://localhost:5173](http://localhost:5173).

Фронтенд проксирует API-запросы на `localhost:8080` через Vite (`/api` → бэкенд).  
Бэкенд должен быть запущен до открытия фронтенда.

Сборка для продакшена:

```bash
cd frontend
npm run build
```

Артефакты появятся в `frontend/dist/`.

При старте приложение автоматически создаёт таблицы в PostgreSQL (`AutoMigrate`).

Собрать бинарник:

```bash
cd backend
go build -o bin/publishier ./cmd
./bin/publishier
```

---

## Настройка Telegram

Для публикации постов нужен **Telegram-бот** и **канал**.

### 1. Создать бота

1. Откройте [@BotFather](https://t.me/BotFather) в Telegram.
2. Отправьте `/newbot`, следуйте инструкциям.
3. Сохраните **bot token** (формат: `123456789:ABCdefGHI...`).

### 2. Создать канал

1. Создайте публичный или приватный канал в Telegram.
2. Добавьте бота в канал как **администратора** с правом публикации сообщений.

### 3. Узнать chat_id канала

**Публичный канал** — используйте `@username` канала (например `@my_channel`).

**Приватный канал** — числовой ID (обычно начинается с `-100`):

1. Опубликуйте любое сообщение в канале.
2. Откройте в браузере:
   ```
   https://api.telegram.org/bot<ВАШ_BOT_TOKEN>/getUpdates
   ```
3. Найдите `"chat":{"id":-100...}` — это ваш `chat_id`.

---

## API

Базовый URL: `http://localhost:8080`

Все защищённые эндпоинты требуют заголовок:

```
Authorization: Bearer <token>
```

### Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"mypassword"}'
```

Ответ:

```json
{"token":"eyJhbGciOiJIUzI1NiIs..."}
```

### Логин

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"mypassword"}'
```

### Привязать Telegram-канал

```bash
curl -X POST http://localhost:8080/accounts/telegram \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"bot_token":"123456789:ABCdef...","chat_id":"@my_channel"}'
```

Ответ:

```json
{"status":"connected"}
```

### Создать черновик поста

```bash
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"content":"Привет из Publishier!"}'
```

Ответ:

```json
{
  "id": "a1b2c3d4-...",
  "content": "Привет из Publishier!",
  "status": "DRAFT",
  "created_at": "2026-06-14T12:00:00Z"
}
```

### Получить пост

```bash
curl http://localhost:8080/posts/<post_id> \
  -H "Authorization: Bearer <token>"
```

### Опубликовать пост в Telegram

```bash
curl -X POST http://localhost:8080/posts/<post_id>/publish \
  -H "Authorization: Bearer <token>"
```

Ответ при успехе:

```json
{
  "id": "a1b2c3d4-...",
  "content": "Привет из Publishier!",
  "status": "PUBLISHED",
  "created_at": "2026-06-14T12:00:00Z",
  "published_at": "2026-06-14T12:01:00Z"
}
```

---

## Полный сценарий одной командой

```bash
cp .env.example .env
# отредактируйте .env — DB_PASSWORD, JWT_SECRET, AES_SECRET

docker compose up --build -d
# откройте http://localhost:5173
```

### Тест через curl (опционально)

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

curl -X POST http://localhost:8080/accounts/telegram \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"bot_token":"YOUR_BOT_TOKEN","chat_id":"@your_channel"}'

POST_ID=$(curl -s -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Мой первый пост!"}' \
  | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -X POST http://localhost:8080/posts/$POST_ID/publish \
  -H "Authorization: Bearer $TOKEN"
```

---

## Переменные окружения

| Переменная | Обязательная | По умолчанию | Описание |
|------------|:---:|--------------|----------|
| `DB_HOST` | да* | — | Хост PostgreSQL |
| `DB_PORT` | нет | `5432` | Порт PostgreSQL |
| `DB_USER` | да* | — | Пользователь БД |
| `DB_PASSWORD` | да* | — | Пароль БД |
| `DB_NAME` | нет | `publishier` | Имя базы данных |
| `DB_SSLMODE` | нет | `disable` | SSL-режим PostgreSQL |
| `DB_DSN` | нет | — | Полная DSN-строка (если задана — `DB_HOST` и др. игнорируются) |
| `JWT_SECRET` | да | — | Секрет для подписи JWT-токенов |
| `AES_SECRET` | да | — | Ключ AES (16/24/32 байта) для шифрования bot-токенов |
| `PORT` | нет | `8080` | Порт HTTP-сервера бэкенда |
| `BACKEND_PORT` | нет | `8080` | Порт бэкенда на хосте (Docker) |
| `FRONTEND_PORT` | нет | `5173` | Порт фронтенда на хосте (Docker) |
| `TELEGRAM_BASE_URL` | нет | `https://api.telegram.org` | Базовый URL Telegram Bot API |

\* Не нужны, если задан `DB_DSN`.

---

## Структура проекта

```
Publishier/
├── backend/
│   ├── cmd/main.go              # Точка входа, wiring зависимостей
│   ├── Dockerfile
│   ├── internal/
│   │   ├── api/                 # HTTP handlers, middleware
│   │   ├── auth/                # JWT
│   │   ├── config/              # Конфигурация из env
│   │   ├── crypto/              # Шифрование токенов
│   │   ├── domain/              # Сущности и ошибки
│   │   ├── publishier/          # Интеграция с Telegram
│   │   ├── repository/          # Интерфейсы репозиториев
│   │   │   └── postgres/        # GORM-реализации
│   │   └── service/             # Бизнес-логика
│   ├── configs/                 # Пример yaml-конфига
│   ├── go.mod
│   └── go.sum
├── frontend/                    # React + TypeScript UI
│   ├── Dockerfile
│   └── nginx.conf               # прокси /api → backend
├── docker-compose.yml           # postgres + backend + frontend
├── .env.example                 # Шаблон переменных окружения
└── README.md
```

---

## Устранение неполадок

### `missing env: JWT_SECRET` / `missing env: DB_HOST`

Переменные окружения не загружены. Выполните:

```bash
set -a && source .env && set +a
```

### `docker compose up` падает или Postgres не стартует

- Убедитесь, что файл `.env` существует (`cp .env.example .env`)
- Проверьте, что заданы `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Если меняли `DB_PASSWORD` после первого запуска — сбросьте volume: `docker compose down -v`

### `failed to connect db`

- Убедитесь, что PostgreSQL запущен: `docker compose ps` или `pg_isready -h localhost -p 5432`
- Проверьте `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` в `.env`
- Если PostgreSQL в Docker, а приложение тоже в Docker — используйте имя сервиса `postgres` вместо `localhost`

### `failed to init encryptor: AES key must be 16, 24 or 32 bytes`

`AES_SECRET` имеет неверную длину. Пример корректного значения: `1234567890123456` (16 символов).

### `telegram api error: ...`

- Бот не добавлен в канал как администратор
- Неверный `chat_id` (для приватного канала нужен числовой ID с `-100`)
- Неверный `bot_token`

### `social account not found`

Сначала привяжите Telegram-канал через `POST /accounts/telegram`.

### `invalid post status`

Пост уже опубликован или завершился с ошибкой. Можно публиковать только посты со статусом `DRAFT`.

### Порт 5432 занят

Если локальный PostgreSQL уже слушает 5432, измените порт в `docker-compose.yml`:

```yaml
ports:
  - "5433:5432"
```

И в `.env`: `DB_PORT=5433`.

---

## Разработка

```bash
# Сборка бэкенда
cd backend && go build ./...

# Статический анализ
cd backend && go vet ./...

# Обновить зависимости
cd backend && go mod tidy
```
