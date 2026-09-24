#!/bin/bash
# install.sh — Instalador de Brain CLI.
#
#   curl -fsSL https://raw.githubusercontent.com/NeRo0128/brain-cli/main/install.sh | bash

set -e

REPO="NeRo0128/brain-cli"
BINARY="brain-cli"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# ─── Detectar OS/arq ────────────────────────────────────────
detect_platform() {
    local os arch
    case "$(uname -s)" in
        Linux)   os="linux" ;;
        Darwin)  os="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *) echo "❌ OS no soportado: $(uname -s)"; exit 1 ;;
    esac
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) echo "❌ Arquitectura no soportada: $(uname -m)"; exit 1 ;;
    esac
    echo "${os}-${arch}"
}

PLATFORM=$(detect_platform)
EXT=""
case "$PLATFORM" in
    windows-*) EXT=".exe" ;;
esac

BIN_NAME="${BINARY}-${PLATFORM}${EXT}"

echo "🧠 Instalando Brain CLI"
echo "   Plataforma: $PLATFORM"
echo "   Destino:    $INSTALL_DIR/$BINARY"
echo ""

# ─── Última versión ─────────────────────────────────────────
echo "🔍 Buscando última versión..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    echo "❌ No se pudo obtener la última versión"
    exit 1
fi
echo "   Versión: $LATEST"
echo ""

# ─── Descargar ──────────────────────────────────────────────
URL="https://github.com/${REPO}/releases/download/${LATEST}/${BIN_NAME}"
TMP=$(mktemp)
echo "⬇️  Descargando $BIN_NAME..."
curl -fsSL "$URL" -o "$TMP" || { echo "❌ Falló la descarga desde $URL"; rm -f "$TMP"; exit 1; }

# ─── Instalar ───────────────────────────────────────────────
mkdir -p "$INSTALL_DIR"
chmod +x "$TMP"
mv "$TMP" "$INSTALL_DIR/$BINARY"
echo "✅ Instalado en $INSTALL_DIR/$BINARY"

# ─── PATH check ─────────────────────────────────────────────
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "⚠️  $INSTALL_DIR no está en tu \$PATH."
    echo "   Añade esto a tu ~/.bashrc o ~/.config/fish/config.fish:"
    echo ""
    echo "     export PATH=\"$INSTALL_DIR:\$PATH\""
    echo ""
fi

echo ""
echo "🎉 ¡Listo! Ejecuta: $BINARY"
