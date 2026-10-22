# Task API: REST API для списка задач на Go

Учебный и портфолио-проект. HTTP API на Go: регистрация и вход по JWT, личный список задач у каждого пользователя. Данные хранятся в PostgreSQL. Приложение и база запускаются через Docker Compose.

Стек: Go, chi, pgx, golang-jwt, bcrypt, log/slog, PostgreSQL 16, Docker Compose.

## Быстрый запуск

1. Запустите Docker Desktop и дождитесь статуса Engine running.
2. В папке проекта выполните:

```powershell
docker compose up --build
```

3. API будет доступен на `http://localhost:8080`. Проверка: `GET http://localhost:8080/healthz` должен вернуть `200`.

Установленный локально Go для запуска через Docker не нужен. Остановить: `Ctrl+C` в окне compose или `docker compose stop`. Данные остаются в томе `pgdata`.

## Настройки

Приложение читает переменные окружения:

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `DATABASE_URL` | строка подключения к PostgreSQL | обязательна |
| `JWT_SECRET` | секрет для подписи токенов | обязательна |
| `ADDR` | адрес и порт HTTP-сервера | `:8080` |

В `docker-compose.yml` они уже заданы для контейнера. Значение `JWT_SECRET` там учебное, для реального использования замените его. Файл `.env.example` служит только шаблоном: приложение само `.env` не читает.

PostgreSQL из Compose опубликован на порту `5433` хоста (`"5433:5432"`), чтобы не конфликтовать с локальным Postgres. Параметры учебной базы: пользователь, пароль и имя базы `taskapi`.

## Миграции

SQL-файлы из `migrations/` подключены к контейнеру базы как `/docker-entrypoint-initdb.d` и выполняются автоматически только при первом запуске пустого тома. Если том уже создан, новые файлы сами не применятся. Применить миграцию вручную, не трогая данные:

```powershell
Get-Content .\migrations\002_tasks.sql | docker compose exec -T db psql -U taskapi -d taskapi
```

Проверить список таблиц (должны быть `users` и `tasks`):

```powershell
docker compose exec db psql -U taskapi -d taskapi -c "\dt"
```

`docker compose down -v` пересоздаёт том, но удаляет все данные, включая пользователей и задачи.

## Работа с API

Все ответы в JSON. Ошибки имеют вид `{"error": "сообщение"}`. Запросы к `/tasks` требуют заголовок `Authorization: Bearer <token>`.

| Метод | Путь | Описание | Успех |
|---|---|---|---|
| GET | `/healthz` | проверка, что сервис жив | 200 |
| POST | `/auth/register` | регистрация | 201 |
| POST | `/auth/login` | вход, возвращает токен | 200 |
| POST | `/tasks` | создать задачу | 201 |
| GET | `/tasks` | список задач пользователя | 200 |
| GET | `/tasks/{id}` | одна задача | 200 |
| PATCH | `/tasks/{id}` | изменить `title` и/или `done` | 200 |
| DELETE | `/tasks/{id}` | удалить задачу | 204 |

Примеры тел запросов:

```json
{"email": "me@example.com", "password": "password123"}
```

```json
{"title": "first task"}
```

```json
{"done": true}
```

Правила:

- Email приводится к нижнему регистру и проверяется на формат, он уникален.
- Пароль от 8 до 72 байт (ограничение bcrypt).
- Токен JWT (HS256) действует 24 часа.
- Название задачи не пустое и не длиннее 200 символов. В `PATCH` можно передать любое из полей, пустое тело `{}` даёт `400`.
- Пользователь видит и меняет только свои задачи. Чужая задача возвращает `404`, как и несуществующая.
- Лишние поля в JSON и битый JSON дают `400`.

Коды ошибок: `400` неверный ввод, `401` нет или неверный токен или пароль, `404` не найдено, `409` email уже занят, `500` внутренняя ошибка (подробности в логах контейнера).

### Пример в PowerShell

```powershell
$base = "http://localhost:8080"
$json = "application/json"

Invoke-RestMethod -Method Post -Uri "$base/auth/register" -ContentType $json -Body '{"email":"me@example.com","password":"password123"}'
$token = (Invoke-RestMethod -Method Post -Uri "$base/auth/login" -ContentType $json -Body '{"email":"me@example.com","password":"password123"}').token
$h = @{ Authorization = "Bearer $token" }

Invoke-RestMethod -Method Post -Uri "$base/tasks" -Headers $h -ContentType $json -Body '{"title":"first task"}'
Invoke-RestMethod -Uri "$base/tasks" -Headers $h
```

В PowerShell `curl` является псевдонимом `Invoke-WebRequest`. Если нужен именно curl, используйте `curl.exe`.

### Postman

Метод и URL как в таблице, тело: Body -> raw -> JSON. Для запросов к `/tasks` во вкладке Authorization выберите Bearer Token и вставьте токен из ответа `/auth/login` без кавычек и без слова `Bearer`.

## Запуск без Docker для приложения

Нужен Go (версия указана в `go.mod`). База поднимается в Docker, приложение запускается локально:

```powershell
docker compose up -d db
$env:DATABASE_URL = "postgres://taskapi:taskapi@localhost:5433/taskapi?sslmode=disable"
$env:JWT_SECRET = "dev-secret"
go run ./cmd/api
```

## Тесты

```powershell
go test -race ./...
```

Если Go локально не установлен, можно запустить в контейнере:

```powershell
docker run --rm -v "${PWD}:/src" -w /src golang:1.27 go test -race ./...
```

Что покрыто:

- `internal/service`: валидация названия задачи, регистрация, вход и разбор токена (на заглушках без базы).
- `internal/handler`: маршруты, middleware авторизации, разбор JSON и коды ответов (через `httptest`, на заглушках сервисов).

Тестов на настоящей базе пока нет, SQL проверяется вручную через API.

## Структура

```text
cmd/api/main.go          запуск, сборка зависимостей, graceful shutdown
internal/
  config/                чтение переменных окружения
  domain/                сущности и общие ошибки
  handler/               HTTP: маршруты, middleware, JSON-ответы
  service/               бизнес-правила, валидация, JWT
  repository/            SQL-запросы через pgx
migrations/              схема базы (users, tasks)
Dockerfile               сборка в два этапа, запуск от обычного пользователя
docker-compose.yml       PostgreSQL и API
Makefile                 run, test, lint, docker-up, docker-down
```

Зависимости направлены в одну сторону: `handler -> service -> repository`. Сервис и обработчики зависят от интерфейсов, поэтому их можно тестировать на заглушках.

## Если не запускается

- `failed to connect to the docker API` или `dockerDesktopLinuxEngine`: Docker Desktop не запущен. Запустите его и проверьте `docker info`.
- `"/go.sum": not found` при сборке: нет файла `go.sum`. Выполните `go mod tidy` (локально или через Docker) и закоммитьте `go.mod` и `go.sum`.
- `package ... is not in std` при сборке: имя модуля в импортах не совпадает с первой строкой `go.mod`. Везде должно быть `taskapi/internal/...`.
- `undefined: jwt`, `chi`, `pgx` и подобное: из файла пропали импорты (часто это автоочистка импортов в IDE при сохранении). Восстановите импорты и отключите Optimize imports on save.
- `relation "tasks" does not exist` (в логах и ответ `500`): таблица не создана, скрипты не выполнились для существующего тома. Примените миграцию вручную (см. раздел «Миграции») или пересоздайте том через `docker compose down -v`.
- `invalid token`: токен вставлен с кавычками или со словом `Bearer`, устарел или выдан до смены `JWT_SECRET`. Залогиньтесь заново.
- `409 email already taken`: пользователь уже зарегистрирован, переходите к `/auth/login`.
- Порт `5433` или `8080` занят: измените публикацию порта в `docker-compose.yml` и `DATABASE_URL` для локального запуска.

Подробности любой ошибки `500` смотрите в логах: `docker compose logs api --tail 20`.

## Ограничения и планы

- Миграции выполняются только скриптами инициализации тома. Планируется переход на goose с отдельным шагом применения.
- Нет пагинации и фильтра списка задач.
- Нет refresh-токенов и ограничения частоты запросов.
- Нет интеграционных тестов на реальной базе (планируется testcontainers).
- CI (GitHub Actions, `golangci-lint`) пока не настроен.
