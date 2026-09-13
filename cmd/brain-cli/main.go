package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NeRo0128/brain-cli/internal/adapters/config"
	"github.com/NeRo0128/brain-cli/internal/adapters/database"
	"github.com/NeRo0128/brain-cli/internal/adapters/database/repositories"
	"github.com/NeRo0128/brain-cli/internal/adapters/executor"
	"github.com/NeRo0128/brain-cli/internal/core/tool"
	"github.com/NeRo0128/brain-cli/internal/ui"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	taskEsxec "github.com/NeRo0128/brain-cli/internal/usecases/task"
	"github.com/NeRo0128/brain-cli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

// Version will be set during build via ldflags
var Version = "dev"

// configPath se puede sobreescribir con BRAIN_CONFIG.
const defaultConfigPath = "configs/config.yaml"

func main() {

	// * Config
	configPath := os.Getenv("BRAIN_CONFIG")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := config.NewYAMLLoader().Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cargando configuración: %v\n", err)
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
		Str("config", configPath).
		Str("db_path", cfg.Database.Path).
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

	// * Use case: la factory de executors se inyecta como función.
	executorUC := taskEsxec.NewExecutor(taskRepo, toolRepo, execRepo, executor.New)
	managerUC := taskEsxec.NewManager(taskRepo, toolRepo)
	interpreters := tool.Available(tool.Detect())
	if len(interpreters) == 0 {
		log.Warn().Msg("no hay intérpretes disponibles; el tipo 'script' estará deshabilitado")
	}

	//  * KeyMap
	keyRegistry := keys.New(nil)

	// * UI

	log.Info().
		Msg("Brain CLI iniciando")

	deps := ui.Deps{
		Version:     Version,
		Cfg:         cfg,
		ExecUC:      executorUC,
		TaskRepo:    taskRepo,
		ToolRepo:    toolRepo,
		ExecRepo:    execRepo,
		Keys:        keyRegistry,
		Log:         log,
		Manager:     managerUC,
		Interpreter: interpreters,
	}

	p := tea.NewProgram(
		ui.NewModels(
			deps,
			screens.NewMainScreen(taskRepo),
		),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal().Err(err).Msg("Fallo la TUI")
	}

	log.Info().Msg("Brain CLI finalizado")
}
