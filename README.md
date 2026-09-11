# 🧠 Brain CLI

> Asistente personal de línea de comandos con capacidades de IA

Brain CLI es un asistente TUI (Terminal User Interface) moderno construido en Go que automatiza tus tareas diarias del sistema y proporciona capacidades de inteligencia artificial.

## ✨ Características

- ⚡ **Automatización Inteligente**: Ejecuta scripts y comandos del sistema organizadamente
- 🤖 **Integración con IA**: Multi-proveedor (OmniRoute, Ollama, OpenAI, DeepSeek)
- 🎨 **Interfaz Elegante**: TUI moderna con Bubble Tea
- 📊 **Historial Completo**: Guarda todas las ejecuciones en SQLite
- 🔧 **Extensible**: Añade nuevas tareas y herramientas fácilmente
- ⚙️ **Dockerizado**: Ejecuta en contenedores con persistencia
- 📦 **Scripts Embebidos**: Binario autosuficiente

## 🚀 Quick Start

### Prerrequisitos

- Go 1.23 o superior
- Docker (opcional, para ejecución en contenedor)
- Bash (para ejecutar scripts)

### Instalación

```bash
# Clonar el repositorio
git clone https://github.com/NeRo0128/brain-cli.git
cd brain-cli

# Instalar dependencias
go mod download

# Compilar
go build -o brain-cli cmd/brain-cli/main.go

# Ejecutar
./brain-cli
```

### Desarrollo

```bash
# Ejecutar en modo desarrollo
go run cmd/brain-cli/main.go

# Ejecutar con debug logging
DEBUG=1 go run cmd/brain-cli/main.go

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...
```

### Docker

```bash
# Build image
docker build -f docker/Dockerfile -t brain-cli .

# Run with docker-compose
docker-compose -f docker/docker-compose.yml up -d

# View logs
docker-compose -f docker/docker-compose.yml logs -f brain-cli
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
- **Chat con IA**: Conversación libre con el proveedor configurado

## ⚙️ Configuración

### Variables de Entorno

```bash
# Database
DB_PATH=/app/data/brain.db

# OmniRoute connection
OMNIROUTE_URL=http://omniroute:20128
OMNIROUTE_API_KEY=your-api-key

# Encryption
ENCRYPTION_KEY=your-encryption-key
```

### Archivos de Configuración

- `configs/config.yaml`: Configuración general de la aplicación
- `embed/migrations/`: Migraciones de base de datos
- `embed/scripts/`: Scripts embebidos en el binario

## 🏗️ Arquitectura

Brain CLI sigue **Clean Architecture** con capas claramente separadas:

```
internal/
├── core/        # Dominio (entidades e interfaces)
├── usecases/    # Lógica de aplicación
├── adapters/    # Implementaciones (SQLite, AI Providers, Executor)
└── ui/          # Presentación (Bubble Tea TUI)
```

### Entidades Principales

- **Task**: Unidad de trabajo ejecutable
- **Execution**: Instancia de ejecución de una tarea
- **Provider**: Proveedor de IA configurado
- **Tool**: Herramienta/script ejecutable

Ver [AGENTS.md](./AGENTS.md) para detalles completos de arquitectura y lógica de negocio.

## 📦 Estructura del Proyecto

```
brain-cli/
├── cmd/brain-cli/              # Entry point
├── internal/
│   ├── core/                   # Domain Layer
│   ├── usecases/              # Application Layer
│   ├── adapters/              # Infrastructure Layer
│   └── ui/                    # Presentation Layer
├── embed/                      # Archivos embebidos
│   ├── scripts/               # Scripts por defecto
│   └── migrations/            # Migraciones SQL
├── docker/                     # Infraestructura Docker
├── configs/                    # Configuración
├── data/                       # Base de datos SQLite
├── scripts/                    # Scripts personalizados
├── docs/                       # Documentación
└── pkg/                        # Paquetes reutilizables
```

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
- [docs/MIGRATION_PLAN.md](./docs/MIGRATION_PLAN.md) - Plan de migración v1 → v2
- [DEVELOPMENT.md](./docs/DEVELOPMENT.md) - Guía de desarrollo

## 🛣️ Roadmap

### Hito 1: Fundación (Semanas 1-2)
- [x] Estructura de carpetas
- [x] Documentación inicial
- [x] Configuración con go.mod
- [x] Dockerización
- [ ] "Hello World" con Bubble Tea

### Hito 2: Core y Tareas Básicas (Semanas 3-4)
- [ ] Entidades del dominio
- [ ] Ejecución de scripts
- [ ] UI básica con lista de tareas

### Hito 3: Integración con IA (Semanas 5-6)
- [ ] Cliente HTTP para OmniRoute
- [ ] Soporte Ollama
- [ ] Tareas con IA

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
- GitHub: [@NeRo0128](https://github.com/NeRo0128)

## 🙏 Agradecimientos

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [SQLite](https://www.sqlite.org/) - Base de datos embebida
- Comunidad de Go por las excelentes herramientas

---

**Estado del proyecto**: 🏗️ En desarrollo activo  
**Versión**: 2.0.0  
**Última actualización**: 2026-09-07
