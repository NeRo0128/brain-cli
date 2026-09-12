-- ============================================
-- brain-cli Seed Data v2.0
-- Migration: 002_seed_data.sql
-- ============================================

-- ============================================
-- Seed: Providers
-- ============================================
INSERT OR IGNORE INTO providers
    (name, type, base_url, model, description, api_key_env, is_active, config, max_tokens, timeout_seconds)
VALUES
    ('OmniRoute Local', 'omniroute', 'http://localhost:20128', 'auto',
     'Gateway multi-modelo local',
     'OMNIROUTE_API_KEY', 1,
     '{"temperature": 0.7}', 4000, 60),

    ('Ollama Local', 'ollama', 'http://localhost:11434', 'llama3.2',
     'Modelos locales vía Ollama',
     '', 0,
     '{"stream": true}', 4000, 120),

    ('DeepSeek Cloud', 'deepseek', 'https://api.deepseek.com', 'deepseek-chat',
     'API de DeepSeek (compatible con OpenAI)',
     'DEEPSEEK_API_KEY', 0,
     '{}', 4000, 90);

-- ============================================
-- Seed: Tools (scripts embebidos)
-- ============================================
INSERT OR IGNORE INTO tools
    (name, description, script_type, category, script_content, command, requires_sudo, timeout_seconds, is_builtin, version)
VALUES
    ('wifi-vpn', 'Conectar WiFi y levantar VPN',
     'bash', 'network',
     '#!/bin/bash
set -e
echo "🔌 Conectando a WiFi..."
# TODO: nmcli device wifi connect "$SSID" password "$PASSWORD"
echo "🔐 Levantando VPN..."
# TODO: openvpn --config "$OVPN_CONFIG" --daemon
echo "✅ Conexión establecida"',
     '', 1, 120, 1, 1),

    ('start-dev', 'Levantar stack de desarrollo',
     'bash', 'dev',
     '#!/bin/bash
set -e
command -v docker >/dev/null || { echo "❌ Docker no instalado"; exit 1; }
echo "📦 Iniciando OmniRoute..."
# docker-compose -f ~/omniroute/docker-compose.yml up -d
echo "💻 Iniciando OpenCode..."
# code ~/projects
echo "✅ Stack listo"',
     '', 0, 180, 1, 1),

    ('cleanup', 'Limpieza de temporales y logs',
     'bash', 'maintenance',
     '#!/bin/bash
set -e
echo "🧹 Limpieza..."
sudo rm -rf /tmp/* 2>/dev/null || true
rm -rf ~/.cache/* 2>/dev/null || true
find /var/log -type f -name "*.log" -mtime +7 -delete 2>/dev/null || true
command -v docker >/dev/null && docker system prune -f
echo "✅ Limpieza completada"',
     '', 1, 300, 1, 1),

    ('backup', 'Backup rápido de configuraciones',
     'bash', 'maintenance',
     '#!/bin/bash
set -e
DEST="$HOME/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$DEST"
cp -r ~/.config "$DEST/config" 2>/dev/null || true
cp ~/.bashrc "$DEST/" 2>/dev/null || true
cp ~/.zshrc  "$DEST/" 2>/dev/null || true
echo "✅ Backup en: $DEST"',
     '', 0, 600, 1, 1),

    ('monitor', 'Monitor de CPU/RAM/Disco',
     'bash', 'system',
     '#!/bin/bash
echo "📊 Recursos del sistema"
echo "======================="
echo "🔥 CPU:"
top -bn1 | grep "Cpu(s)" | awk ''{print "  " $2}'' || true
echo "💾 RAM:"
free -h | awk ''NR==2{print "  " $3 " / " $2}''
echo "💿 Disco:"
df -h / | awk ''NR==2{print "  " $3 " / " $2 " (" $5 ")"}''',
     '', 0, 30, 1, 1),

    ('docker-ps', 'Listar contenedores Docker',
     'native', 'dev', '',
     'docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"',
     0, 10, 1, 1);

-- ============================================
-- Seed: Tags
-- ============================================
INSERT OR IGNORE INTO tags (name, color) VALUES
    ('automation',   '#3b82f6'),
    ('ai',           '#8b5cf6'),
    ('dev',          '#10b981'),
    ('network',      '#f59e0b'),
    ('maintenance',  '#ef4444'),
    ('productivity', '#06b6d4'),
    ('startup',      '#ec4899'),
    ('monitoring',   '#6366f1'),
    ('backup',       '#84cc16'),
    ('email',        '#f97316');

-- ============================================
-- Seed: Tasks
-- ============================================
INSERT OR IGNORE INTO tasks
    (id, name, description, type, tool_id, requires_ai, ai_prompt, params, priority, is_active, is_favorite)
VALUES
    ('conectar-wifi-vpn', 'Conectar WiFi + VPN',
     'Conecta WiFi y levanta VPN automáticamente',
     'script',
     (SELECT id FROM tools WHERE name = 'wifi-vpn'),
     0, '', '{}', 'high', 1, 1),

    ('levantar-stack-dev', 'Levantar Stack Dev',
     'Inicia OmniRoute y OpenCode',
     'script',
     (SELECT id FROM tools WHERE name = 'start-dev'),
     0, '', '{}', 'high', 1, 1),

    ('limpiar-sistema', 'Limpiar Sistema',
     'Elimina temporales y libera espacio',
     'script',
     (SELECT id FROM tools WHERE name = 'cleanup'),
     0, '', '{}', 'medium', 1, 0),

    ('backup-rapido', 'Backup Rápido',
     'Copia configuraciones importantes',
     'script',
     (SELECT id FROM tools WHERE name = 'backup'),
     0, '', '{}', 'medium', 1, 0),

    ('monitor-recursos', 'Monitor de Recursos',
     'Muestra uso de CPU, RAM y Disco',
     'script',
     (SELECT id FROM tools WHERE name = 'monitor'),
     0, '', '{}', 'low', 1, 0),

    ('estado-docker', 'Estado de Docker',
     'Lista contenedores Docker',
     'command',
     (SELECT id FROM tools WHERE name = 'docker-ps'),
     0, '', '{}', 'low', 1, 0),

    ('chat-ia', 'Chat con IA',
     'Conversación libre con el proveedor activo',
     'ai', NULL, 1,
     '{{.UserInput}}',
     '{}', 'low', 1, 0),

    ('traducir-comando', 'Traductor de Comandos',
     'Convierte lenguaje natural a comando de terminal',
     'ai', NULL, 1,
     'Convierte esta instrucción a un comando de terminal Linux/macOS.

Instrucción: {{.Instruction}}

IMPORTANTE: Responde SOLO con el comando, sin explicaciones.',
     '{}', 'medium', 1, 0);

-- ============================================
-- Seed: task_tags (pivote)
-- ============================================
INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'conectar-wifi-vpn', id FROM tags WHERE name IN ('network', 'automation', 'startup');

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'levantar-stack-dev', id FROM tags WHERE name IN ('dev', 'automation', 'startup');

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'limpiar-sistema', id FROM tags WHERE name = 'maintenance';

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'backup-rapido', id FROM tags WHERE name IN ('backup', 'maintenance');

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'monitor-recursos', id FROM tags WHERE name = 'monitoring';

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'estado-docker', id FROM tags WHERE name IN ('dev', 'monitoring');

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'chat-ia', id FROM tags WHERE name = 'ai';

INSERT OR IGNORE INTO task_tags (task_id, tag_id)
SELECT 'traducir-comando', id FROM tags WHERE name IN ('ai', 'productivity');
