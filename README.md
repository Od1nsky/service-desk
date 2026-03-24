# АИС «Контроль текущей успеваемости студентов»

Веб-приложение для регистрации, хранения и анализа результатов учебной деятельности студентов.

---

## Стек

| Слой | Технология |
|------|-----------|
| Frontend | Nuxt 3 (SPA) + TypeScript, Naive UI, Pinia 2, Axios |
| Backend | Go 1.21+, Gin v1.9, raw SQL (database/sql + lib/pq) |
| СУБД | PostgreSQL 15+ |
| Аутентификация | JWT HS256 (Access 15 мин + Refresh 30 дней) |
| E2E-тесты | Playwright |

---

## Быстрый старт

### 1. База данных

```bash
psql -U <пользователь> -d <база> -f backend/migrations/001_init.sql
```

### 2. Backend

```bash
cd backend
cp .env.example .env   # заполнить переменные окружения
go run ./cmd/main.go
```

**Переменные окружения (`.env`):**

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=servicedesk
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h
SERVER_PORT=8080
CORS_ORIGINS=http://localhost:3000
```

API будет доступен на `http://localhost:8080/api/v1/`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Приложение откроется на `http://localhost:3000`.

---

## Роли пользователей

| Роль | Возможности |
|------|-------------|
| `student` | Просмотр своих оценок, истории изменений, добавление комментариев |
| `teacher` | Выставление и редактирование оценок, управление статусами, заметки |
| `admin` | Всё выше + управление пользователями, дисциплинами, типами оценок, дашборд |

---

## Типы оценок (grade_types)

| Тип | Описание | max_score |
|-----|---------|-----------|
| Текущая | Лабораторные, практические работы, тесты | 100 |
| Рубежная | Рубежный контроль (mid-term) | 100 |
| Зачёт | Итоговая аттестация «зачтено / не зачтено» | 1 |
| Экзамен | Итоговая аттестация с баллом 0–100 | 100 |

---

## API

Базовый путь: `/api/v1/`

| Метод | Endpoint | Описание |
|-------|----------|---------|
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход (получение JWT) |
| POST | `/auth/refresh` | Обновление access-токена |
| POST | `/auth/logout` | Выход |
| GET | `/grades` | Список оценок (с фильтрами) |
| POST | `/grades` | Создать оценку (teacher) |
| GET | `/grades/:id` | Детали оценки + история |
| PATCH | `/grades/:id/status` | Изменить статус / балл |
| PATCH | `/grades/:id/assign` | Назначить преподавателя |
| PATCH | `/grades/:id/acknowledge` | Студент подтверждает ознакомление |
| GET/POST | `/grades/:id/comments` | Комментарии к оценке |
| GET | `/disciplines` | Список дисциплин |
| POST/PUT/DELETE | `/disciplines/:id` | CRUD дисциплин (admin) |
| GET | `/grade-types` | Список типов оценок |
| POST/PUT/DELETE | `/grade-types/:id` | CRUD типов оценок (admin) |
| GET | `/users` | Список пользователей (admin) |
| PATCH | `/users/:id/role` | Сменить роль (admin) |
| PATCH | `/users/:id/block` | Заблокировать/разблокировать (admin) |
| GET | `/stats/dashboard` | Дашборд успеваемости (teacher, admin) |

---

## Статусы оценки

```
pending → graded → certified → closed
                ↘ not_certified → graded (пересдача)
```

---

## E2E-тесты

```bash
cd frontend
npx playwright test
```

Скриншоты сохраняются в `reports/screenshots/`.

---

## Документация

| Файл | Содержание |
|------|-----------|
| `01_project_overview.md` | Обзор проекта, стек, структура файлов |
| `02_tasks_projectlibre.md` | План проекта, диаграмма Ганта |
| `03_analytics.md` | Аналитика: use cases, ERD, API, RBAC |
| `04_tz_gost.md` | Техническое задание (ГОСТ 34.602-2020) |
