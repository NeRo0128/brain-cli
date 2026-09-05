# Plan de Desarrollo: TUI Assistant (Brain CLI)

**Fecha de creación:** 2026-09-05  
**Versión:** 1.0  
**Arquitectura:** Clean Architecture  
**Stack principal:** Go + Bubble Tea + OmniRoute AI

---

## 📋 Tabla de Contenidos

1. [Visión General](#visión-general)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [Hitos de Desarrollo](#hitos-de-desarrollo)
4. [Tareas Detalladas por Componente](#tareas-detalladas-por-componente)
5. [Dependencias y Tecnologías](#dependencias-y-tecnologías)
6. [Criterios de Aceptación](#criterios-de-aceptación)
7. [Métricas de Éxito](#métricas-de-éxito)

---

## 🎯 Visión General

Brain CLI es un asistente personal de línea de comandos (TUI) diseñado para automatizar tareas repetitivas del día a día, integrar capacidades de IA mediante OmniRoute, y proporcionar una interfaz elegante y eficiente para gestionar workflows complejos.

### Objetivos Principales

1. **Automatización Inteligente**: Ejecutar scripts y comandos del sistema de forma organizada
2. **Integración con IA**: Aprovechar OmniRoute para clasificación, resúmenes y procesamiento de lenguaje natural
3. **Experiencia de Usuario**: TUI moderna, responsive y con feedback visual inmediato
4. **Extensibilidad**: Arquitectura limpia que permita añadir nuevas tareas fácilmente
5. **Productividad**: Reducir tiempo en tareas manuales repetitivas

### Alcance Inicial (MVP)

- ✅ Gestión de tareas predefinidas (WiFi/VPN, stack dev, limpieza)
- ✅ Ejecución de scripts del sistema
- ✅ Integración básica con OmniRoute (chat, clasificación)
- ✅ TUI con navegación y feedback visual
- ✅ Configuración mediante archivos YAML
- ✅ Logging y historial de ejecuciones

---

## 🏗️ Arquitectura del Sistema

### Capas de Clean Architecture

```
┌─────────────────────────────────────────┐
│         UI Layer (Bubble Tea)           │  ← Presentación
├─────────────────────────────────────────┤
│       Use Cases (Application Logic)     │  ← Lógica de aplicación
├─────────────────────────────────────────┤
│      Domain (Entities & Interfaces)     │  ← Núcleo de negocio
├─────────────────────────────────────────┤
│    Adapters (Infrastructure/External)   │  ← Implementaciones concretas
└─────────────────────────────────────────┘
```

### Flujo de Dependencias

```
UI → UseCases → Domain ← Adapters
```

**Principio:** Las capas internas NO conocen las externas. La dirección de las dependencias siempre apunta hacia el dominio.

### Componentes Principales

#### 1. **Domain Layer** (`internal/core/`)
- **Task Entity**: Representa una tarea ejecutable
- **Execution Entity**: Representa una ejecución con estado
- **OmniRoute Entity**: Modelos de petición/respuesta de IA
- **Repository Interfaces**: Contratos para persistencia

#### 2. **Use Cases Layer** (`internal/usecases/`)
- **Task Management**: Crear, listar, ejecutar tareas
- **Execution Management**: Gestionar estado de ejecuciones
- **OmniRoute Operations**: Chat, clasificación, resumen

#### 3. **Adapters Layer** (`internal/adapters/`)
- **CLI Executor**: Ejecuta scripts y comandos del sistema
- **OmniRoute Client**: Cliente HTTP para la API
- **Email Client**: Integración con Gmail (futuro)
- **Config Manager**: Lee/escribe configuración YAML

#### 4. **UI Layer** (`internal/ui/`)
- **Main Screen**: Vista principal con menú de tareas
- **Task List**: Lista interactiva de tareas
- **Execution View**: Estado en tiempo real de ejecución
- **Settings**: Configuración de la aplicación

---

## 🗓️ Hitos de Desarrollo

### **Hito 1: Fundación y Configuración** (Semana 1-2)

**Objetivo**: Tener un proyecto Go funcional con la estructura completa y un "Hello World" en Bubble Tea.

#### Tareas

- [x] Crear estructura de carpetas según Clean Architecture
- [ ] Inicializar `go.mod` con dependencias
- [ ] Implementar sistema de configuración (YAML)
- [ ] Configurar logger estructurado (zerolog)
- [ ] Crear esqueleto de entidades del dominio
- [ ] Implementar primer TUI "Hello World" con Bubble Tea
- [ ] Documentar arquitectura y decisiones de diseño

#### Entregables

- ✓ Estructura de carpetas completa
- ✓ `go.mod` y `go.sum` configurados
- ✓ Archivos `config.yaml` y `tasks.yaml` de ejemplo
- ✓ Aplicación que compila y muestra "Hello World"
- ✓ README con instrucciones de setup

#### Criterios de Aceptación

- La aplicación compila sin errores
- Se puede ejecutar `go run cmd/tui-assistant/main.go`
- Se muestra una interfaz TUI básica
- La configuración se carga desde YAML

---

### **Hito 2: Core Domain y Tareas Básicas** (Semana 3-4)

**Objetivo**: Implementar la lógica de negocio central y las primeras tareas ejecutables.

#### Tareas del Dominio

- [ ] Definir entidad `Task` con todos sus campos
- [ ] Definir entidad `Execution` con estados (pending, running, completed, failed)
- [ ] Crear interfaces de repositorio (`TaskRepository`, `ExecutionRepository`)
- [ ] Implementar validaciones de negocio en entidades

#### Tareas de Use Cases

- [ ] UC: CreateTask - Crear nueva tarea desde configuración
- [ ] UC: ListTasks - Listar tareas disponibles
- [ ] UC: ExecuteTask - Ejecutar una tarea y gestionar su estado
- [ ] UC: GetExecutionStatus - Obtener estado de ejecución en tiempo real

#### Tareas de Adapters

- [ ] Implementar `CLIExecutor` para ejecutar comandos del sistema
- [ ] Implementar `ScriptExecutor` para ejecutar scripts en `scripts/`
- [ ] Implementar `YAMLTaskRepository` para leer `tasks.yaml`
- [ ] Implementar `InMemoryExecutionRepository` (temporal)

#### Tareas de Scripts

- [ ] Crear `scripts/wifi-vpn.sh` (conectar WiFi + VPN)
- [ ] Crear `scripts/start-dev.sh` (levantar OmniRoute + OpenCode)
- [ ] Crear `scripts/cleanup.sh` (limpiar temporales)
- [ ] Crear `scripts/backup.sh` (backup rápido)

#### Tareas de UI

- [ ] Implementar pantalla principal con menú de tareas
- [ ] Implementar navegación con teclado
- [ ] Mostrar lista de tareas disponibles
- [ ] Mostrar feedback visual al ejecutar tarea

#### Entregables

- ✓ Modelo de dominio completo y validado
- ✓ 4 scripts de automatización funcionales
- ✓ TUI que lista y ejecuta tareas básicas
- ✓ Logs de ejecución visibles en terminal

#### Criterios de Aceptación

- Se pueden listar tareas desde `tasks.yaml`
- Se puede ejecutar una tarea y ver su progreso
- Los scripts se ejecutan correctamente
- La UI responde a teclas de navegación
- Los errores se manejan y se muestran al usuario

---

### **Hito 3: Integración con OmniRoute** (Semana 5-6)

**Objetivo**: Conectar con la API de OmniRoute y añadir tareas con capacidades de IA.

#### Tareas del Dominio

- [ ] Definir entidad `OmniRouteRequest` (prompt, model, params)
- [ ] Definir entidad `OmniRouteResponse` (content, tokens, metadata)
- [ ] Crear interfaz `OmniRouteClient`

#### Tareas de Use Cases

- [ ] UC: ChatWithAI - Enviar mensaje y recibir respuesta
- [ ] UC: ClassifyContent - Clasificar texto (correos, documentos)
- [ ] UC: SummarizeContent - Resumir texto o URLs

#### Tareas de Adapters

- [ ] Implementar `HTTPOmniRouteClient` con autenticación
- [ ] Manejar streaming de respuestas (si aplica)
- [ ] Implementar retry logic y timeouts
- [ ] Implementar caché de respuestas (opcional)

#### Tareas de Integración

- [ ] **Tarea: "Revisar Correos"** - Conectar con Gmail API y clasificar con IA
- [ ] **Tarea: "Resumir URL"** - Fetch contenido web y resumir con IA
- [ ] **Tarea: "Traductor Natural"** - Convertir lenguaje natural a comando

#### Tareas de UI

- [ ] Mostrar respuestas de IA en formato legible (Markdown)
- [ ] Añadir spinner animado durante peticiones a IA
- [ ] Mostrar progress bar para operaciones largas
- [ ] Implementar modal de confirmación para tareas críticas

#### Entregables

- ✓ Cliente HTTP funcional para OmniRoute
- ✓ 3 tareas con IA funcionando end-to-end
- ✓ UI con feedback visual para operaciones asíncronas
- ✓ Manejo robusto de errores de red

#### Criterios de Aceptación

- Las peticiones a OmniRoute funcionan correctamente
- Las respuestas de IA se muestran en la TUI
- Los errores de red se manejan gracefully
- Las operaciones largas muestran progreso visual
- Se puede interrumpir una operación en curso

---

### **Hito 4: UI Avanzada y Experiencia** (Semana 7-8)

**Objetivo**: Pulir la experiencia de usuario y añadir funcionalidades avanzadas de navegación.

#### Tareas de UI

- [ ] Implementar navegación multi-pantalla (stack de vistas)
- [ ] Crear pantalla de configuración editable
- [ ] Crear pantalla de historial de ejecuciones
- [ ] Implementar status bar con información contextual
- [ ] Añadir sistema de notificaciones/toasts
- [ ] Implementar búsqueda/filtrado de tareas
- [ ] Añadir atajos de teclado configurables
- [ ] Implementar temas de color (dark/light)

#### Tareas de Persistencia

- [ ] Implementar `SQLiteExecutionRepository` para historial
- [ ] Guardar logs de ejecución en archivos
- [ ] Implementar sistema de favoritos para tareas
- [ ] Exportar historial a JSON/CSV

#### Tareas de UX

- [ ] Añadir animaciones suaves entre transiciones
- [ ] Implementar help screen con lista de atajos
- [ ] Añadir confirmaciones para tareas destructivas
- [ ] Implementar modo "dry-run" para previsualizar acciones

#### Entregables

- ✓ TUI completa con múltiples pantallas navegables
- ✓ Sistema de persistencia funcional
- ✓ Historial de ejecuciones consultable
- ✓ Configuración editable desde la TUI

#### Criterios de Aceptación

- Se puede navegar entre pantallas sin bugs
- El historial se guarda y se puede consultar
- Los atajos de teclado funcionan correctamente
- La configuración se persiste al editarla
- Los temas de color cambian correctamente

---

### **Hito 5: Tareas Adicionales y Pulido** (Semana 9-10)

**Objetivo**: Completar el conjunto de tareas planeadas y preparar para lanzamiento.

#### Tareas Nuevas a Implementar

- [ ] **Monitor de Recursos** - Mostrar CPU, RAM, Disco en tiempo real
- [ ] **Gestor de Procesos** - Listar y matar procesos
- [ ] **Abrir URLs Frecuentes** - Quick launcher de sitios
- [ ] **Clasificador de Archivos** - Organizar archivos con IA
- [ ] **Git Quick Actions** - Status, pull, push desde TUI
- [ ] **Notes Manager** - Tomar notas rápidas

#### Tareas de Testing

- [ ] Unit tests para domain entities (>80% coverage)
- [ ] Unit tests para use cases (>80% coverage)
- [ ] Integration tests para adapters
- [ ] E2E tests para flujos críticos
- [ ] Benchmark tests para operaciones críticas

#### Tareas de Documentación

- [ ] Completar README con screenshots
- [ ] Documentar cada tarea en `tasks.yaml`
- [ ] Crear guía de contribución
- [ ] Documentar API de OmniRoute
- [ ] Crear video demo (opcional)

#### Tareas de Deployment

- [ ] Crear Makefile con comandos comunes
- [ ] Script de instalación automática
- [ ] Compilar binarios para múltiples plataformas
- [ ] Crear releases en GitHub
- [ ] Publicar en Homebrew (opcional)

#### Entregables

- ✓ Suite completa de tareas implementadas
- ✓ Tests con buena cobertura
- ✓ Documentación completa
- ✓ Binarios listos para distribución

#### Criterios de Aceptación

- Todas las tareas planeadas funcionan correctamente
- Tests pasan exitosamente
- Documentación está completa y clara
- Se puede instalar fácilmente en un sistema limpio
- No hay bugs críticos conocidos

---

## 📦 Tareas Detalladas por Componente

### **1. Core Domain (`internal/core/`)**

#### `task/entity.go`
```go
type Task struct {
    ID          string
    Name        string
    Description string
    Type        TaskType // ScriptTask, CommandTask, AITask
    Command     string
    ScriptPath  string
    RequiresAI  bool
    AIPrompt    string
    Tags        []string
    Priority    Priority
    CreatedAt   time.Time
}
```

#### `execution/entity.go`
```go
type Execution struct {
    ID         string
    TaskID     string
    Status     ExecutionStatus // Pending, Running, Completed, Failed
    Output     string
    Error      error
    StartedAt  time.Time
    FinishedAt time.Time
    Duration   time.Duration
}
```

#### `omniroute/entity.go`
```go
type OmniRouteRequest struct {
    Prompt      string
    Model       string
    Temperature float64
    MaxTokens   int
}

type OmniRouteResponse struct {
    Content   string
    Tokens    int
    Model     string
    CreatedAt time.Time
}
```

---

### **2. Use Cases (`internal/usecases/`)**

#### `task/execute.go`
```
Responsabilidad: Ejecutar una tarea y gestionar su ciclo de vida

Flujo:
1. Validar que la tarea existe
2. Crear una Execution con estado Pending
3. Cambiar estado a Running
4. Ejecutar la tarea según su tipo (Script, Command, AI)
5. Capturar output/error
6. Cambiar estado a Completed o Failed
7. Guardar Execution en repositorio
8. Retornar resultado
```

#### `omniroute/chat.go`
```
Responsabilidad: Interactuar con OmniRoute para chat

Flujo:
1. Validar input del usuario
2. Construir OmniRouteRequest
3. Enviar request al cliente HTTP
4. Manejar streaming si aplica
5. Retornar OmniRouteResponse
6. Loggear para análisis
```

---

### **3. Adapters (`internal/adapters/`)**

#### `cli/executor.go`
```
Responsabilidad: Ejecutar comandos del sistema

Implementación:
- Usar os/exec para ejecutar comandos
- Capturar stdout/stderr
- Manejar timeouts
- Manejar señales de interrupción (Ctrl+C)
```

#### `omniroute/client.go`
```
Responsabilidad: Cliente HTTP para OmniRoute API

Implementación:
- HTTP client con retry logic
- Autenticación (Bearer token)
- Parsing de JSON responses
- Manejo de rate limiting
- Timeouts configurables
```

#### `email/client.go`
```
Responsabilidad: Cliente para Gmail API

Implementación:
- OAuth2 flow para autenticación
- Listar correos no leídos
- Obtener contenido de correos
- Marcar como leído
```

---

### **4. UI (`internal/ui/`)**

#### `screens/main.go`
```
Responsabilidad: Pantalla principal con menú de tareas

Componentes:
- Lista de tareas (usando bubbles.List)
- Status bar con hora actual
- Ayuda de atajos de teclado
- Indicador de estado de conexión
```

#### `screens/execution.go`
```
Responsabilidad: Mostrar ejecución en tiempo real

Componentes:
- Spinner animado
- Progress bar
- Output en tiempo real (streaming)
- Botón de cancelar
```

#### `components/status_bar.go`
```
Responsabilidad: Barra de estado inferior

Información:
- Hora actual
- Estado de conexión OmniRoute
- Teclas de ayuda rápidas
- Indicador de notificaciones
```

---

## 🔧 Dependencias y Tecnologías

### Dependencias Go

#### UI y Terminal
```
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
github.com/charmbracelet/glamour
```

#### Configuración y Parsing
```
gopkg.in/yaml.v3
github.com/spf13/viper
```

#### Logging
```
github.com/rs/zerolog
```

#### HTTP y APIs
```
net/http (stdlib)
github.com/google/go-querystring
```

#### Persistencia
```
github.com/mattn/go-sqlite3
database/sql (stdlib)
```

#### Testing
```
github.com/stretchr/testify
```

#### Utilidades
```
github.com/google/uuid
golang.org/x/oauth2
google.golang.org/api/gmail/v1
```

### Herramientas Externas

- **OmniRoute**: API de IA para chat, clasificación, resumen
- **Gmail API**: Para revisar correos
- **Scripts Bash**: Automatizaciones del sistema

---

## ✅ Criterios de Aceptación Generales

### Funcionales

1. ✓ La aplicación compila sin warnings
2. ✓ Todas las tareas definidas se pueden ejecutar
3. ✓ La integración con OmniRoute funciona correctamente
4. ✓ La configuración se puede modificar sin recompilar
5. ✓ El historial de ejecuciones se persiste correctamente

### No Funcionales

1. ✓ La UI responde en <100ms a acciones del usuario
2. ✓ Las peticiones a OmniRoute tienen timeout de 30s
3. ✓ Los logs se rotan automáticamente (max 100MB)
4. ✓ El binario compilado pesa <20MB
5. ✓ La aplicación consume <50MB de RAM en idle

### Calidad de Código

1. ✓ Cobertura de tests >70%
2. ✓ `go vet` no reporta errores
3. ✓ `golangci-lint` pasa con configuración estándar
4. ✓ Todos los exports tienen documentación
5. ✓ No hay dependencias circulares

---

## 📊 Métricas de Éxito

### KPIs de Desarrollo

- **Velocity**: Completar 1 hito cada 2 semanas
- **Quality**: Mantener >70% code coverage
- **Stability**: <5 bugs críticos en producción

### KPIs de Usuario

- **Time to First Task**: <30 segundos desde instalación
- **Task Success Rate**: >95% de ejecuciones exitosas
- **Response Time**: <2s para operaciones locales, <5s para IA

### KPIs de Producto

- **Extensibilidad**: Añadir nueva tarea en <30 minutos
- **Mantenibilidad**: Onboarding de nuevo dev en <2 horas
- **Portabilidad**: Funciona en Linux, macOS, Windows

---

## 🔄 Proceso de Desarrollo

### Workflow Git

1. `main` branch: Código estable y funcional
2. `develop` branch: Integración de features
3. Feature branches: `feature/nombre-feature`
4. Hotfix branches: `hotfix/descripcion`

### Commits

Seguir Conventional Commits:
```
feat: añadir tarea de revisar correos
fix: corregir crash al cancelar ejecución
docs: actualizar README con screenshots
refactor: mejorar estructura de adapters
test: añadir tests para task executor
```

### Pull Requests

- Título descriptivo
- Descripción con contexto y screenshots (si aplica)
- Tests pasando
- Revisión de al menos 1 persona

---

## 📚 Referencias

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## 📝 Notas de Implementación

### Decisiones de Diseño

1. **¿Por qué Clean Architecture?**
   - Facilita testing independiente de infraestructura
   - Permite cambiar UI, DB o APIs sin afectar lógica de negocio
   - Escalable para proyectos grandes

2. **¿Por qué Bubble Tea?**
   - Modelo declarativo fácil de razonar
   - Excelente ecosistema (lipgloss, bubbles, glamour)
   - Activamente mantenido
   - Excelente documentación

3. **¿Por qué YAML para config?**
   - Fácil de editar manualmente
   - Soporta comentarios
   - Ampliamente usado en el ecosistema

4. **¿Por qué SQLite para historial?**
   - No requiere servidor externo
   - Perfecto para datos locales
   - Excelente performance para consultas
   - Soporta búsqueda y filtrado complejo

### Patrones a Aplicar

- **Repository Pattern**: Abstracción de persistencia
- **Strategy Pattern**: Diferentes tipos de ejecutores de tareas
- **Observer Pattern**: Notificación de cambios de estado
- **Command Pattern**: Encapsulación de tareas ejecutables
- **Factory Pattern**: Creación de tareas según tipo

### Anti-Patrones a Evitar

- ❌ Lógica de negocio en la UI
- ❌ Dependencias entre adapters
- ❌ God objects con demasiada responsabilidad
- ❌ Uso directo de librerías externas en el dominio
- ❌ Estados mutables compartidos sin sincronización

---

**Última actualización:** 2026-09-05  
**Estado del plan:** 🟢 Activo  
**Próxima revisión:** Fin de Hito 1 (Semana 2)
