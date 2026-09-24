# Changelog

Todas las versiones notables de Brain CLI.
Formato basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/).

## [1.0.0] — 2026-09-23

### 🎉 Primera versión estable

**Core**

- Ejecución de tasks con bash, python y comandos nativos
- Tool system con soft-delete y versionado
- Ejecutor con timeouts, cancelación y captura de stdout/stderr
- Persistencia completa en SQLite con migraciones embebidas

**TUI (Bubble Tea v2)**

- Lista de tasks navegable con badges (tipo, prioridad)
- Detalle en 1 o 2 columnas (responsive)
- Formularios de creación/edición con textareas y selectores
- Historial de ejecuciones con filtrado
- Pantalla de resultados con scroll nativo
- Spinner + barra de progreso durante ejecuciones

**Apariencia**

- 6 temas: brain, catppuccin, tokyo-night, nord, rose-pine, kanagawa
- Cambio de tema en runtime con `Ctrl+T`
- 2 sets de iconos: Unicode y Nerd Font
- Header ASCII configurable (minimal / slim / big)
- Sistema de toasts no bloqueante

**Settings**

- Pantalla de ajustes con 4 tabs
- Persistencia de cambios en SQLite (tabla `settings`)
- Jerarquía de config: defaults → YAML → DB → env vars

**Auth**

- Login con GitHub mediante OAuth Device Flow
- Sesión cifrada en el keyring del sistema
- Integración en Settings

**CLI**

- Bootstrap autosuficiente: el binario crea config, DB y directorios XDG solo
- Detecta modo dev (go.mod presente) vs prod (~/.config/brain-cli/)
- Overrides por env: `BRAIN_CONFIG`, `BRAIN_DB_PATH`, `GITHUB_CLIENT_ID`

**Integraciones**

- Widget flotante en Niri + Ghostty

[1.0.0]: https://github.com/NeRo0128/brain-cli/releases/tag/v1.0.0
