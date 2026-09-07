# 📝 Brain CLI - Proyecto y Estado

**Proyecto:** Brain CLI - Asistente Personal TUI  
**Última actualización:** 2026-09-07
**Estado:** En desarrollo activo

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

### ¿Qué hace Brain CLI?

1. **Automatiza tareas del sistema**: Conectar WiFi/VPN, levantar stacks de desarrollo, limpiar temporales, backups.
2. **Integra IA**: Procesa lenguaje natural, clasifica correos, resume contenido, responde preguntas.
3. **Gestiona workflows**: Ejecuta scripts, comandos y operaciones complejas de forma organizada.
4. **Proporciona feedback visual**: Interfaz TUI moderna con estado en tiempo real.
5. **Persiste historial**: Guarda ejecuciones para análisis y repetición.

### Propuesta de Valor

- ⚡ **Rapidez**: Ejecuta tareas complejas con pocos atajos de teclado.
- 🎨 **Elegancia**: Interfaz visual moderna en terminal.
- 🤖 **Inteligencia**: Capacidades de IA para procesamiento de información.
- 🔧 **Extensibilidad**: Añade nuevas tareas fácilmente desde YAML.
- 📊 **Visibilidad**: Historial completo de ejecuciones y logs.

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
│  • OmniRoute: Chat, Classify, Summarize                │
│  • Responsabilidad: Lógica de aplicación               │
└─────────────────────────────────────────────────────────┘
                        ↓ depende de
┌─────────────────────────────────────────────────────────┐
│  Domain Layer (internal/core/)                          │
│  • Entidades: Task, Execution, OmniRouteRequest        │
│  • Interfaces: Repository contracts                     │
│  • Responsabilidad: Reglas de negocio puras            │
└─────────────────────────────────────────────────────────┘
                        ↑ implementado por
┌─────────────────────────────────────────────────────────┐
│  Adapters Layer (internal/adapters/)                    │
│  • CLI Executor: Ejecuta comandos del sistema          │
│  • OmniRoute Client: HTTP client para API              │
│  • Email Client: Integración Gmail                      │
│  • Config Manager: Lee/escribe YAML                     │
│  • Responsabilidad: Implementaciones concretas          │
└─────────────────────────────────────────────────────────┘
```

### Principios de Diseño

1. **Dependency Inversion**: Las capas externas dependen de abstracciones del dominio.
2. **Single Responsibility**: Cada componente tiene una única razón para cambiar.
3. **Open/Closed**: Abierto a extensión (nuevas tareas), cerrado a modificación.
4. **Interface Segregation**: Interfaces pequeñas y específicas.
5. **Separation of Concerns**: UI, lógica de negocio e infraestructura separadas.

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
1. **ScriptTask**: Ejecuta un script bash/python en `scripts/`
2. **CommandTask**: Ejecuta un comando del sistema directamente
3. **AITask**: Procesa entrada del usuario con OmniRoute

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

---

## 🔧 Componentes Principales

### 1. Task Executor (internal/usecases/task/execute.go)
**Responsabilidad**: Orquestar la ejecución de una tarea de principio a fin.
### 2. OmniRoute Client (internal/adapters/omniroute/client.go)
**Responsabilidad**: Comunicación HTTP con la API de OmniRoute.
### 3. CLI Executor (internal/adapters/cli/executor.go)
**Responsabilidad**: Ejecutar comandos del sistema y scripts.
### 4. Config Manager (internal/adapters/config/config.go)
**Responsabilidad**: Cargar y persistir configuración desde YAML.
### 5. UI Application (internal/ui/app.go)
**Responsabilidad**: Gestionar la interfaz de usuario.

---

## ⚙️ Comandos Principales

```bash
# Compilar
go build -o brain-cli cmd/tui-assistant/main.go
```

---

## 📁 Estructura del Proyecto

```
brain-cli/
├── cmd/
│   └── brain-cli/
│       └── main.go
├── internal/
│   ├── core/
│   ├── usecases/
│   ├── adapters/
│   └── ui/
├── scripts/
├── configs/
├── docs/
└── pkg/
```

---

## 📚 Más Detalles

Para más detalles técnicos, consulta los archivos de documentación relevantes en la carpeta `docs/`.
