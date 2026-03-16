# Аналитические материалы — АС «ServiceDesk»

---

## 1. Описание предметной области

### 1.1 Контекст

Техническая поддержка (helpdesk/servicedesk) — это процесс приёма, регистрации и устранения инцидентов и запросов, поступающих от сотрудников организации в адрес IT-подразделения.

Без автоматизации типичные проблемы:
- Заявки теряются в мессенджерах и почте
- Нет единой истории обращений
- Невозможно оценить нагрузку на специалистов
- Отсутствует контроль соблюдения SLA (времени реакции)

### 1.2 Процесс «как есть» (As-Is)

```
Сотрудник звонит / пишет в чат
        ↓
Специалист отвечает когда видит сообщение
        ↓
Решает устно, без документирования
        ↓
Повторные обращения по тем же проблемам
        ↓
Руководство не знает о загрузке отдела
```

### 1.3 Процесс «как должно быть» (To-Be)

```
Сотрудник открывает веб-приложение
        ↓
Создаёт заявку: тема, описание, категория, приоритет
        ↓
Система регистрирует заявку, присваивает номер
        ↓
Специалист видит очередь → берёт заявку в работу
        ↓
Переписка в комментариях, смена статусов
        ↓
Заявка закрывается → сотрудник оценивает решение
        ↓
Менеджер анализирует отчёты
```

---

## 2. Глоссарий

| Термин | Определение |
|--------|------------|
| Заявка (Ticket) | Обращение сотрудника с описанием проблемы или запроса |
| Инцидент | Незапланированный сбой или деградация IT-сервиса |
| Запрос на обслуживание | Стандартный запрос: установка ПО, создание учётки и т.д. |
| SLA | Service Level Agreement — соглашение о времени реакции/решения |
| Статус заявки | Текущее состояние заявки в процессе обработки |
| Приоритет | Степень срочности и критичности заявки |
| Категория | Классификация заявки по типу проблемы |
| Assignee | Специалист, взявший заявку в работу |
| JWT | JSON Web Token — механизм аутентификации |

---

## 3. Роли и прецеденты (Use Cases)

### 3.1 Роль: Сотрудник (Employee)

| ID | Прецедент | Описание |
|----|-----------|---------|
| UC-01 | Регистрация в системе | Создание учётной записи с email и паролем |
| UC-02 | Вход в систему | Аутентификация, получение JWT-токена |
| UC-03 | Создание заявки | Заполнение формы: тема, описание, категория, приоритет |
| UC-04 | Просмотр своих заявок | Список с фильтрацией по статусу и дате |
| UC-05 | Просмотр деталей заявки | История статусов, комментарии специалиста |
| UC-06 | Добавление комментария | Уточнение информации по заявке |
| UC-07 | Оценка решения | Выставление оценки (1-5 звёзд) после закрытия |

### 3.2 Роль: Специалист поддержки (Support)

| ID | Прецедент | Описание |
|----|-----------|---------|
| UC-08 | Просмотр очереди заявок | Все открытые заявки с фильтрами |
| UC-09 | Назначение заявки на себя | Взять в работу |
| UC-10 | Смена статуса заявки | New → In Progress → Waiting → Resolved → Closed |
| UC-11 | Добавление комментария | Запрос доп. информации, уведомление о решении |
| UC-12 | Закрытие заявки | После подтверждения решения |

### 3.3 Роль: Администратор (Admin)

| ID | Прецедент | Описание |
|----|-----------|---------|
| UC-13 | Управление пользователями | Создание, блокировка, смена роли |
| UC-14 | Управление категориями | Добавление/редактирование/удаление категорий |
| UC-15 | Управление приоритетами | Настройка уровней приоритета и цветовых меток |
| UC-16 | Просмотр дашборда | Статистика по заявкам, средние времена |

---

## 4. Схема базы данных (ERD)

### Таблицы

```sql
-- Пользователи
users
  id          UUID PRIMARY KEY
  email       VARCHAR(255) UNIQUE NOT NULL
  password    VARCHAR(255) NOT NULL        -- bcrypt hash
  full_name   VARCHAR(255) NOT NULL
  role        ENUM('employee','support','admin') NOT NULL
  is_active   BOOLEAN DEFAULT true
  created_at  TIMESTAMP
  updated_at  TIMESTAMP

-- Категории заявок
categories
  id          SERIAL PRIMARY KEY
  name        VARCHAR(100) NOT NULL
  description TEXT
  created_at  TIMESTAMP

-- Приоритеты
priorities
  id          SERIAL PRIMARY KEY
  name        VARCHAR(50) NOT NULL         -- Low, Medium, High, Critical
  color       VARCHAR(7)                  -- hex: #FF0000
  sla_hours   INTEGER                     -- время реакции в часах

-- Заявки
tickets
  id          UUID PRIMARY KEY
  number      SERIAL UNIQUE               -- человекочитаемый номер: #1001
  title       VARCHAR(255) NOT NULL
  description TEXT NOT NULL
  status      ENUM('new','in_progress','waiting','resolved','closed') DEFAULT 'new'
  author_id   UUID REFERENCES users(id)
  assignee_id UUID REFERENCES users(id) NULLABLE
  category_id INTEGER REFERENCES categories(id)
  priority_id INTEGER REFERENCES priorities(id)
  rating      SMALLINT NULLABLE           -- оценка 1-5
  created_at  TIMESTAMP
  updated_at  TIMESTAMP
  closed_at   TIMESTAMP NULLABLE

-- Комментарии
comments
  id          UUID PRIMARY KEY
  ticket_id   UUID REFERENCES tickets(id) ON DELETE CASCADE
  author_id   UUID REFERENCES users(id)
  body        TEXT NOT NULL
  is_internal BOOLEAN DEFAULT false       -- внутренние заметки (только support видит)
  created_at  TIMESTAMP

-- История статусов
status_history
  id          UUID PRIMARY KEY
  ticket_id   UUID REFERENCES tickets(id) ON DELETE CASCADE
  changed_by  UUID REFERENCES users(id)
  old_status  VARCHAR(20)
  new_status  VARCHAR(20) NOT NULL
  changed_at  TIMESTAMP DEFAULT NOW()
```

### Связи

```
users ──< tickets (author_id)
users ──< tickets (assignee_id)
categories ──< tickets
priorities ──< tickets
tickets ──< comments
tickets ──< status_history
users ──< comments
users ──< status_history
```

---

## 5. REST API — описание эндпоинтов

### Аутентификация

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| POST | `/api/v1/auth/register` | Регистрация | — |
| POST | `/api/v1/auth/login` | Вход, получение JWT | — |
| POST | `/api/v1/auth/refresh` | Обновление access-токена | authenticated |
| POST | `/api/v1/auth/logout` | Инвалидация refresh-токена | authenticated |

### Заявки (Tickets)

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| GET | `/api/v1/tickets` | Список заявок (с фильтрами) | all |
| POST | `/api/v1/tickets` | Создание заявки | employee |
| GET | `/api/v1/tickets/:id` | Детали заявки | all |
| PATCH | `/api/v1/tickets/:id/status` | Смена статуса | support, admin |
| PATCH | `/api/v1/tickets/:id/assign` | Назначить на себя | support |
| PATCH | `/api/v1/tickets/:id/rate` | Оценить решение | employee (author) |

### Комментарии

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| GET | `/api/v1/tickets/:id/comments` | Список комментариев | all |
| POST | `/api/v1/tickets/:id/comments` | Добавить комментарий | all |

### Справочники

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| GET | `/api/v1/categories` | Список категорий | all |
| POST | `/api/v1/categories` | Создать категорию | admin |
| PUT | `/api/v1/categories/:id` | Обновить категорию | admin |
| DELETE | `/api/v1/categories/:id` | Удалить категорию | admin |
| GET | `/api/v1/priorities` | Список приоритетов | all |

### Пользователи (admin)

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| GET | `/api/v1/users` | Список пользователей | admin |
| PATCH | `/api/v1/users/:id/role` | Сменить роль | admin |
| PATCH | `/api/v1/users/:id/block` | Заблокировать/разблокировать | admin |

### Статистика

| Метод | URL | Описание | Роли |
|-------|-----|---------|------|
| GET | `/api/v1/stats/dashboard` | Сводные показатели | admin, support |

---

## 6. Структура Vue-приложения

### Страницы (pages/)

```
pages/
├── auth/
│   ├── LoginPage.vue          # Форма входа
│   └── RegisterPage.vue       # Форма регистрации
├── tickets/
│   ├── TicketListPage.vue     # Список заявок
│   ├── TicketCreatePage.vue   # Создание заявки
│   └── TicketDetailPage.vue   # Детали заявки + комментарии
├── support/
│   └── QueuePage.vue          # Очередь заявок для support
├── admin/
│   ├── UsersPage.vue          # Управление пользователями
│   └── CategoriesPage.vue     # Управление категориями
└── DashboardPage.vue          # Главная / дашборд
```

### Компоненты (components/)

```
components/
├── layout/
│   ├── AppHeader.vue
│   ├── AppSidebar.vue
│   └── AppLayout.vue
├── tickets/
│   ├── TicketCard.vue         # Карточка заявки в списке
│   ├── TicketStatusBadge.vue  # Цветной бейдж статуса
│   ├── TicketPriorityTag.vue  # Метка приоритета
│   └── TicketFilters.vue      # Панель фильтров
├── comments/
│   ├── CommentList.vue
│   └── CommentForm.vue
└── common/
    ├── LoadingSpinner.vue
    ├── EmptyState.vue
    └── ConfirmModal.vue
```

### Pinia-сторы (stores/)

```
stores/
├── auth.ts      # Текущий пользователь, токены, login/logout
├── tickets.ts   # Список заявок, текущая заявка, CRUD
├── comments.ts  # Комментарии к текущей заявке
└── admin.ts     # Данные для административных страниц
```

---

## 7. Формат JWT-токенов

### Access Token payload

```json
{
  "sub": "uuid пользователя",
  "email": "user@example.com",
  "role": "employee",
  "exp": 1234567890
}
```

### Заголовок запроса

```
Authorization: Bearer <access_token>
```

---

## 8. Диаграмма состояний заявки

```
         [Создана]
             │
          new (●)
             │  support берёт в работу
             ▼
        in_progress
         /        \
        ↓          ↓
   waiting       resolved
  (ждём юзера)  (решена)
        \          │
         \         │ юзер подтверждает / время вышло
          ↓        ▼
              closed (■)
```

Разрешённые переходы:

| Из | В | Кто |
|----|---|-----|
| new | in_progress | support |
| in_progress | waiting | support |
| in_progress | resolved | support |
| waiting | in_progress | support |
| waiting | resolved | support |
| resolved | closed | employee (или авто через 3 дня) |
| resolved | in_progress | employee (если не доволен) |
