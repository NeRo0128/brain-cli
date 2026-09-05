# 🔄 Plan de Migración: brain-cli v1.0 → v2.0

**Fecha:** 2026-09-05  
**Estado:** Propuesto  
**Tipo de cambio:** Arquitectura + Infraestructura

---

## 📊 Resumen Ejecutivo de Cambios

### Cambios Principales

1. **🐳 Contenerización completa con Docker**
2. **💾 SQLite como única fuente de verdad** (adiós YAML para tareas)
3. **🔌 Sistema de proveedores de IA extensible** (no solo OmniRoute)
4. **📦 Scripts embebidos en el binario** (usando `embed`)
5. **🛠️ Nueva entidad: Tool** (separar scripts de tareas)
6. **🗂️ Reestructuración de carpetas**

---

## 🗂️ Comparación de Estructuras

### ❌ Estructura ANTERIOR (v1.0)

```
brain-cli/
├── cmd/tui-assistant/main.go           ← Nombre antiguo
├── internal/
│   ├── core/
│   │   ├── task/
│   │   ├── execution/
│   │   └── omniroute/                  ← Solo OmniRoute
│   ├── usecases/
│   │   ├── task/
│   │   ├── execution/
│   │   └── omniroute/                  ← Específico a OmniRoute
│   ├── adapters/
│   │   ├── cli/                        ← Nombre genérico
│   │   ├── email/
│   │   ├── omniroute/                  ← Cliente específico
│   │   ├── repository/
│   │   │   ├── task_yaml.go           ← ❌ Tareas en YAML
│   │   │   └── execution_sqlite.go    ← Solo executions en DB
│   │   └── config/
│   └── ui/
├── scripts/                             ← Scripts externos
├── configs/
│   ├── config.yaml                      ← Configuración general
│   └── tasks.yaml                       ← ❌ Definición de tareas (YAML)
└── docs/
```

### ✅ Estructura NUEVA (v2.0)

```
brain-cli/
├── cmd/brain-cli/main.go                ← ✅ Nombre simplificado
├── internal/
│   ├── core/
│   │   ├── task/                        ← Mantiene
│   │   ├── execution/                   ← Mantiene
│   │   ├── provider/                    ← 🔥 NUEVO: Proveedores IA
│   │   └── tool/                        ← 🔥 NUEVO: Herramientas/Scripts
│   ├── usecases/
│   │   ├── task/                        ← Mantiene
│   │   ├── execution/                   ← Mantiene
│   │   ├── provider/                    ← 🔥 NUEVO: Gestión proveedores
│   │   └── ai/                          ← ✅ Genérico (no solo OmniRoute)
│   ├── adapters/
│   │   ├── database/                    ← 🔥 NUEVO: SQLite + migrations
│   │   │   ├── sqlite.go
│   │   │   ├── migrations/
│   │   │   └── repositories/           ← ✅ Todos los repos en SQLite
│   │   ├── executor/                    ← ✅ Renombrado (antes "cli")
│   │   │   ├── bash.go
│   │   │   ├── python.go
│   │   │   └── native.go
│   │   ├── ai/                          ← 🔥 NUEVO: Múltiples proveedores
│   │   │   ├── provider.go             ← Interfaz común
│   │   │   ├── omniroute.go
│   │   │   ├── ollama.go               ← Soporte local
│   │   │   └── direct.go               ← OpenAI, DeepSeek
│   │   ├── email/                       ← Mantiene
│   │   └── config/                      ← Mantiene
│   └── ui/
│       ├── screens/
│       │   ├── providers.go             ← 🔥 NUEVO: Gestión proveedores
│       │   └── ...
│       └── ...
├── embed/                                ← 🔥 NUEVO: Archivos embebidos
│   ├── scripts/                         ← Scripts por defecto
│   └── migrations/                      ← Migraciones SQL
├── data/                                 ← 🔥 NUEVO: Volumen persistente
│   └── brain.db                         ← Base de datos principal
├── scripts/                              ← Scripts personalizados
├── docker/                               ← 🔥 NUEVO: Dockerización
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── entrypoint.sh
└── configs/
    └── config.yaml                       ← Solo config general (no tasks)
```

---

## 🔄 Cambios Específicos por Componente

### 1. **Almacenamiento de Tareas: YAML → SQLite**

#### ❌ ANTES (v1.0)
```yaml
# configs/tasks.yaml
tasks:
  - id: "wifi-vpn"
    name: "Conectar WiFi + VPN"
    type: "script"
    script_path: "scripts/wifi-vpn.sh"
    ...
```

```go
// internal/adapters/repository/task_yaml.go
type YAMLTaskRepository struct {
    filePath string
}

func (r *YAMLTaskRepository) List() ([]task.Task, error) {
    // Lee archivo YAML
    // Parsea y retorna tareas
}
```

#### ✅ DESPUÉS (v2.0)
```sql
-- embed/migrations/001_initial.sql
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL,
    tool_id TEXT REFERENCES tools(id),
    params TEXT, -- JSON
    requires_ai BOOLEAN,
    ai_prompt TEXT,
    tags TEXT, -- JSON array
    priority TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tools (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL, -- bash, python, native
    category TEXT,
    script_path TEXT,
    command TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

```go
// internal/adapters/database/repositories/task_repo.go
type SQLiteTaskRepository struct {
    db *sql.DB
}

func (r *SQLiteTaskRepository) List(ctx context.Context) ([]task.Task, error) {
    query := "SELECT id, name, type, ... FROM tasks ORDER BY priority"
    // Query SQLite y retorna tareas
}
```

**Impacto:**
- ✅ Tareas son dinámicas (CRUD completo desde TUI)
- ✅ No necesita reinicio para añadir tareas
- ✅ Historial de cambios (updated_at)
- ✅ Queries complejas (búsqueda, filtrado)
- ❌ Migrar tareas de YAML a SQL (script de migración)

---

### 2. **Sistema de Proveedores de IA**

#### ❌ ANTES (v1.0)
```go
// internal/adapters/omniroute/client.go
type OmniRouteClient struct {
    apiURL string
    apiKey string
}

// Solo hay un cliente, hardcoded
```

```yaml
# configs/config.yaml
omniroute:
  api_url: "http://localhost:8000/v1"
  api_key: "${OMNIROUTE_API_KEY}"
  model: "claude-3-5-sonnet-20241022"
```

#### ✅ DESPUÉS (v2.0)
```sql
-- Nueva tabla de proveedores
CREATE TABLE providers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- omniroute, ollama, openai, deepseek
    base_url TEXT,
    api_key TEXT, -- Encriptado
    model TEXT,
    is_active BOOLEAN DEFAULT 0,
    config TEXT, -- JSON extra config
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Seed data inicial
INSERT INTO providers VALUES
    ('omniroute-default', 'OmniRoute Local', 'omniroute', 'http://localhost:8000/v1', NULL, 'claude-3-5-sonnet-20241022', 1, '{}', CURRENT_TIMESTAMP),
    ('ollama-local', 'Ollama Local', 'ollama', 'http://localhost:11434', NULL, 'llama2', 0, '{}', CURRENT_TIMESTAMP);
```

```go
// internal/core/provider/entity.go
type Provider struct {
    ID       string
    Name     string
    Type     ProviderType // omniroute | ollama | openai | deepseek
    BaseURL  string
    APIKey   string
    Model    string
    IsActive bool
    Config   map[string]interface{}
}

type ProviderType string

const (
    ProviderTypeOmniRoute ProviderType = "omniroute"
    ProviderTypeOllama    ProviderType = "ollama"
    ProviderTypeOpenAI    ProviderType = "openai"
    ProviderTypeDeepSeek  ProviderType = "deepseek"
)

// internal/adapters/ai/provider.go
type AIProvider interface {
    Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
    Classify(ctx context.Context, content string) (Classification, error)
    Summarize(ctx context.Context, content string) (Summary, error)
    IsAvailable(ctx context.Context) bool
}

// internal/adapters/ai/omniroute.go
type OmniRouteProvider struct {
    provider *provider.Provider
    client   *http.Client
}

// internal/adapters/ai/ollama.go
type OllamaProvider struct {
    provider *provider.Provider
    client   *http.Client
}

// Factory para crear proveedores
func NewProvider(p *provider.Provider) (AIProvider, error) {
    switch p.Type {
    case provider.ProviderTypeOmniRoute:
        return NewOmniRouteProvider(p), nil
    case provider.ProviderTypeOllama:
        return NewOllamaProvider(p), nil
    // ...
    }
}
```

**Impacto:**
- ✅ Multi-proveedor sin recompilar
- ✅ Cambiar proveedor activo desde TUI
- ✅ Probar múltiples proveedores
- ✅ Fallback automático si uno falla
- ✅ Configuración por proveedor

---

### 3. **Entidad Tool: Separación de Scripts**

#### ❌ ANTES (v1.0)
```go
// Tasks contienen directamente el script
type Task struct {
    ID          string
    ScriptPath  string  // Ruta hardcoded
    Command     string  // O comando directo
}
```

#### ✅ DESPUÉS (v2.0)
```go
// internal/core/tool/entity.go
type Tool struct {
    ID          string
    Name        string
    Description string
    Type        ToolType     // bash | python | native
    Category    ToolCategory // network | dev | maintenance
    ScriptPath  string
    Command     string
    CreatedAt   time.Time
}

type ToolType string
const (
    ToolTypeBash   ToolType = "bash"
    ToolTypePython ToolType = "python"
    ToolTypeNative ToolType = "native"
)

// internal/core/task/entity.go
type Task struct {
    ID          string
    Name        string
    Type        TaskType // script | command | ai
    ToolID      *string  // ← Referencia a Tool (FK)
    Params      TaskParams
    RequiresAI  bool
    AIPrompt    string
}

// Relación: Task → Tool
// Una tarea puede usar un tool con parámetros diferentes
```

**Impacto:**
- ✅ Reutilización de scripts (1 tool, N tasks)
- ✅ Gestión independiente de herramientas
- ✅ Versionado de tools
- ✅ Parámetros dinámicos por tarea

---

### 4. **Ejecutores de Scripts**

#### ❌ ANTES (v1.0)
```go
// internal/adapters/cli/executor.go
type CLIExecutor interface {
    RunCommand(ctx context.Context, cmd string) (Output, error)
    RunScript(ctx context.Context, scriptPath string) (Output, error)
}
```

#### ✅ DESPUÉS (v2.0)
```go
// internal/adapters/executor/executor.go
type Executor interface {
    Execute(ctx context.Context, tool *tool.Tool, params map[string]string) (Output, error)
    IsAvailable() bool
}

// internal/adapters/executor/bash.go
type BashExecutor struct {
    shell string // /bin/bash
}

// internal/adapters/executor/python.go
type PythonExecutor struct {
    pythonPath string // /usr/bin/python3
}

// internal/adapters/executor/native.go
type NativeExecutor struct {
    // Ejecuta comandos del sistema directamente
}

// Factory
func NewExecutor(toolType tool.ToolType) (Executor, error) {
    switch toolType {
    case tool.ToolTypeBash:
        return &BashExecutor{shell: "/bin/bash"}, nil
    case tool.ToolTypePython:
        return &PythonExecutor{pythonPath: "/usr/bin/python3"}, nil
    case tool.ToolTypeNative:
        return &NativeExecutor{}, nil
    }
}
```

**Impacto:**
- ✅ Ejecutores especializados por tipo
- ✅ Validación de disponibilidad
- ✅ Inyección de parámetros
- ✅ Mejor manejo de errores por tipo

---

### 5. **Scripts Embebidos (`embed`)**

#### ❌ ANTES (v1.0)
```
scripts/
├── wifi-vpn.sh         ← Scripts externos
├── start-dev.sh        ← No están en el binario
└── cleanup.sh          ← Necesitan estar presentes
```

#### ✅ DESPUÉS (v2.0)
```go
// cmd/brain-cli/main.go
import _ "embed"

//go:embed embed/scripts/*
var scriptsFS embed.FS

//go:embed embed/migrations/*.sql
var migrationsFS embed.FS

func main() {
    // Scripts disponibles aunque no exista carpeta scripts/
    // Binario autosuficiente
}
```

```
embed/
├── scripts/
│   ├── wifi-vpn.sh      ← Embebido en binario
│   ├── start-dev.sh
│   └── cleanup.sh
└── migrations/
    ├── 001_initial.sql   ← Migraciones embebidas
    └── 002_seed_data.sql
```

**Impacto:**
- ✅ Binario autosuficiente (no depende de archivos externos)
- ✅ Scripts siempre disponibles
- ✅ Migraciones portables
- ✅ Instalación más simple (1 binario)

---

### 6. **Dockerización**

#### ❌ ANTES (v1.0)
- Sin Docker
- Dependencias del sistema host
- Configuración manual

#### ✅ DESPUÉS (v2.0)
```dockerfile
# docker/Dockerfile (multi-stage)
FROM golang:1.21-alpine AS builder
# ... build con CGO para SQLite

FROM alpine:latest
# Runtime con bash, python, git, sqlite
COPY --from=builder /app/brain-cli .
# Usuario no-root
USER brain
```

```yaml
# docker/docker-compose.yml
services:
  brain-cli:
    volumes:
      - brain-data:/app/data          # SQLite persistente
      - brain-scripts:/app/scripts    # Scripts custom
    environment:
      - DB_PATH=/app/data/brain.db
```

**Impacto:**
- ✅ Aislamiento de dependencias
- ✅ Portabilidad total
- ✅ Persistencia con volúmenes
- ✅ Reproducibilidad

---

## 📋 Matriz de Cambios por Archivo

| Archivo/Carpeta | Acción | Razón |
|-----------------|--------|-------|
| `cmd/tui-assistant/` | ✏️ Renombrar → `cmd/brain-cli/` | Simplificar nombre |
| `configs/tasks.yaml` | ❌ Eliminar | Tareas ahora en SQLite |
| `internal/core/provider/` | ➕ Crear | Nueva entidad Provider |
| `internal/core/tool/` | ➕ Crear | Nueva entidad Tool |
| `internal/usecases/omniroute/` | ✏️ Renombrar → `usecases/ai/` | Genérico, no específico |
| `internal/usecases/provider/` | ➕ Crear | Gestión de proveedores |
| `internal/adapters/cli/` | ✏️ Renombrar → `adapters/executor/` | Más descriptivo |
| `internal/adapters/database/` | ➕ Crear | Nueva capa de persistencia |
| `internal/adapters/omniroute/` | 🔄 Mover → `adapters/ai/omniroute.go` | Parte de sistema multi-provider |
| `internal/adapters/repository/task_yaml.go` | ❌ Eliminar | Reemplazado por SQLite |
| `internal/adapters/repository/execution_sqlite.go` | 🔄 Mover → `database/repositories/execution_repo.go` | Consolidar repos |
| `scripts/` | ⚠️ Mantener | Ahora para scripts custom (no por defecto) |
| `embed/` | ➕ Crear | Scripts y migraciones embebidas |
| `data/` | ➕ Crear | Directorio para brain.db |
| `docker/` | ➕ Crear | Dockerización completa |

---

## 🚧 Riesgos y Mitigaciones

### Riesgo 1: Pérdida de tareas existentes en YAML
**Mitigación:** Crear script de migración `tasks.yaml → brain.db`

```go
// cmd/brain-cli/migrate.go
func MigrateYAMLTasks(yamlPath string, db *sql.DB) error {
    // Lee tasks.yaml
    // Inserta en tabla tasks
    // Inserta tools asociados
}
```

### Riesgo 2: Compatibilidad con scripts existentes
**Mitigación:** Mantener carpeta `scripts/` para scripts personalizados

### Riesgo 3: Configuración de proveedores inicial
**Mitigación:** Seed data con OmniRoute por defecto activo

### Riesgo 4: Complejidad de Docker para usuarios no técnicos
**Mitigación:** 
- Mantener opción de ejecutar binario standalone
- Documentar ambas opciones
- Makefile con comandos simples

---

## ✅ Ventajas del Nuevo Diseño

1. **Persistencia unificada**: Todo en SQLite (no mezclar YAML + DB)
2. **Multi-proveedor**: Soportar Ollama, OpenAI, DeepSeek sin recompilar
3. **Portabilidad**: Docker + binario embebido
4. **Gestión dinámica**: CRUD de tareas/tools desde TUI
5. **Separación clara**: Tools reutilizables, Tasks son instancias
6. **Escalabilidad**: Fácil añadir nuevos proveedores o ejecutores

---

## 📅 Plan de Ejecución

### Fase 1: Preparación (1 día)
- [ ] Crear branches: `develop`, `feature/sqlite-migration`
- [ ] Actualizar AGENTS.md con nueva arquitectura
- [ ] Actualizar plan de desarrollo
- [ ] Crear nuevas carpetas

### Fase 2: Migraciones (1 día)
- [ ] Crear esquema SQL completo
- [ ] Implementar sistema de migraciones
- [ ] Crear seed data
- [ ] Script de migración YAML → SQLite

### Fase 3: Core + Repositorios (2 días)
- [ ] Nuevas entidades: Provider, Tool
- [ ] Actualizar Task entity
- [ ] Implementar repositorios SQLite
- [ ] Tests unitarios

### Fase 4: Adaptadores (2 días)
- [ ] Refactorizar ejecutores
- [ ] Sistema de proveedores AI
- [ ] Cliente Ollama
- [ ] Tests de integración

### Fase 5: UI (2 días)
- [ ] Pantalla de gestión de proveedores
- [ ] Actualizar pantallas de tareas
- [ ] CRUD de tasks desde TUI

### Fase 6: Dockerización (1 día)
- [ ] Dockerfile multi-stage
- [ ] docker-compose.yml
- [ ] Volúmenes persistentes
- [ ] Tests de contenedor

### Fase 7: Documentación (1 día)
- [ ] Actualizar README
- [ ] Guía de migración
- [ ] Docker usage guide
- [ ] API documentation

**Total estimado:** 10 días

---

## 🤔 Decisiones Pendientes

1. **¿Encriptar API keys en la base de datos?**
   - Opción A: Texto plano (más simple, pero menos seguro)
   - Opción B: Encriptar con clave maestra (más seguro, más complejo)
   - **Recomendación:** Empezar con texto plano, añadir encriptación después

2. **¿Mantener soporte para YAML como backup?**
   - Opción A: Solo SQLite (más simple)
   - Opción B: SQLite + export/import YAML (más flexible)
   - **Recomendación:** Solo SQLite, con export a JSON para backups

3. **¿Ejecutar migraciones automáticas al inicio?**
   - Opción A: Automático (más fácil para usuario)
   - Opción B: Manual con comando (más control)
   - **Recomendación:** Automático en primer run, con flag `--skip-migrations`

4. **¿Nombre de carpeta de comandos?**
   - Opción A: `cmd/brain-cli/` (consistente con nombre proyecto)
   - Opción B: `cmd/brain/` (más corto)
   - **Recomendación:** `cmd/brain-cli/` (evita confusión con carpeta raíz)

---

## 📝 Checklist de Aprobación

Antes de comenzar la implementación, confirmar:

- [ ] Arquitectura nueva está clara
- [ ] Cambios en estructura de carpetas aprobados
- [ ] Plan de migración de datos es viable
- [ ] Riesgos identificados y mitigados
- [ ] Timeline de 10 días es realista
- [ ] Decisiones técnicas tomadas

---

**Siguiente paso:** Una vez aprobado este plan, comenzar con Fase 1 (Preparación).

**Fecha de aprobación:** _____________  
**Aprobado por:** _____________
