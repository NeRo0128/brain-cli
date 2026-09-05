# 🧠 Brain CLI

> Asistente personal de línea de comandos con capacidades de IA

Brain CLI es un asistente TUI (Terminal User Interface) moderno construido en Go que automatiza tus tareas diarias del sistema y proporciona capacidades de inteligencia artificial mediante integración con OmniRoute.

## ✨ Características

- ⚡ **Automatización Inteligente**: Ejecuta scripts y comandos del sistema organizadamente
- 🤖 **Integración con IA**: Clasificación de correos, resúmenes, chat con OmniRoute
- 🎨 **Interfaz Elegante**: TUI moderna con Bubble Tea
- 📊 **Historial Completo**: Guarda todas las ejecuciones en SQLite
- 🔧 **Extensible**: Añade nuevas tareas fácilmente desde YAML
- ⚙️ **Configurable**: Personaliza todo desde archivos de configuración

## 🚀 Quick Start

### Prerrequisitos

- Go 1.23 o superior
- OmniRoute API (opcional, para funciones de IA)
- Bash (para ejecutar scripts)

### Instalación

```bash
# Clonar el repositorio
git clone https://github.com/nero/brain-cli.git
cd brain-cli

# Instalar dependencias
go mod download

# Configurar API key de OmniRoute (opcional)
export OMNIROUTE_API_KEY="tu-api-key-aqui"

# Compilar
go build -o brain-cli cmd/tui-assistant/main.go

# Ejecutar
./brain-cli
```

### Desarrollo

```bash
# Ejecutar en modo desarrollo
go run cmd/tui-assistant/main.go

# Ejecutar con debug logging
DEBUG=1 go run cmd/tui-assistant/main.go

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...
```

## 📋 Tareas Disponibles

### Automatización del Sistema

- **Conectar WiFi + VPN**: Conecta automáticamente a red WiFi y levanta VPN
- **Levantar Stack Dev**: Inicia OmniRoute y OpenCode
- **Limpiar Temporales**: Elimina archivos temporales y libera espacio
- **Backup Rápido**: Copia archivos importantes a backup
- **Monitor de Recursos**: Muestra uso de CPU, RAM y Disco

### Tareas con IA

- **Revisar Correos**: Clasifica correos por urgencia e importancia
- **Resumir URL**: Extrae y resume contenido de URLs
- **Traductor Natural**: Convierte lenguaje natural a comandos de terminal
- **Clasificador de Archivos**: Sugiere organización lógica de archivos
- **Chat con IA**: Conversación libre con OmniRoute

## ⚙️ Configuración

Edita `configs/config.yaml` para personalizar:

- URL y credenciales de OmniRoute
- Tema de la UI (dark/light)
- Timeouts de ejecución
- Nivel de logging
- Límites de historial

Edita `configs/tasks.yaml` para:

- Añadir nuevas tareas
- Modificar scripts existentes
- Configurar prompts de IA
- Cambiar prioridades

## 🏗️ Arquitectura

Brain CLI sigue **Clean Architecture** con capas claramente separadas:

```
internal/
├── core/        # Dominio (entidades e interfaces)
├── usecases/    # Lógica de aplicación
├── adapters/    # Implementaciones (CLI, OmniRoute, Config)
└── ui/          # Presentación (Bubble Tea TUI)
```

Ver [AGENTS.md](./AGENTS.md) para detalles completos de arquitectura y lógica de negocio.

Ver [docs/plans/development-roadmap.md](./docs/plans/development-roadmap.md) para el plan de desarrollo detallado.

## 📦 Estructura del Proyecto

```
brain-cli/
├── cmd/tui-assistant/     # Entry point
├── internal/              # Código principal
│   ├── core/             # Domain layer
│   ├── usecases/         # Application layer
│   ├── adapters/         # Infrastructure layer
│   └── ui/               # Presentation layer
├── scripts/              # Scripts de automatización
├── configs/              # Archivos de configuración
├── docs/                 # Documentación
└── pkg/                  # Paquetes reutilizables
```

## 🔧 Añadir una Nueva Tarea

1. Edita `configs/tasks.yaml`:

```yaml
- id: "mi-tarea"
  name: "Mi Tarea Nueva"
  description: "Descripción de qué hace"
  type: "script"  # script | command | ai
  script_path: "scripts/mi-script.sh"
  tags: ["automation"]
  priority: "medium"
```

2. Si es un script, créalo en `scripts/`:

```bash
#!/bin/bash
# scripts/mi-script.sh
echo "Ejecutando mi tarea..."
# Tu lógica aquí
```

3. Reinicia Brain CLI y la tarea aparecerá automáticamente

## 🧪 Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...
```

## 📚 Documentación

- [AGENTS.md](./AGENTS.md) - Arquitectura completa y lógica de negocio
- [docs/plans/development-roadmap.md](./docs/plans/development-roadmap.md) - Plan de desarrollo detallado
- [configs/config.yaml](./configs/config.yaml) - Opciones de configuración
- [configs/tasks.yaml](./configs/tasks.yaml) - Definición de tareas

## 🛣️ Roadmap

### Hito 1: Fundación (Semanas 1-2) - 🏗️ En progreso
- [x] Estructura de carpetas
- [x] Documentación inicial
- [x] Configuración con go.mod
- [ ] "Hello World" con Bubble Tea

### Hito 2: Core y Tareas Básicas (Semanas 3-4)
- [ ] Entidades del dominio
- [ ] Ejecución de scripts
- [ ] UI básica con lista de tareas

### Hito 3: Integración con OmniRoute (Semanas 5-6)
- [ ] Cliente HTTP para OmniRoute
- [ ] Tareas con IA
- [ ] Streaming de respuestas

### Hito 4: UI Avanzada (Semanas 7-8)
- [ ] Navegación multi-pantalla
- [ ] Historial de ejecuciones
- [ ] Sistema de notificaciones

### Hito 5: Pulido y Lanzamiento (Semanas 9-10)
- [ ] Suite completa de tareas
- [ ] Testing exhaustivo
- [ ] Binarios para distribución

## 🤝 Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

MIT License - Ver [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Nero**

- GitHub: [@nero](https://github.com/nero)

## 🙏 Agradecimientos

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [OmniRoute](https://omniroute.ai) - API de IA
- Comunidad de Go por las excelentes herramientas

---

**Estado del proyecto**: 🏗️ En desarrollo activo  
**Última actualización**: 2026-09-05
