module github.com/nero/brain-cli

go 1.23

require (
	// TUI Framework
	github.com/charmbracelet/bubbletea v0.27.1
	github.com/charmbracelet/bubbles v0.20.0
	github.com/charmbracelet/lipgloss v0.13.1
	github.com/charmbracelet/glamour v0.8.0
	github.com/charmbracelet/log v0.4.0

	// Database
	github.com/mattn/go-sqlite3 v1.14.24
	modernc.org/sqlite v1.33.1 // CGO-free alternative

	// HTTP Client
	github.com/go-resty/resty/v2 v2.15.3

	// Configuration
	github.com/spf13/viper v1.19.0

	// Utilities
	github.com/google/uuid v1.6.0
	
	// Testing
	github.com/stretchr/testify v1.9.0
)
