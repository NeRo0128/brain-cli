# AGENTS.md - Brain CLI

**Proyecto:** Brain CLI - Asistente Personal TUI  
**Versión:** 2.0.0  
**Última actualización:** 2026-09-07  
**Estado:** En desarrollo activo

---

## 📋 Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [Lógica de Negocio](#lógica-de-negocio)
4. [Componentes Principales](#componentes-principales)
5. [Flujos de Trabajo](#flujos-de-trabajo)
6. [Decisiones Técnicas](#decisiones-técnicas)
7. [Comandos Principales](#comandos-principales)
8. [Estructura del Proyecto](#estructura-del-proyecto)

---

## 🎯 Resumen Ejecutivo

Brain CLI es un **asistente personal de línea de comandos** (Terminal User Interface) construido en Go que combina automatización de tareas del sistema con capacidades de inteligencia artificial mediante OmniRoute. Su objetivo es maximizar la productividad del usuario automatizando workflows repetitivos y proporcionando una interfaz elegante y eficiente.

### Arquitectura v2.0 - Cambios
- Dockerización completa y gestión de scripts embebidos
- Nueva estructura de carpetas: separación de proveedores y herramientas
- Uso de SQLite como única fuente de datos

### ¿Qué hace Brain CLI?

1. **Automatiza tareas del sistema**: Conectar WiFi/VPN, levantar stacks de desarrollo, limpiar temporales, backups
2. **Integra IA**: Procesa lenguaje natural, clasifica correos, resume contenido, responde preguntas
3. **Gestiona workflows**: Ejecuta scripts, comandos y operaciones complejas de forma organizada
4. **Proporciona feedback visual**: Interfaz TUI moderna con estado en tiempo real
5. **Persiste historial**: Guarda ejecuciones para análisis y repetición

### Propuesta de Valor

- ⚡ **Rapidez**: Ejecuta tareas complejas con pocos atajos de teclado
- 🎨 **Elegancia**: Interfaz visual moderna en terminal
- 🤖 **Inteligencia**: Capacidades de IA para procesamiento de información
- 🔧 **Extensibilidad**: Añade nuevas tareas fácilmente desde YAML
- 📊 **Visibilidad**: Historial completo de ejecuciones y logs

---

## 🏗️ Arquitectura del Sistema

Brain CLI sigue **Clean Architecture** de Robert C. Martin, organizando el código en capas concéntricas con dependencias unidireccionales hacia el centro.

### Capas y Responsabilidades

```
┌─────────────────────────────────────────────────────────┐
│  UI Layer (internal/ui/)                                │
│  • Bubble Tea TUI                                       │
│  • Pantallas: Main, TaskList, Execution, Settings      │
│  • Componentes: Menu, StatusBar, Spinner               │
│  • Responsabilidad: Presentación y captura de input    │
└─────────────────────────────────────────────────────────┘
                        ↓ depende de
┌─────────────────────────────────────────────────────────┐
│  Use Cases Layer (internal/usecases/)                   │
│  • Task: Create, List, Execute                         │
│  • Execution: Manage, Monitor                          │
│  • AI: Chat, Classify, Summarize                       │
│  • Provider: Manage providers                          │
│  • Responsabilidad: Lógica de aplicación               │
└─────────────────────────────────────────────────────────┘
                        ↓ depende de
┌─────────────────────────────────────────────────────────┐
│  Domain Layer (internal/core/)                          │
│  • Entidades: Task, Execution, Provider, Tool          │
│  • Interfaces: Repository contracts                     │
│  • Responsabilidad: Reglas de negocio puras            │
└─────────────────────────────────────────────────────────┘
                        ↑ implementado por
┌─────────────────────────────────────────────────────────┐
│  Adapters Layer (internal/adapters/)                    │
│  • Executor: Ejecuta comandos/scripts                   │
│  • AI Provider: Multi-provider (OmniRoute, Ollama)     │
│  • Database: SQLite with migrations                    │
│  • Email Client: Integración Gmail                      │
│  • Config Manager: Lee/escribe YAML                     │
│  • Responsabilidad: Implementaciones concretas          │
└─────────────────────────────────────────────────────────┘
```

### Principios de Diseño

1. **Dependency Inversion**: Las capas externas dependen de abstracciones del dominio
2. **Single Responsibility**: Cada componente tiene una única razón para cambiar
3. **Open/Closed**: Abierto a extensión (nuevas tareas), cerrado a modificación
4. **Interface Segregation**: Interfaces pequeñas y específicas
5. **Separation of Concerns**: UI, lógica de negocio e infraestructura separadas

---

## 💼 Lógica de Negocio

### Entidad Central: Task

Una **Task** representa una unidad de trabajo ejecutable por el sistema.

```go
type Task struct {
    ID          string      // Identificador único
    Name        string      // Nombre visible al usuario
    Description string      // Descripción de qué hace
    Type        TaskType    // ScriptTask | CommandTask | AITask
    Command     string      // Comando a ejecutar (si es CommandTask)
    ScriptPath  string      // Ruta al script (si es ScriptTask)
    RequiresAI  bool        // Si necesita OmniRoute
    AIPrompt    string      // Prompt para OmniRoute (si RequiresAI)
    Tags        []string    // Etiquetas para filtrado
    Priority    Priority    // Low | Medium | High
    CreatedAt   time.Time
}
```

#### Tipos de Tareas

1. **ScriptTask**: Ejecuta un script bash/python embebido o local
   - Ejemplo: `wifi-vpn` (script embebido en DB)
   
2. **CommandTask**: Ejecuta un comando del sistema directamente
   - Ejemplo: `docker ps -a`
   
3. **AITask**: Procesa entrada del usuario con IA
   - Ejemplo: Clasificar correos, resumir URL

### Entidad: Provider

Un **Provider** representa un proveedor de IA configurado en el sistema.

```go
type Provider struct {
    ID              int
    Name            string
    Type            ProviderType  // omniroute | ollama | openai | deepseek
    Endpoint        string
    APIKey          string
    Model           string
    IsActive        bool
    Config          string        // JSON con configuración adicional
    CreatedAt       time.Time
    UpdatedAt       *time.Time
}
```

### Entidad: Tool

Un **Tool** representa una herramienta o script ejecutable.

```go
type Tool struct {
    ID             int
    Name           string
    Description    string
    ScriptContent  string
    ScriptType     ToolType    // bash | python | native | go
    Category       Category    // system | dev | ai | utils | custom
    RequiresSudo   bool
    TimeoutSeconds int
    IsBuiltin      bool
    Version        int
    CreatedAt      time.Time
    UpdatedAt      *time.Time
}
```

### Entidad: Execution

Una **Execution** representa una instancia de ejecución de una Task.

```go
type Execution struct {
    ID         string           // UUID único
    TaskID     string           // ID de la tarea ejecutada
    Status     ExecutionStatus  // Pending | Running | Completed | Failed
    Output     string           // Stdout capturado
    Error      error            // Error si falló
    StartedAt  time.Time
    FinishedAt time.Time
    Duration   time.Duration
}
```

#### Estados de Ejecución

```
Pending → Running → Completed
                 ↘ Failed
```

- **Pending**: Ejecución creada pero no iniciada
- **Running**: En progreso, mostrando spinner en UI
- **Completed**: Finalizada exitosamente
- **Failed**: Falló con error

### Reglas de Negocio

1. **Unicidad de Ejecución**: Solo una ejecución activa por vez (por seguridad)
2. **Validación de Tasks**: Toda Task debe tener ID, Name y Type válidos
3. **Timeout**: Las ejecuciones tienen timeout configurable (default: 5 minutos)
4. **Cancelación**: El usuario puede cancelar ejecuciones en Running
5. **Persistencia**: Todas las ejecuciones se guardan en historial
6. **Retry**: Tasks fallidas se pueden re-intentar desde historial

---

## 🔧 Componentes Principales

### 1. Task Executor (internal/usecases/task/execute.go)

**Responsabilidad**: Orquestar la ejecución de una tarea de principio a fin.

**Flujo de Ejecución**:

```
1. ValidateTask(task)
   ↓
2. CreateExecution(task) → status: Pending
   ↓
3. UpdateStatus(execution, Running)
   ↓
4. ExecuteByType(task)
   ├─ ScriptTask → CLIExecutor.RunScript()
   ├─ CommandTask → CLIExecutor.RunCommand()
   └─ AITask → OmniRouteClient.Process()
   ↓
5. CaptureOutput(stdout, stderr)
   ↓
6. UpdateStatus(execution, Completed/Failed)
   ↓
7. SaveToHistory(execution)
   ↓
8. Return Result
```

**Dependencias**:
- `TaskRepository`: Para obtener definición de tarea
- `ExecutionRepository`: Para guardar estado
- `Executor` o `AIProvider`: Según tipo de tarea

### 2. AI Provider (internal/adapters/ai/provider.go)

**Responsabilidad**: Comunicación con múltiples proveedores de IA.

**Operaciones**:

```go
type Provider interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    ChatStreaming(ctx context.Context, req *ChatRequest, handler StreamHandler) error
    IsAvailable(ctx context.Context) bool
    GetType() string
}
```

**Proveedores soportados**:
- `OmniRouteProvider`: Cliente HTTP para OmniRoute API
- `OllamaProvider`: Soporte local con Ollama
- OpenAI, DeepSeek, Anthropic (planificados)

### 3. Executor (internal/adapters/executor/executor.go)

**Responsabilidad**: Ejecutar scripts y comandos del sistema.

**Operaciones**:

```go
type Executor interface {
    Execute(ctx context.Context, scriptContent string, params map[string]string) (*Output, error)
    IsAvailable() bool
    GetType() string
}
```

**Características**:
- Captura de stdout/stderr en tiempo real
- Manejo de señales (SIGINT, SIGTERM)
- Variables de entorno configurables
- Working directory configurable
- Timeout por comando
- Soporte para bash, python, native, go

### 4. AI Provider (internal/adapters/ai/provider.go)

**Responsabilidad**: Comunicación con múltiples proveedores de IA.

**Operaciones**:

```go
type Provider interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    ChatStreaming(ctx context.Context, req *ChatRequest, handler StreamHandler) error
    IsAvailable(ctx context.Context) bool
    GetType() string
}
```

**Proveedores soportados**:
- `OmniRouteProvider`: Cliente HTTP para OmniRoute API
- `OllamaProvider`: Soporte local con Ollama
- OpenAI, DeepSeek, Anthropic (planificados)

### 4. Config Manager (internal/adapters/config/config.go)

**Responsabilidad**: Cargar y persistir configuración desde YAML.

**Archivos de Configuración**:

#### `configs/config.yaml`
```yaml
app:
  name: "Brain CLI"
  version: "2.0.0"
  log_level: "info"

execution:
  max_concurrent: 1
  default_timeout: 300
  history_limit: 1000

ui:
  theme: "dark"
  animations: true
```

#### `docker/.env.example`
```bash
# Database
DB_PATH=/app/data/brain.db

# OmniRoute connection
OMNIROUTE_URL=http://omniroute:20128
OMNIROUTE_API_KEY=your-api-key

# Encryption
ENCRYPTION_KEY=your-encryption-key
```

### 5. TUI Application (internal/ui/app.go)

**Responsabilidad**: Gestionar la interfaz de usuario con Bubble Tea.

**Modelo de Estado**:

```go
type Model struct {
    currentScreen Screen           // main | taskList | execution | settings
    tasks         []task.Task       // Tareas disponibles
    executions    []execution.Execution  // Historial
    currentExec   *execution.Execution   // Ejecución en curso
    taskUseCase   usecases.TaskUseCase
    config        *config.Config
    spinner       spinner.Model
    list          list.Model
}
```

**Comandos Bubble Tea**:
- `LoadTasksCmd`: Carga tareas desde SQLite
- `ExecuteTaskCmd`: Ejecuta una tarea (async)
- `PollExecutionCmd`: Consulta estado de ejecución
- `LoadHistoryCmd`: Carga historial desde DB

---

## 🔄 Flujos de Trabajo

### Flujo 1: Usuario Ejecuta Tarea Manual

```
1. Usuario abre Brain CLI
   ↓
2. UI muestra lista de tareas disponibles
   ↓
3. Usuario selecciona "Conectar WiFi + VPN" (Enter)
   ↓
4. UI envía comando ExecuteTask("wifi-vpn")
   ↓
5. TaskExecutor valida tarea y crea Execution
   ↓
6. UI cambia a pantalla "Execution" con spinner
   ↓
7. CLIExecutor ejecuta scripts/wifi-vpn.sh
   ↓
8. Output se captura en tiempo real
   ↓
9. UI muestra output línea por línea
   ↓
10. Script termina exitosamente
    ↓
11. Execution status → Completed
    ↓
12. UI muestra mensaje "✓ Completado" (3s)
    ↓
13. UI regresa a lista de tareas
```

### Flujo 2: Usuario Consulta IA (OmniRoute)

```
1. Usuario selecciona "Revisar Correos"
   ↓
2. UI solicita input adicional (opcional)
   ↓
3. EmailClient obtiene correos no leídos de Gmail
   ↓
4. Se construye prompt: "Clasifica estos correos: [lista]"
   ↓
5. OmniRouteClient envía request a API
   ↓
6. UI muestra spinner "Consultando IA..."
   ↓
7. OmniRoute procesa y retorna clasificación
   ↓
8. UI renderiza respuesta con Markdown (glamour)
   ↓
9. Usuario puede:
    - Ver detalles (Enter)
    - Copiar respuesta (c)
    - Guardar en archivo (s)
    - Volver (Esc)
```

### Flujo 3: Usuario Consulta Historial

```
1. Usuario presiona 'h' (hotkey de historial)
   ↓
2. UI carga últimas 100 ejecuciones desde SQLite
   ↓
3. UI muestra tabla con:
    - Fecha/hora
    - Nombre de tarea
    - Status (✓ Completed / ✗ Failed)
    - Duración
   ↓
4. Usuario selecciona una ejecución (Enter)
   ↓
5. UI muestra detalles:
    - Output completo
    - Error (si hubo)
    - Metadata
   ↓
6. Usuario puede:
    - Re-ejecutar tarea (r)
    - Copiar output (c)
    - Exportar a JSON (e)
    - Volver (Esc)
```

---

## 🎯 Decisiones Técnicas

### ¿Por qué Clean Architecture?

**Problema**: Proyectos Go suelen mezclar lógica de negocio con infraestructura, dificultando testing y evolución.

**Solución**: Clean Architecture separa concerns en capas, permitiendo:
- Tests unitarios sin dependencias externas
- Cambiar UI (TUI → web API) sin tocar lógica
- Cambiar DB (SQLite → PostgreSQL) sin tocar domain
- Facilita onboarding de nuevos developers

**Tradeoff**: Más archivos y abstracciones, pero mayor mantenibilidad a largo plazo.

### ¿Por qué Bubble Tea?

**Alternativas evaluadas**:
- `tview`: Más orientado a widgets, menos control fino
- `termui`: Basado en eventos, menos activo
- `gocui`: Bajo nivel, más complejo

**Razones para Bubble Tea**:
1. **Modelo declarativo**: Similar a React, fácil de razonar
2. **Ecosistema**: lipgloss (estilos), bubbles (componentes), glamour (markdown)
3. **Comunidad activa**: Mantenido por Charm, actualizaciones frecuentes
4. **Concurrencia**: Manejo robusto de operaciones async
5. **Testing**: Fácil de testear gracias a modelo funcional

### ¿Por qué YAML para Configuración?

**Alternativas**: JSON, TOML, HCL

**Razones**:
- Soporta comentarios (importante para documentar)
- Legible por humanos
- Ampliamente adoptado (Kubernetes, Docker Compose, etc.)
- Fácil de editar manualmente

### ¿Por qué SQLite para Historial?

**Alternativas**: JSON files, PostgreSQL, BoltDB

**Razones**:
- No requiere servidor externo (embebido)
- Excelente para datos locales (<1GB)
- Soporta búsqueda y filtrado SQL complejo
- Portable (un solo archivo)
- Maduro y estable

---

## ⚙️ Comandos Principales

### Build y Ejecución

```bash
# Compilar
go build -o brain-cli cmd/tui-assistant/main.go

# Ejecutar en desarrollo
go run cmd/tui-assistant/main.go

# Ejecutar con debug logging
DEBUG=1 go run cmd/tui-assistant/main.go

# Instalar dependencias
go mod download
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/usecases/task/...

# Run with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

### Linting y Calidad

```bash
# Run go vet
go vet ./...

# Run golangci-lint
golangci-lint run

# Format code
go fmt ./...

# Check imports
goimports -w .

# Run staticcheck
staticcheck ./...
```

### Docker

```bash
# Build image
docker build -f docker/Dockerfile -t brain-cli .

# Run with docker-compose
cd docker && docker-compose up -d

# View logs
docker-compose logs -f brain-cli
```

### Utilidades

```bash
# Generate mocks (usando mockery)
mockery --all --output internal/mocks

# Update dependencies
go get -u ./...
go mod tidy

# Vendor dependencies (opcional)
go mod vendor

# Clean build cache
go clean -cache
```

---

## 📁 Estructura del Proyecto

```
brain-cli/
├── cmd/
│   └── brain-cli/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── core/                       # Domain Layer
│   │   ├── task/                   # Task entity & repository interface
│   │   ├── execution/              # Execution entity & repository interface
│   │   ├── provider/               # Provider entity (AI providers)
│   │   └── tool/                   # Tool entity (scripts)
│   │
│   ├── usecases/                   # Application Layer
│   │   ├── task/                   # Task CRUD & execution
│   │   ├── execution/              # Execution management
│   │   ├── ai/                     # AI operations (chat, classify, summarize)
│   │   └── provider/               # Provider management
│   │
│   ├── adapters/                   # Infrastructure Layer
│   │   ├── database/               # SQLite & migrations
│   │   │   ├── sqlite.go           # Database connection
│   │   │   ├── migrations/         # SQL migrations
│   │   │   └── repositories/       # Repository implementations
│   │   ├── executor/               # Script/command executors
│   │   │   ├── bash.go             # Bash executor
│   │   │   ├── python.go           # Python executor
│   │   │   └── native.go           # Native commands
│   │   ├── ai/                     # AI providers
│   │   │   ├── provider.go         # Provider interface
│   │   │   ├── omniroute.go        # OmniRoute provider
│   │   │   └── ollama.go           # Ollama provider
│   │   ├── email/                  # Gmail integration
│   │   └── config/                 # Configuration manager
│   │
│   └── ui/                         # Presentation Layer
│       ├── app.go                  # Main Bubble Tea application
│       ├── screens/
│       │   ├── main.go             # Main menu screen
│       │   ├── task_list.go        # Task list screen
│       │   ├── execution.go        # Execution details screen
│       │   ├── history.go          # History screen
│       │   ├── settings.go         # Settings screen
│       │   └── providers.go        # Provider management screen
│       ├── components/
│       │   ├── menu.go             # Menu component
│       │   ├── status_bar.go       # Status bar component
│       │   ├── spinner.go          # Loading spinner
│       │   └── table.go            # Table component
│       └── styles/
│           └── theme.go            # lipgloss styles
│
├── embed/                          # Embedded files
│   ├── scripts/                    # Default scripts (embedded in binary)
│   └── migrations/                 # SQL migrations (embedded)
│
├── pkg/                            # Public packages (reusable)
│   └── utils/
│       ├── logger.go               # Structured logger
│       ├── validator.go            # Input validation
│       └── formatter.go            # Output formatting
│
├── scripts/                        # Custom scripts (user-defined)
│
├── docker/                         # Docker infrastructure
│   ├── Dockerfile                  # Multi-stage build
│   ├── docker-compose.yml          # Container orchestration
│   ├── entrypoint.sh               # Initialization script
│   └── .env.example                # Environment variables example
│
├── data/                           # Persistent data (SQLite DB)
│
├── configs/
│   └── config.yaml                 # Main configuration (no tasks)
│
├── docs/
│   ├── GITHUB_ISSUES.md            # GitHub issues template
│   ├── MIGRATION_PLAN.md           # v1→v2 migration plan
│   └── plans/
│       └── development-roadmap.md  # Development plan
│
├── go.mod                          # Go modules
├── go.sum                          # Dependency checksums
├── AGENTS.md                       # This file
├── README.md                       # Project readme
├── DEVELOPMENT.md                  # Development guide
└── Makefile                        # Build automation
```

---

## 🔐 Seguridad y Buenas Prácticas

### Manejo de Secrets

- ❌ **NO** hardcodear API keys en código
- ✅ Usar variables de entorno: `${OMNIROUTE_API_KEY}`
- ✅ Archivo `.env` para desarrollo (gitignored)
- ✅ Validar que secrets existan al inicio

### Manejo de Errores

```go
// ✅ BIEN: Errores contextuales
if err != nil {
    return fmt.Errorf("failed to execute task %s: %w", taskID, err)
}

// ❌ MAL: Errores genéricos
if err != nil {
    return err
}
```

### Logging

- **Debug**: Información detallada para desarrollo
- **Info**: Eventos importantes (task started, completed)
- **Warn**: Situaciones inesperadas pero manejadas
- **Error**: Fallos que impiden continuar

### Concurrencia

- Usar `context.Context` para cancelación
- Proteger estado compartido con `sync.Mutex`
- Preferir canales a variables compartidas
- Siempre cerrar goroutines correctamente

---

## 📊 Métricas y Observabilidad

### Logs Estructurados

```json
{
  "level": "info",
  "timestamp": "2026-09-05T15:30:00Z",
  "task_id": "wifi-vpn",
  "execution_id": "550e8400-e29b-41d4-a716-446655440000",
  "duration_ms": 1250,
  "status": "completed",
  "message": "Task executed successfully"
}
```

### Historial de Ejecuciones

Guardado en SQLite para análisis:
- Tasa de éxito por tarea
- Tiempo promedio de ejecución
- Errores más frecuentes
- Uso de recursos

---

## 🚀 Próximos Pasos

1. ✅ Estructura de carpetas creada
2. ✅ Plan de desarrollo documentado
3. ✅ go.mod con dependencias configurado
4. ✅ Infrastructure Docker completada
5. ✅ Entidades Core (Provider, Tool) implementadas
6. ⏳ Implementar repositorios SQLite
7. ⏳ Implementar UI con Bubble Tea
8. ⏳ Implementar adaptadores de ejecución
9. ⏳ Crear issues en GitHub
10. ⏳ Implementar Hito 1: Fundación

---

**Mantenedores**: Nero  
**Licencia**: MIT  
**Repositorio**: https://github.com/NeRo0128/brain-cli
