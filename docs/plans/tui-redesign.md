# Plan de Mejora UI – Brain CLI TUI Profesional

**Fecha:** 2026-09-13
**Decisiones:** minimal premium, responsive, intuitivo + vistoso
**Skills aplicadas:** `tui-design` + `gentleman-bubbletea`
**Estado:** aprobado para implementar por fases

## 0. Objetivo

Dejar Brain CLI con estética minimal premium, responsive al tamaño de terminal,
intuitiva y vistosa, sin romper arquitectura `ScreenI` ni lógica de negocio.

Solo se toca `internal/ui/` + `styles`. No se tocan DB, usecases, providers IA, config YAML.

## 1. Diagnóstico (nota 6.2/10)

- Doble título compitiendo (`🧠 Brain` + `list.Title`), sin layout global.
- Paleta 7 colores saturados, fatiga visual. Foco solo por color, invisible en tema claro.
- Sin bordes/paneles, todo texto suelto. `maxLines=15` fijo, `height-6` mágico, se rompe en resize.
- Delegates duplicados x3 con micro-diferencias.
- Footers hardcodeados desincronizados de `keys.Registry`.
- Toast actual fondo sólido que grita + `Overlay()` por concatenación que se corta si la línea base ocupa todo el ancho. Sin icono/título/timeout.
- Executing solo `spinner.Dot`, parece colgado en tasks largas.
- Result con `scroll int` manual, sin viewport, sin re-ejecutar.
- Form 9 campos idénticos, foco sutil, prompt 1 línea, error solo al final.
- `app.go` mezcla nav + exec + Key→Action + toast fan-out.

## 2. Dirección de diseño

- Monocromo base + 1 acento. Máx 4 colores.
- `Text #E5E7EB, Muted #6B7280, Primary Adaptive Light #4C1D95 / Dark #A78BFA, bg seleccionado #1F2937`.
- `Success/Warning/Error` solo para iconos de estado.
- Bordes `Rounded Muted` solo para agrupar. Separación normal = espaciado.
- Foco = `fondo + ▶ + bold`, nunca solo color.
- Footer único desde `Registry`. Header único con estado IA.
- Breakpoints: `S <80 1col compacto, M 80-119 1col aireado, L >=120 2col`.
- Gate si `width<70` o `height<20`.

Mockup frame:

```
┌────────────────────────────────────────────────────────┐
| 🧠 Brain CLI v2.0               ● IA ready  12 tasks 3★ |
├────────────────────────────────────────────────────────┤
| CONTENIDO SCREEN                                       |
├────────────────────────────────────────────────────────┤
| ↑↓ nav  Enter run  d detalle  ? ayuda  q salir         |
└────────────────────────────────────────────────────────┘
```

Toast objetivo top-right:

```
                                               ┌─ ✓ Guardado ───────┐
                                               │ Task guardada      │
                                               │ ━━━━━━━━░░░░ 2s   │
                                               └────────────────────┘
```

Executing objetivo:

```
◐ Ejecutando Backup diario        00:12
  › dumping db…
  ━━━━━━━━━━╺━━━━━━━ 45%  1.2 MB/s
  Esc cancelar
```

## 3. Fases

### S0 – Base
**Por qué:** desbloquea todo, cero riesgo lógica.
- `styles/tokens.go`: Primary/Text/Muted/SelectedRow/Border/InputFocused/InputBlurred/StatusBar + AdaptiveColor.
- `styles/status.go`: `StatusIcon(), Badge(), HumanTime()`.
- `layout/responsive.go`: `Breakpoint(), ContentSize(), ShouldGate()`.
- `components/frame/header.go footer.go`: header IA status, footer desde Registry.
- `app/view.go`: frame + gate + error screen con borde rojo `e logs / r reintentar / Esc`.
- Verif: `go build ./...`, probar 70x20,100x30,160x50, claro/oscuro.

### S1 – Toasts premium
**Por qué:** quick win visual, feedback intuitivo sin leer.
- Split `components/toast/model.go|view.go|overlay.go`.
- Caja 30w, RoundedBorder, bg #111827, borde-izq + icono color, título bold + msg Muted 2 líneas + barra timeout 2px.
- Overlay con `lipgloss.Place` top-right, `colStart=width-toastW-2`. Si `width<80` → top-center.
- Max 3, error 5s resto 3s.
- Verif: guardar + ejecutar, comprobar no tapa lista.

### S2 – Main + History
**Por qué:** mayor impacto vistoso.
- `components/list/delegate.go` único reemplaza 3 copias.
- Fila: `▶ ★ Nombre [TYPE][PRIO]` + `✓ hace 2h · 1.2s`. Badges tenues.
- Main empty `Sin tareas — n crear`, keys `x favorito, ctrl+r`.
- History agrupado `HOY/AYER`, tiempo relativo, `tab filtra todos/✓/✗/⊘`.
- Verif: `go test ./internal/ui/...` + screenshots.

### S3 – Result + Executing + Progress
**Por qué:** flujo core ejecutar→ver→reintentar.
- Result: `scroll int` → `viewport.Model`, caja bordeada, highlight error|fail, `r re-ejecutar, c copiar`.
- `components/progress/bar.go`: wrapper `progress Gradient #7C3AED→#06B6D4 + spinner.Jump`, `SetPercent/SetStage`, tail 6 líneas.
- `app/exec.go`: extraer startExecution/runTask/handleFinished.
- Verif: task larga, Esc cancela, r re-ejecuta.

### S4 – Detail + Form + Picker/Help + split final
**Por qué:** cierra intuitivo + deuda arquitectura.
- Detail: `>=100` 2 col meta|preview glamour, `<100` apilado. Tags pills.
- Form: inputs con borde, textarea 5 líneas prompt, validación inline, steps, up/down alias tab.
- Picker badge `★ actual`. Help 2 col por Grouped, `?/Esc` cierra.
- `app/update.go` split updateNav/translateKey/delegateTop. `components/states/empty.go loading.go error.go`.
- Desambiguar `enter`: NavConfirm solo pickers, ActionExecute solo main/detail.
- Verif: `go vet ./...`, teatest, checklist tui-design.

## 4. Archivos a modificar

**Modificar:**
- `internal/ui/app.go` → dividir en `internal/ui/app/model.go|update.go|exec.go|view.go`
- `internal/ui/styles/styles.go` → ampliar a `tokens.go + status.go`
- `internal/ui/screens/main.go` – delegate + badges + empty + Keys
- `internal/ui/screens/history.go` – agrupado + relativo + filtro
- `internal/ui/screens/results.go` – viewport + re-run
- `internal/ui/screens/executing.go` – progress + log tail
- `internal/ui/screens/details.go` – responsive 1/2 col + glamour
- `internal/ui/screens/form.go` – bordes + textarea + validación
- `internal/ui/screens/tool_picker.go` – badge actual
- `internal/ui/screens/help.go` – 2 columnas
- `internal/ui/components/toast.go` – split + estilo minimal + overlay Place
- `internal/ui/keys/defaults.go` – desambiguar enter

**Crear:**
- `internal/ui/layout/responsive.go`
- `internal/ui/components/frame/header.go`
- `internal/ui/components/frame/footer.go`
- `internal/ui/components/toast/model.go|view.go|overlay.go` (split del actual)
- `internal/ui/components/progress/bar.go`
- `internal/ui/components/list/delegate.go`
- `internal/ui/components/states/empty.go|loading.go|error.go`
- `internal/ui/app/model.go|update.go|exec.go|view.go` (split de app.go)

**No tocar:**
- `internal/core/*`, `internal/usecases/*`, `internal/adapters/*`, `configs/*`, `go.mod` (bubbles ya incluye progress).

## 5. Criterio de done

Build verde, footers = Registry, foco visible claro/oscuro, toasts arriba-derecha con timeout, progress visible, responsive 70/100/160, sin maxLines fijos ni delegates duplicados.
