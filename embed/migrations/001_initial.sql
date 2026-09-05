-- ============================================
-- brain-cli Database Schema v2.0
-- ============================================
-- Migration: 001_initial.sql
-- Description: Schema completo para brain-cli v2.0
-- Author: Nero
-- Date: 2026-09-05
-- ============================================

-- ============================================
-- Table: providers (Motores de IA)
-- Description: Almacena la configuración de los diferentes motores de IA
-- ============================================
CREATE TABLE IF NOT EXISTS providers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK(type IN ('omniroute', 'ollama', 'openai', 'deepseek', 'anthropic', 'custom')),
    endpoint TEXT NOT NULL,
    api_key TEXT,  -- Se encriptará antes de guardar
    model TEXT,
    is_active BOOLEAN DEFAULT 0,
    config TEXT,  -- JSON con configuración adicional
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE INDEX idx_providers_type ON providers(type);
CREATE INDEX idx_providers_active ON providers(is_active);

-- ============================================
-- Table: tools (Scripts y Herramientas)
-- Description: Guarda el código de los scripts que brain-cli puede ejecutar
-- ============================================
CREATE TABLE IF NOT EXISTS tools (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    script_content TEXT NOT NULL,
    script_type TEXT NOT NULL CHECK(script_type IN ('bash', 'python', 'native', 'go')),
    category TEXT CHECK(category IN ('system', 'dev', 'ai', 'utils', 'custom', 'network', 'maintenance')),
    requires_sudo BOOLEAN DEFAULT 0,
    timeout_seconds INTEGER DEFAULT 300,
    is_builtin BOOLEAN DEFAULT 0,
    version INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE INDEX idx_tools_category ON tools(category);
CREATE INDEX idx_tools_type ON tools(script_type);
CREATE INDEX idx_tools_builtin ON tools(is_builtin);

-- ============================================
-- Table: tasks (Tareas)
-- Description: El nivel más alto de acción: une herramientas con IA y define qué hacer
-- ============================================
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    tool_id INTEGER NOT NULL,
    provider_id INTEGER,
    prompt_template TEXT,
    parameters TEXT,  -- JSON con parámetros
    schedule TEXT,  -- Expresión CRON para ejecución automática
    is_active BOOLEAN DEFAULT 1,
    is_favorite BOOLEAN DEFAULT 0,
    tags TEXT,  -- JSON array de etiquetas
    last_execution DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers(id) ON DELETE SET NULL
);

CREATE INDEX idx_tasks_tool_id ON tasks(tool_id);
CREATE INDEX idx_tasks_provider_id ON tasks(provider_id);
CREATE INDEX idx_tasks_active ON tasks(is_active);
CREATE INDEX idx_tasks_favorite ON tasks(is_favorite);
CREATE INDEX idx_tasks_schedule ON tasks(schedule);

-- ============================================
-- Table: executions (Historial de Ejecuciones)
-- Description: Registra todo lo que brain-cli ha hecho
-- ============================================
CREATE TABLE IF NOT EXISTS executions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'running', 'success', 'failed', 'cancelled')),
    output TEXT,
    error TEXT,
    exit_code INTEGER,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    duration_ms INTEGER,
    parameters_used TEXT,  -- JSON con parámetros reales usados
    provider_used TEXT,  -- Nombre del proveedor usado
    model_used TEXT,  -- Modelo específico usado
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_executions_task_id ON executions(task_id);
CREATE INDEX idx_executions_status ON executions(status);
CREATE INDEX idx_executions_started ON executions(started_at);
CREATE INDEX idx_executions_task_status ON executions(task_id, status);

-- ============================================
-- Table: settings (Configuración Global)
-- Description: Almacena preferencias del usuario
-- ============================================
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    category TEXT DEFAULT 'general' CHECK(category IN ('general', 'ui', 'ai', 'system')),
    is_encrypted BOOLEAN DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_settings_category ON settings(category);

-- ============================================
-- Table: tags (Etiquetas para Organizar)
-- Description: Permite categorizar tareas de forma flexible
-- ============================================
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    color TEXT,  -- Color en formato hex (#ff0000)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- Table: task_tags (Relación Muchos-a-Muchos)
-- Description: Tabla pivote entre tasks y tags
-- ============================================
CREATE TABLE IF NOT EXISTS task_tags (
    task_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (task_id, tag_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX idx_task_tags_task ON task_tags(task_id);
CREATE INDEX idx_task_tags_tag ON task_tags(tag_id);

-- ============================================
-- Table: execution_artifacts (Artefactos de Ejecución)
-- Description: Guarda archivos o resultados generados durante ejecuciones
-- ============================================
CREATE TABLE IF NOT EXISTS execution_artifacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    execution_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('file', 'json', 'image', 'log', 'html', 'text')),
    content TEXT,
    file_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (execution_id) REFERENCES executions(id) ON DELETE CASCADE
);

CREATE INDEX idx_artifacts_execution ON execution_artifacts(execution_id);
CREATE INDEX idx_artifacts_type ON execution_artifacts(type);

-- ============================================
-- Triggers: Auto-update timestamps
-- ============================================

-- Trigger para providers
CREATE TRIGGER IF NOT EXISTS update_providers_timestamp 
AFTER UPDATE ON providers
FOR EACH ROW
BEGIN
    UPDATE providers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Trigger para tools
CREATE TRIGGER IF NOT EXISTS update_tools_timestamp 
AFTER UPDATE ON tools
FOR EACH ROW
BEGIN
    UPDATE tools SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Trigger para tasks
CREATE TRIGGER IF NOT EXISTS update_tasks_timestamp 
AFTER UPDATE ON tasks
FOR EACH ROW
BEGIN
    UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Trigger para actualizar last_execution en tasks
CREATE TRIGGER IF NOT EXISTS update_task_last_execution
AFTER INSERT ON executions
FOR EACH ROW
BEGIN
    UPDATE tasks 
    SET last_execution = NEW.started_at
    WHERE id = NEW.task_id;
END;

-- Trigger para asegurar solo un proveedor activo
CREATE TRIGGER IF NOT EXISTS ensure_single_active_provider
BEFORE UPDATE ON providers
FOR EACH ROW
WHEN NEW.is_active = 1 AND OLD.is_active = 0
BEGIN
    UPDATE providers SET is_active = 0 WHERE is_active = 1;
END;

-- ============================================
-- Views: Utilidades para consultas comunes
-- ============================================

-- View: Tareas con información completa
CREATE VIEW IF NOT EXISTS v_tasks_full AS
SELECT 
    t.id,
    t.name,
    t.description,
    t.is_active,
    t.is_favorite,
    t.last_execution,
    t.tags,
    t.schedule,
    tool.name AS tool_name,
    tool.script_type AS tool_type,
    tool.category AS tool_category,
    p.name AS provider_name,
    p.type AS provider_type,
    p.model AS provider_model,
    (SELECT COUNT(*) FROM executions e WHERE e.task_id = t.id) AS execution_count,
    (SELECT COUNT(*) FROM executions e WHERE e.task_id = t.id AND e.status = 'success') AS success_count,
    (SELECT COUNT(*) FROM executions e WHERE e.task_id = t.id AND e.status = 'failed') AS failed_count
FROM tasks t
LEFT JOIN tools tool ON t.tool_id = tool.id
LEFT JOIN providers p ON t.provider_id = p.id;

-- View: Últimas ejecuciones con detalles
CREATE VIEW IF NOT EXISTS v_recent_executions AS
SELECT 
    e.id AS execution_id,
    e.status,
    e.started_at,
    e.finished_at,
    e.duration_ms,
    e.exit_code,
    t.id AS task_id,
    t.name AS task_name,
    tool.name AS tool_name,
    tool.script_type AS tool_type,
    e.provider_used,
    e.model_used
FROM executions e
JOIN tasks t ON e.task_id = t.id
JOIN tools tool ON t.tool_id = tool.id
ORDER BY e.started_at DESC;

-- View: Estadísticas de proveedores
CREATE VIEW IF NOT EXISTS v_provider_stats AS
SELECT 
    p.id,
    p.name,
    p.type,
    p.is_active,
    p.model,
    COUNT(DISTINCT e.task_id) AS tasks_using,
    COUNT(e.id) AS total_executions,
    SUM(CASE WHEN e.status = 'success' THEN 1 ELSE 0 END) AS successful_executions,
    SUM(CASE WHEN e.status = 'failed' THEN 1 ELSE 0 END) AS failed_executions,
    AVG(e.duration_ms) AS avg_duration_ms
FROM providers p
LEFT JOIN tasks t ON p.id = t.provider_id
LEFT JOIN executions e ON t.id = e.task_id AND e.provider_used = p.name
GROUP BY p.id;

-- View: Estadísticas de tools
CREATE VIEW IF NOT EXISTS v_tool_stats AS
SELECT 
    tool.id,
    tool.name,
    tool.script_type,
    tool.category,
    tool.is_builtin,
    COUNT(DISTINCT t.id) AS tasks_using,
    COUNT(e.id) AS total_executions,
    SUM(CASE WHEN e.status = 'success' THEN 1 ELSE 0 END) AS successful_executions,
    SUM(CASE WHEN e.status = 'failed' THEN 1 ELSE 0 END) AS failed_executions,
    AVG(e.duration_ms) AS avg_duration_ms,
    MAX(e.started_at) AS last_used
FROM tools tool
LEFT JOIN tasks t ON tool.id = t.tool_id
LEFT JOIN executions e ON t.id = e.task_id
GROUP BY tool.id;

-- View: Tareas favoritas ordenadas por uso
CREATE VIEW IF NOT EXISTS v_favorite_tasks AS
SELECT 
    t.id,
    t.name,
    t.description,
    t.is_active,
    t.last_execution,
    tool.name AS tool_name,
    COUNT(e.id) AS execution_count
FROM tasks t
JOIN tools tool ON t.tool_id = tool.id
LEFT JOIN executions e ON t.id = e.task_id
WHERE t.is_favorite = 1
GROUP BY t.id
ORDER BY execution_count DESC, t.last_execution DESC;

-- ============================================
-- Funciones de utilidad (comentadas como referencia)
-- ============================================

-- SQLite no soporta funciones personalizadas nativamente,
-- pero estas operaciones se implementarán en Go:
--
-- - EncryptAPIKey(plaintext string) string
-- - DecryptAPIKey(encrypted string) string
-- - ParseJSON(json_text string) map[string]interface{}
-- - FormatDuration(ms int) string
-- - CalculateSuccessRate(task_id int) float64

-- ============================================
-- Índices adicionales para performance
-- ============================================

-- Índices compuestos para queries comunes
CREATE INDEX IF NOT EXISTS idx_executions_task_started ON executions(task_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_active_favorite ON tasks(is_active, is_favorite);

-- ============================================
-- Constraints adicionales (vía triggers)
-- ============================================

-- Validar que prompt_template esté presente si provider_id no es NULL
CREATE TRIGGER IF NOT EXISTS validate_ai_task
BEFORE INSERT ON tasks
FOR EACH ROW
WHEN NEW.provider_id IS NOT NULL AND (NEW.prompt_template IS NULL OR NEW.prompt_template = '')
BEGIN
    SELECT RAISE(ABORT, 'Tasks con provider_id deben tener prompt_template');
END;

-- Validar que tool_id sea válido
CREATE TRIGGER IF NOT EXISTS validate_tool_exists
BEFORE INSERT ON tasks
FOR EACH ROW
WHEN NEW.tool_id NOT IN (SELECT id FROM tools)
BEGIN
    SELECT RAISE(ABORT, 'tool_id debe referenciar un tool existente');
END;

-- ============================================
-- Pragma settings para optimización
-- ============================================

-- Estas se ejecutarán al conectar a la DB desde Go:
-- PRAGMA foreign_keys = ON;
-- PRAGMA journal_mode = WAL;
-- PRAGMA synchronous = NORMAL;
-- PRAGMA cache_size = -64000;  -- 64MB cache
-- PRAGMA temp_store = MEMORY;
-- PRAGMA mmap_size = 30000000000;  -- 30GB memory-mapped I/O

-- ============================================
-- End of Migration 001
-- ============================================
