# GitHub Issues para Brain CLI

Este archivo contiene las issues sugeridas para crear en GitHub, organizadas por hitos según el plan de desarrollo.

**Instrucciones**: Copia y pega cada issue en GitHub siguiendo el formato proporcionado.

---

## 📋 Hito 1: Fundación y Configuración

### Issue #1: Implementar sistema de configuración con YAML
**Labels**: `enhancement`, `hito-1`, `core`  
**Milestone**: Hito 1: Fundación

**Descripción**:
Implementar el sistema de carga y gestión de configuración desde archivos YAML.

**Tareas**:
- [ ] Crear interfaz `ConfigManager` en `internal/core/`
- [ ] Implementar `YAMLConfigLoader` en `internal/adapters/config/`
- [ ] Añadir soporte para variables de entorno (ej: `${OMNIROUTE_API_KEY}`)
- [ ] Implementar validación de configuración
- [ ] Añadir tests unitarios para config loader
- [ ] Documentar opciones de configuración en README

**Archivos involucrados**:
- `internal/adapters/config/config.go`
- `internal/adapters/config/loader.go`
- `configs/config.yaml`

**Criterios de aceptación**:
- La aplicación puede cargar `configs/config.yaml` exitosamente
- Variables de entorno se resuelven correctamente
- Errores de configuración se reportan claramente
- Tests cubren casos normales y edge cases

---

### Issue #2: Implementar logger estructurado con zerolog
**Labels**: `enhancement`, `hito-1`, `infrastructure`  
**Milestone**: Hito 1: Fundación

**Descripción**:
Configurar sistema de logging estructurado con diferentes niveles y output a archivo.

**Tareas**:
- [ ] Crear wrapper de zerolog en `pkg/utils/logger.go`
- [ ] Implementar niveles de log (debug, info, warn, error)
- [ ] Añadir rotación de logs
- [ ] Configurar formato de logs (JSON para producción, pretty para dev)
- [ ] Añadir context logging para tracing
- [ ] Documentar uso del logger

**Archivos involucrados**:
- `pkg/utils/logger.go`
- `configs/config.yaml` (sección de logging)

**Criterios de aceptación**:
- Logs se escriben en archivo y consola según configuración
- Diferentes niveles funcionan correctamente
- Logs incluyen timestamp, level, mensaje y contexto
- Logs rotan automáticamente al alcanzar tamaño máximo

---

### Issue #3: Crear entidades del dominio (Task, Execution)
**Labels**: `enhancement`, `hito-1`, `core`, `domain`  
**Milestone**: Hito 1: Fundación

**Descripción**:
Implementar las entidades principales del dominio según Clean Architecture.

**Tareas**:
- [ ] Crear entidad `Task` en `internal/core/task/entity.go`
- [ ] Crear enums `TaskType` y `Priority` en `internal/core/task/types.go`
- [ ] Crear entidad `Execution` en `internal/core/execution/entity.go`
- [ ] Crear enum `ExecutionStatus` en `internal/core/execution/types.go`
- [ ] Implementar métodos de validación en entidades
- [ ] Añadir tests unitarios para entidades
- [ ] Documentar estructura de entidades

**Archivos involucrados**:
- `internal/core/task/entity.go`
- `internal/core/task/types.go`
- `internal/core/execution/entity.go`
- `internal/core/execution/types.go`

**Criterios de aceptación**:
- Entidades tienen todos los campos requeridos
- Validaciones funcionan correctamente
- Tests cubren casos válidos e inválidos
- Documentación clara de cada campo

---

### Issue #4: Implementar "Hello World" con Bubble Tea
**Labels**: `enhancement`, `hito-1`, `ui`  
**Milestone**: Hito 1: Fundación

**Descripción**:
Crear una aplicación TUI básica con Bubble Tea como proof of concept.

**Tareas**:
- [ ] Crear estructura básica de modelo en `internal/ui/app.go`
- [ ] Implementar métodos Init, Update, View
- [ ] Añadir componente de bienvenida simple
- [ ] Implementar manejo de tecla 'q' para salir
- [ ] Añadir estilos básicos con lipgloss
- [ ] Integrar con main.go
- [ ] Documentar arquitectura del TUI

**Archivos involucrados**:
- `internal/ui/app.go`
- `cmd/tui-assistant/main.go`

**Criterios de aceptación**:
- La aplicación muestra una pantalla de bienvenida
- Se puede salir con 'q' o Ctrl+C
- No hay crashes ni panics
- El código sigue el patrón Elm de Bubble Tea

---

## 📋 Hito 2: Core y Tareas Básicas

### Issue #5: Implementar repositorios de Task y Execution
**Labels**: `enhancement`, `hito-2`, `core`, `repository`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Crear interfaces de repositorio y sus implementaciones.

**Tareas**:
- [ ] Crear interfaz `TaskRepository` en `internal/core/task/repository.go`
- [ ] Crear interfaz `ExecutionRepository` en `internal/core/execution/repository.go`
- [ ] Implementar `YAMLTaskRepository` en `internal/adapters/repository/task_yaml.go`
- [ ] Implementar `InMemoryExecutionRepository` en `internal/adapters/repository/execution_memory.go`
- [ ] Añadir tests para repositorios
- [ ] Documentar contratos de repositorios

**Archivos involucrados**:
- `internal/core/task/repository.go`
- `internal/core/execution/repository.go`
- `internal/adapters/repository/task_yaml.go`
- `internal/adapters/repository/execution_memory.go`

**Criterios de aceptación**:
- TaskRepository puede cargar tareas desde YAML
- ExecutionRepository puede guardar y recuperar ejecuciones
- Interfaces están bien definidas
- Tests cubren todas las operaciones CRUD

---

### Issue #6: Implementar CLI Executor para comandos y scripts
**Labels**: `enhancement`, `hito-2`, `infrastructure`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Crear adaptador para ejecutar comandos del sistema y scripts bash.

**Tareas**:
- [ ] Crear interfaz `CLIExecutor` en `internal/core/`
- [ ] Implementar `SystemCLIExecutor` en `internal/adapters/cli/executor.go`
- [ ] Añadir captura de stdout/stderr en tiempo real
- [ ] Implementar manejo de timeouts
- [ ] Implementar manejo de señales (Ctrl+C)
- [ ] Añadir soporte para variables de entorno
- [ ] Crear tests con comandos mock
- [ ] Documentar uso del executor

**Archivos involucrados**:
- `internal/adapters/cli/executor.go`
- `internal/adapters/cli/scripts.go`

**Criterios de aceptación**:
- Puede ejecutar comandos simples
- Puede ejecutar scripts bash
- Captura output correctamente
- Timeouts funcionan correctamente
- Ctrl+C cancela ejecución gracefully

---

### Issue #7: Implementar casos de uso de Task (Create, List, Execute)
**Labels**: `enhancement`, `hito-2`, `core`, `use-case`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Crear la lógica de aplicación para gestionar tareas.

**Tareas**:
- [ ] Implementar `CreateTask` use case en `internal/usecases/task/create.go`
- [ ] Implementar `ListTasks` use case en `internal/usecases/task/list.go`
- [ ] Implementar `ExecuteTask` use case en `internal/usecases/task/execute.go`
- [ ] Crear `TaskExecutor` orchestrator en `internal/usecases/task/executor.go`
- [ ] Añadir manejo de errores robusto
- [ ] Implementar tests unitarios
- [ ] Documentar flujos de ejecución

**Archivos involucrados**:
- `internal/usecases/task/create.go`
- `internal/usecases/task/list.go`
- `internal/usecases/task/execute.go`
- `internal/usecases/task/executor.go`

**Criterios de aceptación**:
- Se pueden listar tareas desde YAML
- Se pueden ejecutar tareas según su tipo
- Estados de ejecución se actualizan correctamente
- Errores se propagan adecuadamente

---

### Issue #8: Crear scripts de automatización básicos
**Labels**: `enhancement`, `hito-2`, `scripts`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Implementar scripts bash para las tareas básicas de automatización.

**Tareas**:
- [ ] Crear `scripts/wifi-vpn.sh` (conectar WiFi y VPN)
- [ ] Crear `scripts/start-dev.sh` (levantar stack de desarrollo)
- [ ] Crear `scripts/cleanup.sh` (limpiar temporales)
- [ ] Crear `scripts/backup.sh` (backup rápido)
- [ ] Crear `scripts/monitor.sh` (monitorear recursos)
- [ ] Añadir manejo de errores en scripts
- [ ] Hacer scripts ejecutables
- [ ] Documentar uso de cada script

**Archivos involucrados**:
- `scripts/wifi-vpn.sh`
- `scripts/start-dev.sh`
- `scripts/cleanup.sh`
- `scripts/backup.sh`
- `scripts/monitor.sh`

**Criterios de aceptación**:
- Todos los scripts son ejecutables
- Scripts reportan errores claramente
- Scripts tienen output informativo
- Scripts están documentados con comentarios

---

### Issue #9: Implementar UI de lista de tareas con Bubble Tea
**Labels**: `enhancement`, `hito-2`, `ui`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Crear pantalla principal que muestra lista de tareas disponibles.

**Tareas**:
- [ ] Crear `MainScreen` en `internal/ui/screens/main.go`
- [ ] Integrar `bubbles.List` para mostrar tareas
- [ ] Implementar navegación con flechas y Enter
- [ ] Añadir filtrado por tags (opcional)
- [ ] Implementar status bar con información útil
- [ ] Añadir estilos con lipgloss
- [ ] Conectar con use cases de Task
- [ ] Documentar componentes UI

**Archivos involucrados**:
- `internal/ui/screens/main.go`
- `internal/ui/screens/task_list.go`
- `internal/ui/components/status_bar.go`

**Criterios de aceptación**:
- Lista muestra tareas cargadas desde YAML
- Navegación funciona con teclado
- Se puede seleccionar una tarea
- UI es responsive y no tiene glitches

---

### Issue #10: Implementar pantalla de ejecución con feedback en tiempo real
**Labels**: `enhancement`, `hito-2`, `ui`  
**Milestone**: Hito 2: Core y Tareas Básicas

**Descripción**:
Crear pantalla que muestra el progreso de ejecución de una tarea.

**Tareas**:
- [ ] Crear `ExecutionScreen` en `internal/ui/screens/execution.go`
- [ ] Integrar `bubbles.Spinner` para mostrar progreso
- [ ] Mostrar output de comando en tiempo real
- [ ] Implementar indicador de estado (running, completed, failed)
- [ ] Añadir botón de cancelar
- [ ] Mostrar duración de ejecución
- [ ] Manejar transición de vuelta a main screen
- [ ] Documentar estados de la pantalla

**Archivos involucrados**:
- `internal/ui/screens/execution.go`
- `internal/ui/components/spinner.go`

**Criterios de aceptación**:
- Spinner muestra que la tarea está ejecutándose
- Output se muestra línea por línea
- Se puede cancelar con Ctrl+C
- Estado final se muestra claramente (éxito/error)
- Se regresa automáticamente a main screen

---

## 📋 Hito 3: Integración con OmniRoute

### Issue #11: Implementar cliente HTTP para OmniRoute
**Labels**: `enhancement`, `hito-3`, `infrastructure`, `ai`  
**Milestone**: Hito 3: Integración con OmniRoute

**Descripción**:
Crear cliente HTTP para comunicarse con la API de OmniRoute.

**Tareas**:
- [ ] Crear interfaz `OmniRouteClient` en `internal/core/omniroute/client.go`
- [ ] Implementar `HTTPOmniRouteClient` en `internal/adapters/omniroute/client.go`
- [ ] Añadir autenticación con Bearer token
- [ ] Implementar retry logic con exponential backoff
- [ ] Manejar rate limiting (429)
- [ ] Añadir timeouts configurables
- [ ] Implementar tests con server mock
- [ ] Documentar uso del cliente

**Archivos involucrados**:
- `internal/core/omniroute/client.go`
- `internal/core/omniroute/entity.go`
- `internal/adapters/omniroute/client.go`
- `internal/adapters/omniroute/models.go`
- `internal/adapters/omniroute/retry.go`

**Criterios de aceptación**:
- Cliente puede hacer peticiones exitosas a OmniRoute
- Autenticación funciona correctamente
- Retry logic maneja errores transitorios
- Timeouts previenen colgones
- Tests cubren casos normales y error cases

---

### Issue #12: Implementar casos de uso de OmniRoute (Chat, Classify, Summarize)
**Labels**: `enhancement`, `hito-3`, `core`, `use-case`, `ai`  
**Milestone**: Hito 3: Integración con OmniRoute

**Descripción**:
Crear lógica de aplicación para operaciones de IA.

**Tareas**:
- [ ] Implementar `Chat` use case en `internal/usecases/omniroute/chat.go`
- [ ] Implementar `Classify` use case en `internal/usecases/omniroute/classify.go`
- [ ] Implementar `Summarize` use case en `internal/usecases/omniroute/summarize.go`
- [ ] Añadir construcción de prompts desde templates
- [ ] Implementar manejo de streaming (opcional)
- [ ] Añadir tests con cliente mock
- [ ] Documentar casos de uso

**Archivos involucrados**:
- `internal/usecases/omniroute/chat.go`
- `internal/usecases/omniroute/classify.go`
- `internal/usecases/omniroute/summarize.go`

**Criterios de aceptación**:
- Se pueden enviar mensajes a OmniRoute
- Respuestas se procesan correctamente
- Prompts se construyen desde templates en YAML
- Tests validan lógica sin hacer llamadas reales

---

### Issue #13: Implementar tarea "Revisar Correos" con Gmail API
**Labels**: `enhancement`, `hito-3`, `feature`, `ai`  
**Milestone**: Hito 3: Integración con OmniRoute

**Descripción**:
Integrar con Gmail para leer y clasificar correos con IA.

**Tareas**:
- [ ] Implementar `GmailClient` en `internal/adapters/email/client.go`
- [ ] Configurar OAuth2 flow en `internal/adapters/email/auth.go`
- [ ] Implementar lectura de correos no leídos
- [ ] Integrar con OmniRoute para clasificación
- [ ] Crear tarea "review-emails" en `tasks.yaml`
- [ ] Mostrar resultados en UI con formato
- [ ] Documentar configuración de Gmail API

**Archivos involucrados**:
- `internal/adapters/email/client.go`
- `internal/adapters/email/auth.go`
- `configs/tasks.yaml`

**Criterios de aceptación**:
- OAuth2 flow funciona correctamente
- Se pueden listar correos no leídos
- OmniRoute clasifica correos por urgencia
- Resultados se muestran de forma legible

---

### Issue #14: Implementar tarea "Resumir URL"
**Labels**: `enhancement`, `hito-3`, `feature`, `ai`  
**Milestone**: Hito 3: Integración con OmniRoute

**Descripción**:
Crear tarea que extrae contenido de URL y lo resume con IA.

**Tareas**:
- [ ] Implementar web scraper para extraer contenido
- [ ] Limpiar HTML y extraer texto relevante
- [ ] Integrar con OmniRoute para resumen
- [ ] Crear tarea "summarize-url" en `tasks.yaml`
- [ ] Manejar errores de red y URLs inválidas
- [ ] Mostrar resumen en UI con Markdown
- [ ] Documentar uso de la tarea

**Archivos involucrados**:
- `internal/adapters/web/scraper.go`
- `configs/tasks.yaml`

**Criterios de aceptación**:
- Se puede extraer contenido de URLs válidas
- OmniRoute genera resumen conciso
- Errores se manejan gracefully
- Resumen se renderiza con glamour

---

### Issue #15: Implementar renderizado de Markdown en UI
**Labels**: `enhancement`, `hito-3`, `ui`  
**Milestone**: Hito 3: Integración con OmniRoute

**Descripción**:
Añadir capacidad de mostrar respuestas de IA en formato Markdown.

**Tareas**:
- [ ] Integrar glamour para renderizado de Markdown
- [ ] Crear componente de vista de Markdown
- [ ] Añadir scroll para contenido largo
- [ ] Implementar tema que combine con TUI
- [ ] Añadir acciones (copiar, guardar)
- [ ] Documentar componente

**Archivos involucrados**:
- `internal/ui/components/markdown.go`
- `internal/ui/screens/ai_response.go`

**Criterios de aceptación**:
- Markdown se renderiza correctamente en terminal
- Scroll funciona para contenido largo
- Tema es consistente con resto de UI
- Se puede copiar contenido al clipboard

---

## 📋 Hito 4: UI Avanzada y Experiencia

### Issue #16: Implementar navegación multi-pantalla
**Labels**: `enhancement`, `hito-4`, `ui`  
**Milestone**: Hito 4: UI Avanzada

**Descripción**:
Crear sistema de navegación entre múltiples pantallas.

**Tareas**:
- [ ] Implementar stack de pantallas en `internal/ui/app.go`
- [ ] Añadir transiciones entre pantallas
- [ ] Implementar navegación back con Esc
- [ ] Crear pantalla de Settings
- [ ] Crear pantalla de History
- [ ] Añadir animaciones de transición (opcional)
- [ ] Documentar arquitectura de navegación

**Archivos involucrados**:
- `internal/ui/app.go`
- `internal/ui/screens/settings.go`
- `internal/ui/screens/history.go`

**Criterios de aceptación**:
- Se puede navegar entre pantallas sin bugs
- Back button funciona correctamente
- Estado se preserva entre navegaciones
- No hay memory leaks

---

### Issue #17: Implementar persistencia de historial con SQLite
**Labels**: `enhancement`, `hito-4`, `infrastructure`  
**Milestone**: Hito 4: UI Avanzada

**Descripción**:
Guardar historial de ejecuciones en base de datos SQLite.

**Tareas**:
- [ ] Crear esquema de base de datos
- [ ] Implementar migrations en `internal/adapters/repository/migrations.go`
- [ ] Implementar `SQLiteExecutionRepository`
- [ ] Añadir índices para búsquedas rápidas
- [ ] Implementar queries de historial
- [ ] Añadir límite de registros (auto-cleanup)
- [ ] Documentar esquema y queries

**Archivos involucrados**:
- `internal/adapters/repository/execution_sqlite.go`
- `internal/adapters/repository/migrations.go`

**Criterios de aceptación**:
- Ejecuciones se guardan correctamente
- Se pueden consultar por fecha, estado, tarea
- Auto-cleanup funciona según límite configurado
- No hay race conditions con escrituras concurrentes

---

### Issue #18: Implementar pantalla de historial consultable
**Labels**: `enhancement`, `hito-4`, `ui`  
**Milestone**: Hito 4: UI Avanzada

**Descripción**:
Crear UI para consultar historial de ejecuciones.

**Tareas**:
- [ ] Crear `HistoryScreen` en `internal/ui/screens/history.go`
- [ ] Mostrar tabla con ejecuciones recientes
- [ ] Implementar filtrado por estado, fecha, tarea
- [ ] Añadir vista de detalles de ejecución
- [ ] Implementar acción de re-ejecutar tarea
- [ ] Añadir exportación a JSON/CSV
- [ ] Documentar pantalla

**Archivos involucrados**:
- `internal/ui/screens/history.go`
- `internal/ui/components/table.go`

**Criterios de aceptación**:
- Historial se carga correctamente desde DB
- Filtros funcionan sin lag
- Se pueden ver detalles de cualquier ejecución
- Re-ejecutar funciona correctamente
- Exportación genera archivo válido

---

### Issue #19: Implementar pantalla de configuración editable
**Labels**: `enhancement`, `hito-4`, `ui`  
**Milestone**: Hito 4: UI Avanzada

**Descripción**:
Crear UI para editar configuración sin tocar archivos.

**Tareas**:
- [ ] Crear `SettingsScreen` en `internal/ui/screens/settings.go`
- [ ] Mostrar opciones de configuración actual
- [ ] Implementar forms para editar valores
- [ ] Validar inputs antes de guardar
- [ ] Persistir cambios a `config.yaml`
- [ ] Añadir opción de reset a defaults
- [ ] Documentar opciones configurables

**Archivos involucrados**:
- `internal/ui/screens/settings.go`
- `internal/adapters/config/writer.go`

**Criterios de aceptación**:
- Se muestran todas las opciones editables
- Cambios se validan correctamente
- Configuración se guarda a archivo YAML
- Cambios se aplican inmediatamente (o con restart)

---

### Issue #20: Implementar sistema de atajos de teclado configurables
**Labels**: `enhancement`, `hito-4`, `ui`  
**Milestone**: Hito 4: UI Avanzada

**Descripción**:
Añadir soporte para hotkeys configurables por el usuario.

**Tareas**:
- [ ] Crear sistema de keybindings en `internal/ui/keybindings.go`
- [ ] Cargar keybindings desde configuración
- [ ] Implementar pantalla de help con lista de atajos
- [ ] Añadir conflictos detection
- [ ] Permitir edición desde Settings
- [ ] Documentar atajos por defecto

**Archivos involucrados**:
- `internal/ui/keybindings.go`
- `internal/ui/screens/help.go`
- `configs/config.yaml`

**Criterios de aceptación**:
- Atajos se pueden configurar desde YAML
- Help screen muestra atajos actuales
- Conflictos se detectan y reportan
- Atajos funcionan en todas las pantallas relevantes

---

## 📋 Hito 5: Pulido y Lanzamiento

### Issue #21: Implementar suite completa de tests
**Labels**: `testing`, `hito-5`  
**Milestone**: Hito 5: Pulido y Lanzamiento

**Descripción**:
Asegurar cobertura de tests >70% en todo el código.

**Tareas**:
- [ ] Unit tests para todas las entidades
- [ ] Unit tests para todos los use cases
- [ ] Integration tests para adapters
- [ ] E2E tests para flujos críticos
- [ ] Benchmark tests para operaciones críticas
- [ ] Configurar CI para ejecutar tests
- [ ] Generar reportes de cobertura

**Criterios de aceptación**:
- Cobertura >70% en toda la base de código
- Todos los tests pasan
- CI ejecuta tests automáticamente
- Reportes de cobertura se generan

---

### Issue #22: Completar documentación del proyecto
**Labels**: `documentation`, `hito-5`  
**Milestone**: Hito 5: Pulido y Lanzamiento

**Descripción**:
Documentar completamente el proyecto para usuarios y contributors.

**Tareas**:
- [ ] Actualizar README con screenshots
- [ ] Crear guía de usuario
- [ ] Crear guía de contributor
- [ ] Documentar API interna
- [ ] Crear ejemplos de uso
- [ ] Grabar video demo (opcional)
- [ ] Crear FAQ

**Archivos involucrados**:
- `README.md`
- `docs/user-guide.md`
- `docs/contributing.md`
- `docs/api.md`

**Criterios de aceptación**:
- README tiene screenshots y ejemplos
- Guías están completas y claras
- Todos los exports públicos están documentados
- Ejemplos funcionan correctamente

---

### Issue #23: Crear script de instalación y releases
**Labels**: `infrastructure`, `hito-5`  
**Milestone**: Hito 5: Pulido y Lanzamiento

**Descripción**:
Preparar proceso de instalación y releases automáticas.

**Tareas**:
- [ ] Crear script de instalación bash
- [ ] Configurar GitHub Actions para releases
- [ ] Compilar binarios para múltiples plataformas
- [ ] Crear checksums para verificación
- [ ] Configurar versionado semántico
- [ ] Documentar proceso de release
- [ ] Publicar primer release

**Archivos involucrados**:
- `install.sh`
- `.github/workflows/release.yml`
- `Makefile`

**Criterios de aceptación**:
- Script de instalación funciona en Linux/macOS
- Releases se crean automáticamente con tags
- Binarios están disponibles para múltiples plataformas
- Checksums se generan correctamente

---

### Issue #24: Implementar tareas adicionales planificadas
**Labels**: `enhancement`, `hito-5`, `feature`  
**Milestone**: Hito 5: Pulido y Lanzamiento

**Descripción**:
Completar el conjunto de tareas planificadas en el roadmap.

**Tareas**:
- [ ] Implementar "Git Quick Actions"
- [ ] Implementar "Notes Manager"
- [ ] Implementar "Kill Process by Port"
- [ ] Crear scripts correspondientes
- [ ] Añadir tareas a `tasks.yaml`
- [ ] Documentar cada tarea
- [ ] Añadir tests

**Archivos involucrados**:
- `configs/tasks.yaml`
- `scripts/git-quick.sh`
- `scripts/take-note.sh`
- `scripts/kill-port.sh`

**Criterios de aceptación**:
- Todas las tareas funcionan correctamente
- Scripts están documentados
- Tareas aparecen en UI
- Tests cubren funcionalidad básica

---

## 📌 Issues Transversales

### Issue #25: Configurar CI/CD con GitHub Actions
**Labels**: `infrastructure`, `ci-cd`

**Descripción**:
Configurar pipeline de CI/CD para automatizar tests y releases.

**Tareas**:
- [ ] Crear workflow para tests en push
- [ ] Crear workflow para releases con tags
- [ ] Configurar cache de dependencias
- [ ] Añadir linting en CI
- [ ] Configurar matrix para múltiples versiones de Go
- [ ] Documentar pipelines

**Archivos involucrados**:
- `.github/workflows/test.yml`
- `.github/workflows/release.yml`

---

### Issue #26: Optimización de performance
**Labels**: `enhancement`, `performance`

**Descripción**:
Identificar y optimizar bottlenecks de performance.

**Tareas**:
- [ ] Perfilar aplicación con pprof
- [ ] Optimizar queries de SQLite
- [ ] Reducir allocations en hot paths
- [ ] Optimizar renderizado de UI
- [ ] Añadir benchmarks
- [ ] Documentar resultados

---

### Issue #27: Seguridad y manejo de secrets
**Labels**: `security`, `high-priority`

**Descripción**:
Asegurar que secrets se manejan correctamente.

**Tareas**:
- [ ] Validar que no se loggean secrets
- [ ] Implementar encriptación de tokens (opcional)
- [ ] Añadir validación de permisos de archivos
- [ ] Documentar mejores prácticas
- [ ] Añadir security checklist

---

## 📝 Notas

- Cada issue debe tener labels apropiados para filtrado
- Asignar issues a milestones correspondientes
- Priorizar issues de fundación antes de avanzar
- Mantener issues actualizadas con progreso
- Cerrar issues solo cuando todos los criterios de aceptación se cumplan

---

**Total de issues**: 27  
**Hito 1**: 4 issues  
**Hito 2**: 6 issues  
**Hito 3**: 5 issues  
**Hito 4**: 5 issues  
**Hito 5**: 4 issues  
**Transversales**: 3 issues
