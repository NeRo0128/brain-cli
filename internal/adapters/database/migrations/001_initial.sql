-- ============================================
-- brain-cli Database Schema v2.0
-- Migration: 001_initial.sql
-- Description: Schema alineado con las entidades del dominio
-- Date: 2026-09-11
-- ============================================

-- ============================================
-- Table: providers
-- ============================================
CREATE TABLE providers (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT    NOT NULL UNIQUE,
    type            TEXT    NOT NULL CHECK(type IN ('omniroute', 'ollama', 'openai', 'deepseek')),
    base_url        TEXT    NOT NULL DEFAULT '',
    model           TEXT    NOT NULL DEFAULT '',
    description     TEXT    NOT NULL DEFAULT '',
    api_key_env     TEXT    NOT NULL DEFAULT '',
    api_key         TEXT    NOT NULL DEFAULT '',
    is_active       INTEGER NOT NULL DEFAULT 0 CHECK(is_active IN (0, 1)),
    config          TEXT    NOT NULL DEFAULT '{}',
    max_tokens      INTEGER NOT NULL DEFAULT 0,
    timeout_seconds INTEGER NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME,
    deleted_at      DATETIME
);

-- Invariante: como máximo UN provider activo (partial unique index).
CREATE UNIQUE INDEX idx_providers_active ON providers(is_active) WHERE is_active = 1;
CREATE INDEX idx_providers_type ON providers(type);

-- ============================================
-- Table: tools
-- ============================================
CREATE TABLE tools (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT    NOT NULL UNIQUE,
    description     TEXT    NOT NULL DEFAULT '',
    script_type     TEXT    NOT NULL CHECK(script_type IN ('bash', 'python', 'native', 'go')),
    category        TEXT    NOT NULL DEFAULT 'custom'
                            CHECK(category IN ('system', 'dev', 'ai', 'utils', 'custom', 'network', 'maintenance')),
    script_content  TEXT    NOT NULL DEFAULT '',
    script_path     TEXT    NOT NULL DEFAULT '',
    command         TEXT    NOT NULL DEFAULT '',
    requires_sudo   INTEGER NOT NULL DEFAULT 0 CHECK(requires_sudo IN (0, 1)),
    timeout_seconds INTEGER NOT NULL DEFAULT 300,
    is_builtin      INTEGER NOT NULL DEFAULT 0 CHECK(is_builtin IN (0, 1)),
    version         INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME,
    deleted_at      DATETIME,

    -- Consistencia: una y solo una fuente según el tipo.
    CHECK (
        (script_type IN ('bash', 'python', 'go')
            AND ((script_content != '' AND script_path = '') OR (script_content = '' AND script_path != ''))
            AND command = '')
        OR
        (script_type = 'native'
            AND command != ''
            AND script_content = ''
            AND script_path = '')
    )
);

CREATE INDEX idx_tools_category ON tools(category);
CREATE INDEX idx_tools_type     ON tools(script_type);
CREATE INDEX idx_tools_builtin  ON tools(is_builtin);

-- ============================================
-- Table: tasks
-- ============================================
CREATE TABLE tasks (
    id           TEXT     PRIMARY KEY,                          -- slug
    name         TEXT     NOT NULL UNIQUE,
    description  TEXT     NOT NULL DEFAULT '',
    type         TEXT     NOT NULL CHECK(type IN ('script', 'command', 'ai')),
    tool_id      INTEGER  REFERENCES tools(id) ON DELETE RESTRICT,
    requires_ai  INTEGER  NOT NULL DEFAULT 0 CHECK(requires_ai IN (0, 1)),
    ai_prompt    TEXT     NOT NULL DEFAULT '',
    params       TEXT     NOT NULL DEFAULT '{}',                -- JSON
    priority     TEXT     NOT NULL DEFAULT 'medium' CHECK(priority IN ('low', 'medium', 'high')),
    is_active    INTEGER  NOT NULL DEFAULT 1 CHECK(is_active IN (0, 1)),
    is_favorite  INTEGER  NOT NULL DEFAULT 0 CHECK(is_favorite IN (0, 1)),
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME,
    deleted_at   DATETIME,

    -- Consistencia: type ↔ tool_id y type ↔ ai_prompt.
    CHECK (
        (type IN ('script', 'command') AND tool_id IS NOT NULL AND ai_prompt = '')
        OR
        (type = 'ai' AND tool_id IS NULL AND ai_prompt != '')
    )
);

CREATE INDEX idx_tasks_tool_id  ON tasks(tool_id);
CREATE INDEX idx_tasks_active   ON tasks(is_active);
CREATE INDEX idx_tasks_favorite ON tasks(is_favorite);
CREATE INDEX idx_tasks_priority ON tasks(priority);

-- ============================================
-- Table: tags
-- ============================================
CREATE TABLE tags (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL UNIQUE,
    color      TEXT    NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- Table: task_tags (pivote)
-- ============================================
CREATE TABLE task_tags (
    task_id TEXT    NOT NULL,
    tag_id  INTEGER NOT NULL,
    PRIMARY KEY (task_id, tag_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id)  REFERENCES tags(id)  ON DELETE CASCADE
);

CREATE INDEX idx_task_tags_tag ON task_tags(tag_id);

-- ============================================
-- Table: executions
-- ============================================
CREATE TABLE executions (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id      TEXT    NOT NULL,
    status       TEXT    NOT NULL CHECK(status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    triggered_by TEXT    NOT NULL DEFAULT 'manual' CHECK(triggered_by IN ('manual', 'schedule', 'chain')),
    output       TEXT    NOT NULL DEFAULT '',
    error        TEXT    NOT NULL DEFAULT '',
    exit_code    INTEGER,
    provider_id  INTEGER REFERENCES providers(id) ON DELETE SET NULL,
    params       TEXT    NOT NULL DEFAULT '{}',                 -- JSON
    started_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at  DATETIME,

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_executions_task_id         ON executions(task_id);
CREATE INDEX idx_executions_status          ON executions(status);
CREATE INDEX idx_executions_started         ON executions(started_at);
CREATE INDEX idx_executions_task_started    ON executions(task_id, started_at DESC);
