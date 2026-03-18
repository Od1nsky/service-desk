-- ============================================================
-- АИС контроля текущей успеваемости студентов — database schema
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
    role        VARCHAR(20)  NOT NULL DEFAULT 'student'
                    CHECK (role IN ('student', 'teacher', 'admin')),
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
-- disciplines
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS disciplines (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- grade_types
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS grade_types (
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(50) NOT NULL,
    color     VARCHAR(7)  NOT NULL,
    max_score INTEGER     NOT NULL DEFAULT 100
);

-- ------------------------------------------------------------
-- grades
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS grades (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number         SERIAL UNIQUE,
    task_name      VARCHAR(255) NOT NULL,
    notes          TEXT         NOT NULL DEFAULT '',
    status         VARCHAR(20)  NOT NULL DEFAULT 'pending'
                       CHECK (status IN ('pending', 'graded', 'not_certified', 'certified', 'closed')),
    student_id     UUID         NOT NULL REFERENCES users(id),
    teacher_id     UUID                  REFERENCES users(id),
    discipline_id  INTEGER      NOT NULL REFERENCES disciplines(id),
    grade_type_id  INTEGER      NOT NULL REFERENCES grade_types(id),
    score          SMALLINT              CHECK (score BETWEEN 0 AND 100),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    graded_at      TIMESTAMPTZ
);

-- ------------------------------------------------------------
-- comments
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS comments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grade_id    UUID         NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    author_id   UUID         NOT NULL REFERENCES users(id),
    body        TEXT         NOT NULL,
    is_internal BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- grade_history
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS grade_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grade_id    UUID         NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    changed_by  UUID         NOT NULL REFERENCES users(id),
    old_value   VARCHAR(20)  NOT NULL,
    new_value   VARCHAR(20)  NOT NULL,
    changed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------
-- Indexes
-- ------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_grades_student_id    ON grades(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_teacher_id    ON grades(teacher_id);
CREATE INDEX IF NOT EXISTS idx_grades_status        ON grades(status);
CREATE INDEX IF NOT EXISTS idx_comments_grade_id    ON comments(grade_id);
CREATE INDEX IF NOT EXISTS idx_grade_history_grade_id ON grade_history(grade_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- ------------------------------------------------------------
-- Seed data – grade_types
-- ------------------------------------------------------------
INSERT INTO grade_types (name, color, max_score) VALUES
    ('Текущая',     '#0dcaf0', 100),
    ('Рубежная',    '#fd7e14', 100),
    ('Итоговая',    '#0d6efd', 100),
    ('Посещаемость','#198754',   1)
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------
-- Seed data – disciplines
-- ------------------------------------------------------------
INSERT INTO disciplines (name, description) VALUES
    ('Математический анализ',   'Фундаментальный математический курс'),
    ('Базы данных',             'Теория и практика проектирования БД'),
    ('Информационные системы',  'Разработка и анализ ИС'),
    ('Программирование на Go',  'Системное программирование на Go'),
    ('Веб-разработка',          'Фронтенд и бэкенд веб-приложений')
ON CONFLICT DO NOTHING;
