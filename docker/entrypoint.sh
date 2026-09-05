#!/bin/sh
# ============================================
# brain-cli Entrypoint Script
# ============================================
# This script runs before the main application
# to perform initialization tasks
# ============================================

set -e

echo "🧠 Iniciando brain-cli..."
echo "Versión: 2.0.0"
echo "Fecha: $(date)"
echo ""

# ============================================
# 1. Verificar directorios necesarios
# ============================================
echo "📁 Verificando directorios..."
mkdir -p /app/data /app/scripts /app/logs

# ============================================
# 2. Verificar permisos
# ============================================
echo "🔐 Verificando permisos..."
if [ ! -w /app/data ]; then
    echo "❌ Error: No hay permisos de escritura en /app/data"
    exit 1
fi

# ============================================
# 3. Verificar si existe la base de datos
# ============================================
if [ ! -f "$DB_PATH" ]; then
    echo "📦 Base de datos no encontrada. Se creará automáticamente al iniciar."
else
    echo "✅ Base de datos encontrada: $DB_PATH"
    
    # Mostrar tamaño de la base de datos
    DB_SIZE=$(du -h "$DB_PATH" | cut -f1)
    echo "   Tamaño: $DB_SIZE"
fi

# ============================================
# 4. Verificar configuración
# ============================================
if [ -f /app/configs/config.local.yaml ]; then
    export CONFIG_PATH=/app/configs/config.local.yaml
    echo "📝 Usando configuración local"
elif [ -f "$CONFIG_PATH" ]; then
    echo "📝 Usando configuración: $CONFIG_PATH"
else
    echo "⚠️  Advertencia: No se encontró archivo de configuración"
fi

# ============================================
# 5. Verificar variable de encriptación
# ============================================
if [ -z "$ENCRYPTION_KEY" ] || [ "$ENCRYPTION_KEY" = "change-this-secret-key-in-production" ]; then
    echo "⚠️  ADVERTENCIA: ENCRYPTION_KEY no está configurada o usa el valor por defecto"
    echo "   Las API keys no estarán seguras. Configura ENCRYPTION_KEY en producción."
fi

# ============================================
# 6. Mostrar configuración del entorno
# ============================================
echo ""
echo "⚙️  Configuración del entorno:"
echo "   DB_PATH: $DB_PATH"
echo "   CONFIG_PATH: $CONFIG_PATH"
echo "   LOG_LEVEL: $LOG_LEVEL"
echo "   TZ: $TZ"
echo ""

# ============================================
# 7. Ejecutar migraciones (si es necesario)
# ============================================
# Las migraciones se ejecutarán automáticamente
# desde el código Go al iniciar

# ============================================
# 8. Ejecutar el comando principal
# ============================================
echo "🚀 Iniciando brain-cli..."
echo "=========================================="
echo ""

# Ejecutar el comando pasado como argumentos
exec "$@"
