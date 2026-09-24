# Brain CLI — Makefile

APP_NAME    := brain-cli
CMD_PATH    := ./cmd/brain-cli
BUILD_DIR   := build
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE        ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS     := -s -w \
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.BuildDate=$(DATE)

GOFLAGS     := -trimpath

# ─── Build local ─────────────────────────────────────────────
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)
	@echo "✓ $(BUILD_DIR)/$(APP_NAME) ($(VERSION))"

# ─── Run ─────────────────────────────────────────────────────
.PHONY: run
run:
	go run $(CMD_PATH)

# ─── Tests ───────────────────────────────────────────────────
.PHONY: test
test:
	go test ./... -cover

.PHONY: test-verbose
test-verbose:
	go test ./... -v -cover

# ─── Lint ────────────────────────────────────────────────────
.PHONY: vet
vet:
	go vet ./...

# ─── Cross-compile (todos los targets) ───────────────────────
PLATFORMS := \
	linux/amd64 linux/arm64 \
	darwin/amd64 darwin/arm64 \
	windows/amd64 windows/arm64

.PHONY: release
release:
	@mkdir -p $(BUILD_DIR)/release
	@for p in $(PLATFORMS); do \
		GOOS=$${p%/*}; GOARCH=$${p#*/}; \
		EXT=""; [ "$$GOOS" = "windows" ] && EXT=".exe"; \
		OUT="$(BUILD_DIR)/release/$(APP_NAME)-$$GOOS-$$GOARCH$$EXT"; \
		echo "→ $$OUT"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build $(GOFLAGS) \
			-ldflags="$(LDFLAGS)" -o "$$OUT" $(CMD_PATH) || exit 1; \
	done
	@echo "✓ cross-compile completo"

.PHONY: checksums
checksums: release
	@cd $(BUILD_DIR)/release && \
		sha256sum * > SHA256SUMS 2>/dev/null || shasum -a 256 * > SHA256SUMS
	@echo "✓ checksums generados"

# ─── Limpieza ────────────────────────────────────────────────
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# ─── Ayuda ───────────────────────────────────────────────────
.PHONY: help
help:
	@echo "Targets disponibles:"
	@echo "  build          → compila el binario local"
	@echo "  run            → ejecuta la app"
	@echo "  test           → corre tests"
	@echo "  vet            → corre go vet"
	@echo "  release        → cross-compile para 6 plataformas"
	@echo "  checksums      → release + SHA256SUMS"
	@echo "  clean          → borra build/"
	@echo ""
	@echo "Version actual: $(VERSION)"
