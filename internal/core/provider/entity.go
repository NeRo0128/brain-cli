package provider

import "time"

// Provider representa un proveedor de IA configurado en el sistema
type Provider struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Type      ProviderType `db:"type"`
	Endpoint  string       `db:"endpoint"`
	APIKey    string       `db:"api_key"` // Encriptado
	Model     string       `db:"model"`
	IsActive  bool         `db:"is_active"`
	Config    string       `db:"config"` // JSON
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt *time.Time   `db:"updated_at"`
}

// ProviderType define los tipos de proveedores soportados
type ProviderType string

const (
	ProviderTypeOmniRoute ProviderType = "omniroute"
	ProviderTypeOllama    ProviderType = "ollama"
	ProviderTypeOpenAI    ProviderType = "openai"
	ProviderTypeDeepSeek  ProviderType = "deepseek"
	ProviderTypeAnthropic ProviderType = "anthropic"
	ProviderTypeCustom    ProviderType = "custom"
)

// RequiresAPIKey indica si este tipo necesita autenticación.
// OmniRoute y Ollama suelen correr locales sin key.
// todo
