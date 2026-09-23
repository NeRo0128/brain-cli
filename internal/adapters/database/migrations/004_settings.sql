-- 004_settings.sql
-- Tabla de overrides de configuración editables desde la UI.

CREATE TABLE IF NOT EXISTS settings (
    key          TEXT     PRIMARY KEY,
    value        TEXT     NOT NULL,
    category     TEXT     NOT NULL DEFAULT 'general',
    is_encrypted INTEGER  NOT NULL DEFAULT 0 CHECK(is_encrypted IN (0, 1)),
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
