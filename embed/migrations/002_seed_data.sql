-- ============================================
-- brain-cli Seed Data v2.0
-- ============================================
-- Migration: 002_seed_data.sql
-- Description: Datos iniciales para brain-cli v2.0
-- Author: Nero
-- Date: 2026-09-05
-- ============================================

-- ============================================
-- Seed: Default Providers (Motores de IA)
-- ============================================
INSERT OR IGNORE INTO providers (name, type, endpoint, api_key, model, is_active, config) VALUES
    (
        'OmniRoute Local',
        'omniroute',
        'http://localhost:20128',
        NULL,  -- API key será configurada por el usuario
        'auto',
        1,  -- Proveedor activo por defecto
        '{"timeout": 60, "max_retries": 3, "retry_delay": 2}'
    ),
    (
        'Ollama Local',
        'ollama',
        'http://localhost:11434',
        NULL,
        'llama3.2',
        0,
        '{"timeout": 120, "stream": true}'
    ),
    (
        'DeepSeek Cloud',
        'deepseek',
        'https://api.deepseek.com',
        NULL,  -- Usuario debe configurar su API key
        'deepseek-chat',
        0,
        '{"timeout": 90, "max_tokens": 4000}'
    );

-- ============================================
-- Seed: Default Tools (Scripts Embebidos)
-- ============================================
INSERT OR IGNORE INTO tools (name, description, script_content, script_type, category, requires_sudo, timeout_seconds, is_builtin, version) VALUES
    (
        'wifi-vpn',
        'Conectar WiFi y levantar VPN automáticamente',
        '#!/bin/bash
# Script para conectar WiFi y VPN
# Versión: 1.0

set -e

echo "🔌 Conectando a WiFi..."
# TODO: Implementar lógica de conexión WiFi
# nmcli device wifi connect "SSID" password "PASSWORD"

echo "🔐 Levantando VPN..."
# TODO: Implementar lógica de VPN
# openvpn --config /path/to/config.ovpn --daemon

echo "✅ Conexión establecida"
exit 0',
        'bash',
        'network',
        1,  -- Requiere sudo
        120,
        1,  -- Script builtin
        1
    ),
    (
        'start-dev',
        'Levantar stack de desarrollo (OmniRoute + OpenCode)',
        '#!/bin/bash
# Script para iniciar entorno de desarrollo
# Versión: 1.0

set -e

echo "🚀 Iniciando stack de desarrollo..."

# Verificar Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado"
    exit 1
fi

echo "📦 Iniciando OmniRoute..."
# TODO: Implementar inicio de OmniRoute
# docker-compose -f ~/omniroute/docker-compose.yml up -d

echo "💻 Iniciando OpenCode..."
# TODO: Implementar inicio de OpenCode
# code ~/projects

echo "✅ Stack de desarrollo listo"
exit 0',
        'bash',
        'dev',
        0,
        180,
        1,
        1
    ),
    (
        'cleanup',
        'Limpieza de sistema: temporales, caché, logs antiguos',
        '#!/bin/bash
# Script de limpieza del sistema
# Versión: 1.0

set -e

echo "🧹 Iniciando limpieza del sistema..."

# Limpiar caché de apt (si existe)
if command -v apt-get &> /dev/null; then
    echo "  Limpiando caché de apt..."
    sudo apt-get clean
fi

# Limpiar temporales
echo "  Limpiando archivos temporales..."
sudo rm -rf /tmp/*
rm -rf ~/.cache/*

# Limpiar logs antiguos (mayores a 7 días)
echo "  Limpiando logs antiguos..."
find /var/log -type f -name "*.log" -mtime +7 -delete 2>/dev/null || true

# Docker cleanup (si está instalado)
if command -v docker &> /dev/null; then
    echo "  Limpiando contenedores Docker..."
    docker system prune -f
fi

echo "✅ Limpieza completada"
exit 0',
        'bash',
        'maintenance',
        1,  -- Requiere sudo
        300,
        1,
        1
    ),
    (
        'backup',
        'Backup rápido de archivos importantes',
        '#!/bin/bash
# Script de backup rápido
# Versión: 1.0

set -e

BACKUP_DIR="$HOME/backups/$(date +%Y%m%d_%H%M%S)"

echo "💾 Iniciando backup..."
echo "  Destino: $BACKUP_DIR"

mkdir -p "$BACKUP_DIR"

# Backup de configuraciones
echo "  Respaldando configuraciones..."
cp -r ~/.config "$BACKUP_DIR/config" 2>/dev/null || true
cp ~/.bashrc "$BACKUP_DIR/" 2>/dev/null || true
cp ~/.zshrc "$BACKUP_DIR/" 2>/dev/null || true

# Backup de proyectos (lista personalizable)
echo "  Respaldando proyectos..."
# TODO: Configurar directorios de proyectos
# rsync -av ~/projects "$BACKUP_DIR/"

echo "✅ Backup completado en: $BACKUP_DIR"
exit 0',
        'bash',
        'maintenance',
        0,
        600,
        1,
        1
    ),
    (
        'monitor',
        'Monitor de recursos del sistema (CPU, RAM, Disco)',
        '#!/bin/bash
# Monitor de recursos del sistema
# Versión: 1.0

echo "📊 Monitor de Recursos del Sistema"
echo "=================================="
echo ""

# CPU
echo "🔥 CPU:"
top -bn1 | grep "Cpu(s)" | awk ''{print "  Uso: " $2 " (user), " $4 " (system)"}''

# RAM
echo ""
echo "💾 Memoria RAM:"
free -h | awk ''NR==2{printf "  Usado: %s / %s (%.2f%%)\n", $3, $2, ($3/$2)*100}''

# Disco
echo ""
echo "💿 Disco:"
df -h / | awk ''NR==2{printf "  Usado: %s / %s (%s)\n", $3, $2, $5}''

# Procesos top
echo ""
echo "🔝 Procesos con mayor uso de CPU:"
ps aux --sort=-%cpu | head -6 | tail -5 | awk ''{printf "  %s: %.1f%%\n", $11, $3}''

exit 0',
        'bash',
        'system',
        0,
        30,
        1,
        1
    ),
    (
        'docker-ps',
        'Listar contenedores Docker activos',
        'docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"',
        'native',
        'dev',
        0,
        10,
        1,
        1
    ),
    (
        'email-check',
        'Revisar correos nuevos y clasificar con IA',
        '#!/usr/bin/env python3
# Script para revisar Gmail y clasificar correos
# Versión: 1.0

import sys
import json

def main():
    # TODO: Implementar integración con Gmail API
    print("📧 Revisando correos nuevos...")
    
    # Placeholder - será reemplazado con Gmail API
    emails = [
        {
            "from": "ejemplo@email.com",
            "subject": "Reunión importante",
            "preview": "Hola, necesitamos coordinar..."
        }
    ]
    
    # Retornar JSON para que brain-cli lo procese con IA
    print(json.dumps({"emails": emails}))
    return 0

if __name__ == "__main__":
    sys.exit(main())
',
        'python',
        'ai',
        0,
        60,
        1,
        1
    );

-- ============================================
-- Seed: Default Tasks
-- ============================================
INSERT OR IGNORE INTO tasks (name, description, tool_id, provider_id, prompt_template, parameters, is_active, is_favorite, tags) VALUES
    -- Tareas de automatización del sistema
    (
        'Conectar WiFi + VPN',
        'Conecta automáticamente a la red WiFi y levanta la VPN',
        (SELECT id FROM tools WHERE name = 'wifi-vpn'),
        NULL,
        NULL,
        '{}',
        1,
        1,
        '["network", "automation", "vpn", "startup"]'
    ),
    (
        'Levantar Stack Dev',
        'Inicia OmniRoute y OpenCode para desarrollo',
        (SELECT id FROM tools WHERE name = 'start-dev'),
        NULL,
        NULL,
        '{}',
        1,
        1,
        '["dev", "automation", "docker", "startup"]'
    ),
    (
        'Limpiar Sistema',
        'Elimina archivos temporales y libera espacio en disco',
        (SELECT id FROM tools WHERE name = 'cleanup'),
        NULL,
        NULL,
        '{}',
        1,
        0,
        '["maintenance", "disk", "cleanup"]'
    ),
    (
        'Backup Rápido',
        'Copia archivos importantes a directorio de backup',
        (SELECT id FROM tools WHERE name = 'backup'),
        NULL,
        NULL,
        '{}',
        1,
        0,
        '["backup", "safety", "maintenance"]'
    ),
    (
        'Monitor de Recursos',
        'Muestra uso de CPU, RAM y Disco en tiempo real',
        (SELECT id FROM tools WHERE name = 'monitor'),
        NULL,
        NULL,
        '{}',
        1,
        0,
        '["monitoring", "system", "diagnostics"]'
    ),
    
    -- Tareas con IA
    (
        'Revisar Correos',
        'Lee correos nuevos y los clasifica con IA por urgencia e importancia',
        (SELECT id FROM tools WHERE name = 'email-check'),
        (SELECT id FROM providers WHERE name = 'OmniRoute Local'),
        'Analiza estos correos y clasifícalos por urgencia e importancia.

Para cada correo indica:
1. Nivel de urgencia: Alta / Media / Baja
2. Importancia: Alta / Media / Baja
3. Resumen en una línea
4. Acción recomendada (si aplica)

Correos:
{{.Emails}}

Responde en formato JSON con esta estructura:
{
  "clasificacion": [
    {
      "de": "remitente",
      "asunto": "asunto del correo",
      "urgencia": "alta/media/baja",
      "importancia": "alta/media/baja",
      "resumen": "resumen breve",
      "accion": "acción recomendada o null"
    }
  ]
}',
        '{"max_emails": 10}',
        1,
        1,
        '["email", "productivity", "ai", "classification"]'
    ),
    (
        'Resumir URL',
        'Extrae contenido de una URL y lo resume en ideas clave',
        (SELECT id FROM tools WHERE name = 'email-check'),  -- Placeholder, crear tool específico
        (SELECT id FROM providers WHERE name = 'OmniRoute Local'),
        'Resume el siguiente contenido web en 3-5 ideas clave.
Sé conciso, directo y enfócate en lo más importante.

URL: {{.URL}}
Contenido:
{{.Content}}

Responde en este formato:
## Resumen

**Ideas principales:**
1. [Idea 1]
2. [Idea 2]
3. [Idea 3]

**Conclusión:** [Una frase resumen]',
        '{}',
        1,
        0,
        '["ai", "productivity", "research", "summarization"]'
    ),
    (
        'Chat con IA',
        'Conversación libre con el proveedor de IA configurado',
        (SELECT id FROM tools WHERE name = 'email-check'),  -- Placeholder
        (SELECT id FROM providers WHERE name = 'OmniRoute Local'),
        '{{.UserInput}}',
        '{}',
        1,
        0,
        '["ai", "chat", "general", "assistant"]'
    ),
    (
        'Traductor de Comandos',
        'Convierte lenguaje natural a comandos de terminal',
        (SELECT id FROM tools WHERE name = 'email-check'),  -- Placeholder
        (SELECT id FROM providers WHERE name = 'OmniRoute Local'),
        'Convierte esta instrucción en lenguaje natural a un comando de terminal de Linux/macOS.

Instrucción: {{.Instruction}}

IMPORTANTE: Responde SOLO con el comando, sin explicaciones adicionales.
Si la instrucción no es clara o es peligrosa, responde con un error descriptivo.',
        '{}',
        1,
        0,
        '["ai", "terminal", "productivity", "translation"]'
    ),
    
    -- Tareas de comandos nativos
    (
        'Estado de Docker',
        'Muestra contenedores Docker en ejecución',
        (SELECT id FROM tools WHERE name = 'docker-ps'),
        NULL,
        NULL,
        '{}',
        1,
        0,
        '["docker", "dev", "monitoring"]'
    );

-- ============================================
-- Seed: Default Settings
-- ============================================
INSERT OR IGNORE INTO settings (key, value, category, is_encrypted) VALUES
    -- General
    ('app.name', 'Brain CLI', 'general', 0),
    ('app.version', '2.0.0', 'general', 0),
    ('app.log_level', 'info', 'general', 0),
    ('app.debug', 'false', 'general', 0),
    
    -- UI
    ('ui.theme', 'dark', 'ui', 0),
    ('ui.animations', 'true', 'ui', 0),
    ('ui.show_tips', 'true', 'ui', 0),
    ('ui.confirm_dangerous_actions', 'true', 'ui', 0),
    
    -- AI
    ('ai.default_provider_id', '1', 'ai', 0),
    ('ai.default_temperature', '0.7', 'ai', 0),
    ('ai.max_tokens', '2000', 'ai', 0),
    ('ai.stream_responses', 'true', 'ai', 0),
    
    -- System
    ('system.auto_update', 'true', 'system', 0),
    ('system.backup_enabled', 'true', 'system', 0),
    ('system.backup_interval_days', '7', 'system', 0),
    ('system.max_log_size_mb', '100', 'system', 0),
    
    -- Execution
    ('execution.max_concurrent', '1', 'system', 0),
    ('execution.default_timeout', '300', 'system', 0),
    ('execution.history_limit', '1000', 'system', 0),
    ('execution.auto_cleanup_days', '30', 'system', 0),
    
    -- Database
    ('db.auto_vacuum', 'true', 'system', 0),
    ('db.wal_mode', 'true', 'system', 0),
    
    -- Security
    ('security.encryption_enabled', 'true', 'system', 0),
    ('security.require_confirmation', 'true', 'system', 0);

-- ============================================
-- Seed: Default Tags
-- ============================================
INSERT OR IGNORE INTO tags (name, color) VALUES
    ('automation', '#3b82f6'),    -- Azul
    ('ai', '#8b5cf6'),            -- Púrpura
    ('dev', '#10b981'),           -- Verde
    ('network', '#f59e0b'),       -- Naranja
    ('maintenance', '#ef4444'),   -- Rojo
    ('productivity', '#06b6d4'),  -- Cyan
    ('startup', '#ec4899'),       -- Rosa
    ('monitoring', '#6366f1'),    -- Índigo
    ('backup', '#84cc16'),        -- Lima
    ('email', '#f97316');         -- Naranja oscuro

-- ============================================
-- Vincular tags con tasks (task_tags)
-- ============================================

-- Tags para "Conectar WiFi + VPN"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Conectar WiFi + VPN'),
    id 
FROM tags 
WHERE name IN ('network', 'automation', 'startup');

-- Tags para "Levantar Stack Dev"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Levantar Stack Dev'),
    id 
FROM tags 
WHERE name IN ('dev', 'automation', 'startup');

-- Tags para "Limpiar Sistema"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Limpiar Sistema'),
    id 
FROM tags 
WHERE name IN ('maintenance');

-- Tags para "Backup Rápido"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Backup Rápido'),
    id 
FROM tags 
WHERE name IN ('backup', 'maintenance');

-- Tags para "Monitor de Recursos"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Monitor de Recursos'),
    id 
FROM tags 
WHERE name IN ('monitoring');

-- Tags para "Revisar Correos"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Revisar Correos'),
    id 
FROM tags 
WHERE name IN ('email', 'productivity', 'ai');

-- Tags para "Resumir URL"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Resumir URL'),
    id 
FROM tags 
WHERE name IN ('ai', 'productivity');

-- Tags para "Chat con IA"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Chat con IA'),
    id 
FROM tags 
WHERE name IN ('ai');

-- Tags para "Traductor de Comandos"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Traductor de Comandos'),
    id 
FROM tags 
WHERE name IN ('ai', 'productivity');

-- Tags para "Estado de Docker"
INSERT OR IGNORE INTO task_tags (task_id, tag_id) 
SELECT 
    (SELECT id FROM tasks WHERE name = 'Estado de Docker'),
    id 
FROM tags 
WHERE name IN ('dev', 'monitoring');

-- ============================================
-- End of Migration 002
-- ============================================
