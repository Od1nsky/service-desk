-- ============================================================
-- ServiceDesk database schema
-- ============================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ------------------------------------------------------------
-- users
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    full_name   VARCHAR(255) NOT NULL,
    role        VARCHAR(20)  NOT NULL DEFAULT 'employee'
                    CHECK (role IN ('employee', 'support', 'admin')),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- refresh_tokens
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       TEXT         NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- categories
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- priorities
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS priorities (
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(50) NOT NULL,
    color     VARCHAR(7)  NOT NULL,
    sla_hours INTEGER     NOT NULL
);

-- ------------------------------------------------------------
-- tickets
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS tickets (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number       SERIAL UNIQUE,
    title        VARCHAR(255) NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    status       VARCHAR(20)  NOT NULL DEFAULT 'new'
                     CHECK (status IN ('new', 'in_progress', 'waiting', 'resolved', 'closed')),
    author_id    UUID         NOT NULL REFERENCES users(id),
    assignee_id  UUID                  REFERENCES users(id),
    category_id  INTEGER      NOT NULL REFERENCES categories(id),
    priority_id  INTEGER      NOT NULL REFERENCES priorities(id),
    rating       SMALLINT              CHECK (rating BETWEEN 1 AND 5),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    closed_at    TIMESTAMPTZ
);

-- ------------------------------------------------------------
-- comments
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS comments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID         NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    author_id   UUID         NOT NULL REFERENCES users(id),
    body        TEXT         NOT NULL,
    is_internal BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- status_history
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS status_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID         NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    changed_by  UUID         NOT NULL REFERENCES users(id),
    old_status  VARCHAR(20)  NOT NULL,
    new_status  VARCHAR(20)  NOT NULL,
    changed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- Indexes
-- ------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_tickets_author_id    ON tickets(author_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee_id  ON tickets(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status       ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_comments_ticket_id   ON comments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_status_history_ticket_id ON status_history(ticket_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id   ON refresh_tokens(user_id);

-- ------------------------------------------------------------
-- Seed data – priorities
-- ------------------------------------------------------------
INSERT INTO priorities (name, color, sla_hours) VALUES
    ('Low',      '#6c757d', 72),
    ('Medium',   '#0d6efd', 24),
    ('High',     '#fd7e14',  8),
    ('Critical', '#dc3545',  4)
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------
-- Seed data – categories
-- ------------------------------------------------------------
INSERT INTO categories (name, description) VALUES
    ('Hardware', ''),
    ('Software', ''),
    ('Network',  ''),
    ('Access',   ''),
    ('Other',    '')
ON CONFLICT DO NOTHING;
