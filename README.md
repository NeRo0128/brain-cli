# Brain CLI

> Asistente personal de terminal con inteligencia artificial y automatización de sistemas.

Brain CLI es una interfaz de línea de comandos (TUI) construida en Go que centraliza la ejecución de tareas del sistema, la gestión de herramientas y la interacción con proveedores de IA — todo desde la terminal, con una interfaz responsive y minimalista.

---

## Características

- **Ejecución de tareas** — Scripts bash, comandos del sistema y operaciones con IA, ejecutados desde una lista centralizada.
- **Gestión de herramientas** — CRUD completo de scripts y comandos reutilizables con soporte para bash, python, go y comandos nativos.
- **Multi-proveedor de IA** — Integración con OmniRoute, Ollama y compatibilidad con OpenAI/DeepSeek/Anthropic.
- **Historial persistente** — Cada ejecución se registra en SQLite con output, errores y métricas de rendimiento.
- **6 temas de color** — Brain (default), Catppuccin, Tokyo Night, Nord, Rose Pine y Kanagawa. Cambio en tiempo real con `Ctrl+T`.
- **Iconos nerd-font** — Glyphs para estados de ejecución, badges de tipo y prioridad.
- **Diseño responsive** — Layout adaptativo con breakpoints para terminales pequeñas (S<80), medianas (M 80-119) y grandes (L>=120).
- **Hotkeys configurables** — Bindings reasignables desde archivo de configuración YAML.

---

## Instalación

### Requisitos previos

- **Go 1.23+**
- **Bash** (para ejecutar scripts del sistema)
- **Nerd Font** (opcional, para iconos — ej. JetBrainsMono Nerd Font)

### Compilar desde fuente

```bash
git clone https://github.com/NeRo0128/brain-cli.git
cd brain-cli
go mod download
go build -o brain-cli cmd/brain-cli/main.go
./brain-cli
```

### Modo desarrollo

```bash
go run cmd/brain-cli/main.go
```

### Docker

```bash
docker build -f docker/Dockerfile -t brain-cli .
docker compose -f docker/docker-compose.yml up -d
```

---

## Uso

Al iniciar Brain CLI se presenta la **pantalla principal** con la lista de tareas disponibles. Navega con las teclas, selecciona una tarea y presiona `Enter` para ejecutarla.

```
┌──────────────────────────────────────────────────────────┐
│  🧠 Brain CLI v2.0           ● IA ready   12 tasks   3★ │
├──────────────────────────────────────────────────────────┤
│                                                          │
│   Conectar WiFi + VPN          [bash]     ★★★           │
│   Levantar Stack Dev           [bash]     ★★☆           │
│   Revisar Correos              [ai]       ★★★           │
│   Limpiar Temporales          [command]   ★☆☆           │
│   Resumir URL                  [ai]       ★★☆           │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  enter ejecutar  n nueva  d detalle  h historial  ? ayuda│
└──────────────────────────────────────────────────────────┘
```

### Flujo típico

1. Selecciona una tarea con `↑↓` o `j/k`
2. Presiona `Enter` para ejecutar
3. Observa el progreso en tiempo real con spinner y barra de progreso
4. Al finalizar, revisa el resultado con scroll
5. Presiona `r` para re-ejecutar o `Esc` para volver

---

## Atajos de teclado (Hotkeys)

### Globales

Disponibles desde cualquier pantalla.

| Tecla | Acción |
|-------|--------|
| `↑` / `k` | Mover cursor arriba |
| `↓` / `j` | Mover cursor abajo |
| `PgUp` | Página arriba |
| `PgDn` | Página abajo |
| `g` | Ir al inicio de la lista |
| `G` | Ir al final de la lista |
| `Enter` | Confirmar / Ejecutar |
| `Esc` | Volver / Cancelar |
| `/` | Filtrar lista |
| `Ctrl+T` | Cambiar tema de color |
| `Ctrl+C` | Salir de la aplicación |
| `q` | Salir (fuera de ejecución) |
| `?` | Abrir pantalla de ayuda |

### Pantalla Principal

| Tecla | Acción |
|-------|--------|
| `Enter` | Ejecutar tarea seleccionada |
| `n` | Crear nueva tarea |
| `e` | Editar tarea seleccionada |
| `Ctrl+D` | Eliminar tarea |
| `d` | Ver detalle de la tarea |
| `h` | Abrir historial de ejecuciones |

### Pantalla de Detalle

| Tecla | Acción |
|-------|--------|
| `Enter` | Ejecutar tarea |
| `e` | Editar tarea |
| `Esc` | Volver a la lista |

### Pantalla de Resultado

| Tecla | Acción |
|-------|--------|
| `r` | Re-ejecutar la tarea |
| `Esc` | Volver |

### Pantalla de Ejecución

| Tecla | Acción |
|-------|--------|
| `Esc` | Cancelar ejecución en curso |

### Pantalla de Historial

| Tecla | Acción |
|-------|--------|
| `Enter` | Ver detalle de la ejecución seleccionada |
| `Esc` | Volver a la pantalla principal |

### Formulario de Tarea / Herramienta

| Tecla | Acción |
|-------|--------|
| `Tab` / `Shift+Tab` | Mover foco entre campos |
| `←` / `→` | Cambiar valor en selects (tipo, prioridad, etc.) |
| `Ctrl+S` | Guardar cambios |
| `Esc` | Cancelar y volver |

---

## Configuración

Brain CLI se configura desde `configs/config.yaml`. Las variables de entorno se resuelven automáticamente con soporte para valores por defecto (`${VAR:-default}`).

### Estructura del archivo

```yaml
app:
  name: brain-cli

database:
  path: ${BRAIN_DB_PATH:-data/brain.db}
  auto_migrate: true

logging:
  level: ${BRAIN_LOG_LEVEL:-info}        # debug | info | warn | error
  format: ${BRAIN_LOG_FORMAT:-pretty}    # json | pretty
  output_path: /tmp/brain-cli.log
  max_size_mb: 100

ui:
  theme: ${BRAIN_THEME:-brain}           # brain | catppuccin | tokyo-night | nord | rose-pine | kanagawa
  brand_style: ascii                     # ascii | minimal | none
  icons: ${BRAIN_ICONS:-nerd-b}          # unicode | nerd-b | nerd-c

scripts:
  custom_dir: scripts
```

### Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `BRAIN_DB_PATH` | `data/brain.db` | Ruta al archivo SQLite |
| `BRAIN_LOG_LEVEL` | `info` | Nivel de logging |
| `BRAIN_LOG_FORMAT` | `pretty` | Formato de logs |
| `BRAIN_THEME` | `brain` | Tema de color activo |
| `BRAIN_ICONS` | `nerd-b` | Set de iconografía |
| `BRAIN_CONFIG` | `configs/config.yaml` | Ruta al archivo de configuración |

---

## Temas

Brain CLI incluye 6 paletas de color. Cambia entre ellas en cualquier momento con `Ctrl+T`.

| Tema | Estilo | Colores principales |
|------|--------|-------------------|
| **brain** | Default, dark-only | Cyan primary, violet secondary, rose accents |
| **catppuccin** | Macchiato (dark) / Latte (light) | Mauve primary, pink secondary |
| **tokyo-night** | Dark-only | Blue primary, magenta secondary |
| **nord** | Dark-only | Frost blue primary, aurora green accents |
| **rose-pine** | Moon (dark) / Dawn (light) | Pine primary, gold secondary |
| **kanagawa** | Dark-only | Wave blue primary, carp red accents |

### Brand Style

El header de la aplicación soporta 3 estilos de marca:

- `ascii` — Texto plano `BrainCLI` con color de acento
- `minimal` — Icono + nombre `🧠 Brain CLI` (default)
- `none` — Sin marca visible

---

## Arquitectura

Brain CLI sigue **Clean Architecture** con dependencias unidireccionales hacia el dominio.

```
┌─────────────────────────────────────────────────────┐
│  UI (internal/ui/)                                  │
│  Bubble Tea TUI · Lipgloss · Bubbles                │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│  Use Cases (internal/usecases/)                     │
│  Task CRUD · Execution · AI · Provider · Tool       │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│  Domain (internal/core/)                            │
│  Task · Execution · Provider · Tool · Config        │
└──────────────────────▲──────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────┐
│  Adapters (internal/adapters/)                      │
│  SQLite · Executor · AI Providers · Email · Config  │
└─────────────────────────────────────────────────────┘
```

### Paquetes principales

| Paquete | Responsabilidad |
|---------|----------------|
| `internal/core/` | Entidades, interfaces de repositorio, reglas de negocio |
| `internal/usecases/` | Lógica de aplicación (CRUD, ejecución, IA) |
| `internal/adapters/` | Implementaciones concretas (DB, executor, proveedores) |
| `internal/ui/` | Presentación TUI, pantallas, componentes, estilos |
| `internal/ui/theme/` | 6 paletas de color con soporte dark/light |
| `internal/ui/icons/` | Sets de iconografía (unicode, nerd-font) |
| `internal/ui/styles/` | Estilos Lipgloss derivados del tema activo |
| `pkg/utils/` | Paquetes reutilizables (logger, validador, formateador) |

---

## Desarrollo

### Ejecutar en desarrollo

```bash
go run cmd/brain-cli/main.go
```

### Ejecutar con debug

```bash
DEBUG=1 go run cmd/brain-cli/main.go
```

### Tests

```bash
go test ./...                         # todos los tests
go test -cover ./...                  # con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out      # reporte HTML
go test -bench=. ./...                # benchmarks
```

### Lint y verificación

```bash
go vet ./...
go fmt ./...
goimports -w .
```

### Previsualizar iconos

```bash
go run cmd/icon-preview/main.go
```

---

## Estructura del proyecto

```
brain-cli/
├── cmd/
│   ├── brain-cli/main.go            # Entry point principal
│   └── icon-preview/main.go         # Dev tool: preview de iconos
├── internal/
│   ├── core/                         # Domain Layer
│   │   ├── config/                   # Config entity + validation
│   │   ├── execution/                # Execution entity
│   │   ├── provider/                 # AI Provider entity
│   │   ├── task/                     # Task entity
│   │   └── tool/                     # Tool entity
│   ├── usecases/                     # Application Layer
│   │   ├── task/                     # Task CRUD + execution
│   │   ├── tool/                     # Tool CRUD
│   │   └── provider/                 # Provider management
│   ├── adapters/                     # Infrastructure Layer
│   │   ├── database/                 # SQLite + migrations + repos
│   │   ├── executor/                 # Script/command executors
│   │   ├── ai/                       # AI provider implementations
│   │   └── config/                   # YAML config loader
│   └── ui/                           # Presentation Layer
│       ├── app.go                    # Root Bubble Tea model
│       ├── update.go                 # Message routing
│       ├── view.go                   # View composition
│       ├── model.go                  # App state
│       ├── screens/                  # 10 screens (main, detail, result, etc.)
│       ├── components/               # Reusable components (toast, header, list)
│       ├── styles/                   # Lipgloss styles from theme
│       ├── theme/                    # 6 color palettes
│       ├── icons/                    # Nerd-font icon sets
│       └── keys/                     # Keybinding registry
├── configs/config.yaml               # Configuración
├── docker/                           # Docker infrastructure
├── docs/                             # Documentación y planes
└── pkg/utils/                        # Public packages
```

---

## Widget flotante (Niri + Ghostty)

Brain CLI puede ejecutarse como un **widget de escritorio desplegable** al estilo Quake/Guake. Se abre con un atajo de teclado, aparece centrado flotando sobre el escritorio, y se cierra desde la propia TUI.

**Flujo:** `Super+B` → Brain CLI aparece → trabajas → `q` → se cierra → `Super+B` → se reabre.

### Requisitos

| Componente | Propósito |
|------------|-----------|
| **Niri** | Compositor Wayland (soporte flotante) |
| **Ghostty** | Terminal con config dedicada |
| **`jq`** | Parsear JSON de `niri msg` |

```bash
command -v niri ghostty jq  # verificar dependencias
```

### Ghostty: config dedicada

Crea un archivo de configuración exclusivo para el widget. Esto evita que las opciones afecten a tus terminales normales.

**`~/.config/ghostty/config-drop`**:

```ini
# Título único — Niri lo usa como identificador de ventana.
title = BrainCLIDrop

# Sin decoraciones de ventana (más limpio para flotante).
window-decoration = false

# Tamaño en celdas (103x32 ≈ 800x700 px con fuente por defecto).
window-width = 103
window-height = 32

# Fondo 100% opaco para legibilidad.
background-opacity = 1
```

> El título es la clave: Ghostty/GTK4 en Wayland ignora `--class` y `GDK_APP_ID`. El título sí se respeta y Niri puede matchearlo.

### Niri: window-rule

Añade esta regla a `~/.config/niri/config.kdl`:

```kdl
window-rule {
    match title="BrainCLIDrop"
    open-floating true
    opacity 1.0
    default-column-width { fixed 800; }
    default-window-height { fixed 700; }
}
```

| Campo | Descripción |
|-------|-------------|
| `match title=` | Matchea por título (no `app-id`, que Ghostty no respeta) |
| `open-floating true` | Abre la ventana flotante |
| `opacity 1.0` | Sobreescribe reglas globales de transparencia |
| `default-column-width` | Ancho fijo en píxeles |
| `default-window-height` | Alto fijo en píxeles |

Valida y recarga:

```bash
niri validate
niri msg action load-config-file
```

### Script de toggle

Crea `~/.local/bin/tdrop-niri.sh`:

```bash
#!/bin/bash
# tdrop-niri.sh — Drop-down toggle para Brain CLI en Niri.
set -e

DROP_TITLE="BrainCLIDrop"
BRAIN_DIR="${BRAIN_DIR:-$HOME/brain-cli}"
BRAIN_CLI="$BRAIN_DIR/build/brain-cli"
GHOSTTY_CONFIG="${GHOSTTY_CONFIG:-$HOME/.config/ghostty/config-drop}"

if [ ! -x "$BRAIN_CLI" ]; then
    echo "brain-cli no encontrado en $BRAIN_CLI"
    echo "Compila: cd $BRAIN_DIR && go build -o build/brain-cli ./cmd/brain-cli"
    exit 1
fi

WINDOW_JSON=$(niri msg --json windows 2>/dev/null \
    | jq -c --arg title "$DROP_TITLE" '.[] | select(.title == $title)')

# No existe → lanzar
if [ -z "$WINDOW_JSON" ]; then
    niri msg action spawn -- \
        ghostty \
        --config-file="$GHOSTTY_CONFIG" \
        -e bash -c "cd '$BRAIN_DIR' && exec '$BRAIN_CLI'"
    exit 0
fi

# Existe → enfocar
WIN_ID=$(echo "$WINDOW_JSON" | jq -r '.id')
niri msg action focus-window --id "$WIN_ID"
```

```bash
chmod +x ~/.local/bin/tdrop-niri.sh
```

### Atajo de teclado

En `~/.config/niri/keybinds.kdl`:

```kdl
Mod+B { spawn-sh "/home/$USER/.local/bin/tdrop-niri.sh"; }
```

> Usa ruta absoluta — Niri no expande `~` dentro de `spawn-sh`.

### Evitar que scripts roben el foco

Si desde Brain CLI lanzas apps (editores, terminales, navegadores), puedes evitar que roben el foco. El truco: capturar el foco antes y restaurarlo después.

**`~/.local/bin/lib/refocus.sh`**:

```bash
#!/bin/bash
capture_focus() {
    if command -v niri >/dev/null 2>&1 && command -v jq >/dev/null 2>&1; then
        niri msg --json focused-window 2>/dev/null | jq -r '.id // empty'
    fi
}

restore_focus() {
    local win_id="$1"
    if [ -n "$win_id" ]; then
        sleep 1
        niri msg action focus-window --id "$win_id" 2>/dev/null || true
    fi
}
```

Uso en cualquier script:

```bash
source ~/.local/bin/lib/refocus.sh
FOCUS=$(capture_focus)

setsid my-editor </dev/null >/dev/null 2>&1 &
setsid my-terminal </dev/null >/dev/null 2>&1 &

restore_focus "$FOCUS"
```

`setsid` crea una nueva sesión que sobrevive al padre, evitando que los hijos reciban `SIGHUP` al terminar el script.

### Ejemplo: entorno dev sin perder foco

**`scripts/open-dev.sh`**:

```bash
#!/bin/bash
set -e
WORK_DIR="${1:-$PWD}"

source ~/.local/bin/lib/refocus.sh
FOCUS=$(capture_focus)

setsid zed "$WORK_DIR" </dev/null >/dev/null 2>&1 &
setsid ghostty --working-directory="$WORK_DIR" </dev/null >/dev/null 2>&1 &

restore_focus "$FOCUS"
echo "Entorno lanzado en $WORK_DIR"
```

Regístralo como Tool en Brain CLI (`n` → tipo bash), crea una Task que lo use. Al ejecutarlo desde la TUI: Zed + Ghostty se abren, Brain CLI mantiene el foco, `q` cierra la TUI y Zed queda listo.

### Troubleshooting

| Síntoma | Causa | Fix |
|---------|-------|-----|
| `niri validate` falla | Usaste `#` como comentario | KDL solo acepta `//` o `/* */` |
| Ventana no es flotante | El título no matchea | `niri msg --json windows \| jq '.[].title'` para ver títulos reales |
| Ventana no se centra | Niri recuerda último tamaño | Redimensiónala, ciérrala, reábrela |
| Tamaño del config no se aplica | Niri guarda tamaño por ventana | Cambia el título (`BrainCLIDropV2`) y matchea el nuevo |
| App lanzada roba el foco | Comportamiento por defecto | Usa `capture_focus` / `restore_focus` |

---

## Contribuir

1. Fork el repositorio
2. Crea una rama (`git checkout -b feature/nueva-funcionalidad`)
3. Haz commit con convención (`feat:`, `fix:`, `refactor:`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

---

## Licencia

MIT — Ver [LICENSE](LICENSE).

---

**Autor:** [Nero](https://github.com/NeRo0128) · **Versión:** 2.0.0 · **Estado:** En desarrollo activo
