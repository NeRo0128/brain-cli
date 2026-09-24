package main

import (
	"context"
	"fmt"
	"os"
	"time"

	authadapter "github.com/NeRo0128/brain-cli/internal/adapters/auth"
	"github.com/NeRo0128/brain-cli/internal/adapters/config"
	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/adapters/executor"
	"github.com/NeRo0128/brain-cli/internal/core/paths"
	"github.com/NeRo0128/brain-cli/internal/core/settings"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui"
	"github.com/NeRo0128/brain-cli/internal/ui/icons"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	taskScreens "github.com/NeRo0128/brain-cli/internal/ui/screens/tasks"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	"github.com/NeRo0128/brain-cli/internal/ui/theme"
	authuc "github.com/NeRo0128/brain-cli/internal/usecases/auth"
	settingsuc "github.com/NeRo0128/brain-cli/internal/usecases/settings"
	taskEsxec "github.com/NeRo0128/brain-cli/internal/usecases/task"
	tooluc "github.com/NeRo0128/brain-cli/internal/usecases/tool"
	"github.com/NeRo0128/brain-cli/pkg/utils"

	tea "charm.land/bubbletea/v2"
)

// Version will be set during build via ldflags
// Inyectadas en build time vía -ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// configPath se puede sobreescribir con BRAIN_CONFIG.
const defaultConfigPath = "configs/config.yaml"

func resolveConfigPath() string {
	if p := os.Getenv("BRAIN_CONFIG"); p != "" {
		return p
	}
	if paths.DevMode() {
		return "configs/config.yaml"
	}
	return paths.ConfigFile()
}

func main() {

	// * Paths
	if err := paths.EnsureDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creando directorios XDG: %v\n", err)
		os.Exit(1)
	}

	// * Config
	configPath := resolveConfigPath()

	if !paths.DevMode() && os.Getenv("BRAIN_DB_PATH") == "" {
		_ = os.Setenv("BRAIN_DB_PATH", paths.DBFile())
	}

	defaultDB := "data/brain.db"
	if !paths.DevMode() {
		defaultDB = paths.DBFile()
	}

	if err := config.EnsureConfig(configPath, defaultDB); err != nil {
		fmt.Fprintf(os.Stderr, "Error creando config: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.NewYAMLLoader().Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cargando config: %v\n", err)
		os.Exit(1)
	}

	// * Logger
	log, err := utils.NewLogger(utils.LoggerOptions{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		MaxSizeMB:  cfg.Logging.MaxSizeMB,
		OutputPath: cfg.Logging.OutputPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configurando logger: %v\n", err)
		os.Exit(1)
	}
	log.Debug().
		Str("version", Version).
		Str("commit", Commit).
		Str("build_date", BuildDate).
		Str("config", configPath).
		Str("db_path", cfg.Database.Path).
		Bool("dev_mode", paths.DevMode()).
		Msg("Brain CLI - debug")
	// * DB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.New(ctx, database.Options{
		Path:        cfg.Database.Path,
		AutoMigrate: cfg.Database.AutoMigrate,
	})
	if err != nil {
		log.Err(err).Msg("abriendo database")
		os.Exit(1)
	}

	defer db.Close()

	// * Repositories

	taskRepo := repositories.NewTaskRepository(db.DB())
	toolRepo := repositories.NewToolRepository(db.DB())
	execRepo := repositories.NewExecutionRepository(db.DB())

	// * Use cases
	executorUC := taskEsxec.NewExecutor(taskRepo, toolRepo, execRepo, executor.New)
	managerUC := taskEsxec.NewManager(taskRepo, toolRepo)
	toolManagerUC := tooluc.NewManager(toolRepo)
	interpreters := tool.Available(tool.Detect())
	if len(interpreters) == 0 {
		log.Warn().Msg("no hay intérpretes disponibles; el tipo 'script' estará deshabilitado")
	}

	// * Settings
	settingsRepo := repositories.NewSettingsRepository(db.DB())

	schema := settings.NewSchema([]settings.Field{
		{Key: "ui.theme", Category: "ui",
			Values: theme.Names()},
		{Key: "ui.icons", Category: "ui",
			Values: icons.Names()},
		{Key: "ui.brand_style", Category: "ui",
			Values: []string{"minimal", "slim", "big"}},
		{Key: "logging.level", Category: "logging",
			Values: []string{"debug", "info", "warn", "error"}},
		{Key: "logging.format", Category: "logging",
			Values: []string{"pretty", "json"}},
	})

	settingsMgr := settingsuc.NewManager(settingsRepo, schema, cfg, log)

	// Aplicar overrides sobre la config base (no muta cfg).
	cfgResolved, err := settingsMgr.Resolve(ctx)
	if err != nil {
		log.Error().Err(err).Msg("aplicando overrides; usando config base")
	} else {
		cfg = cfgResolved
	}

	// * Auth
	var oauthClient authuc.OAuthClient
	if cfg.GitHub.ClientID != "" {
		oauthClient = authadapter.NewGitHubOAuth(cfg.GitHub.ClientID)
		log.Debug().Msg("GitHub OAuth configurado")
	} else {
		log.Warn().Msg("GitHub ClientID no configurado; auth deshabilitada")
	}

	tokenRepo := authadapter.NewKeyringRepository()
	authManager := authuc.NewManager(tokenRepo, oauthClient, log)

	// * UI

	// Styles
	appStyles := styles.New(
		theme.Get(cfg.UI.Theme),
		true,
		icons.Get(cfg.UI.Icons),
	)

	log.Info().Msg("Brain CLI iniciando")
	// Deps
	deps := ui.Deps{
		Version:         Version,
		Cfg:             cfg,
		ExecUC:          executorUC,
		TaskRepo:        taskRepo,
		ToolRepo:        toolRepo,
		AuthManager:     authManager,
		ExecRepo:        execRepo,
		SettingsManager: settingsMgr,
		Keys:            keys.New(nil),
		Log:             log,
		Manager:         managerUC,
		ToolManager:     toolManagerUC,
		Interpreter:     interpreters,
		Styles:          &appStyles,
	}

	log.Debug().
		Str("theme", cfg.UI.Theme).
		Str("icons", cfg.UI.Icons).
		Msg("apariencia cargada")

	mainScreen := taskScreens.NewListScreen(taskRepo, &appStyles)
	p := tea.NewProgram(ui.NewModels(deps, mainScreen))

	if _, err := p.Run(); err != nil {
		log.Fatal().Err(err).Msg("Fallo la TUI")
	}

	log.Info().Msg("Brain CLI finalizado")
}
