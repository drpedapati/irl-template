# IRL Template Project Generator
# Usage:
#   make                    - Show welcome message
#   make irl project-name   - Create a new IRL project
#   make irl project-name -t template-name - Use specific template

TEMPLATE_DIR := 01-plans/templates
DEFAULT_TEMPLATE := irl-basic-template

# Colors for beautiful output
BLUE := \033[0;34m
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
MAGENTA := \033[0;35m
BOLD := \033[1m
RESET := \033[0m

# Default target - show welcome message
.PHONY: default
default:
	@echo ""
	@echo "${BOLD}${CYAN}╔══════════════════════════════════════════════════════════════╗${RESET}"
	@echo "${BOLD}${CYAN}║${RESET}                                                              ${BOLD}${CYAN}║${RESET}"
	@echo "${BOLD}${CYAN}║${RESET}  ${BOLD}${GREEN}✨ IRL Template Project Generator ✨${RESET}                      ${BOLD}${CYAN}║${RESET}"
	@echo "${BOLD}${CYAN}║${RESET}                                                              ${BOLD}${CYAN}║${RESET}"
	@echo "${BOLD}${CYAN}╚══════════════════════════════════════════════════════════════╝${RESET}"
	@echo ""
	@echo "${BOLD}${BLUE}Welcome to the Idempotent Research Loop (IRL) Template!${RESET}"
	@echo ""
	@echo "${CYAN}Create reproducible, auditable research workflows with AI assistants.${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Quick Start:${RESET}"
	@echo ""
	@echo "  ${GREEN}make irl${RESET}                        ${CYAN}→ Interactive setup (recommended)${RESET}"
	@echo "  ${GREEN}make irl my-project${RESET}              ${CYAN}→ Create project with default template${RESET}"
	@echo "  ${GREEN}make irl my-project -t meeting-abstract${RESET}"
	@echo "                                      ${CYAN}→ Use specific template${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Available Templates:${RESET}"
	@echo ""
	@echo "  ${MAGENTA}•${RESET} ${BOLD}irl-basic-template${RESET}      ${CYAN}General purpose research workflow${RESET}"
	@echo "  ${MAGENTA}•${RESET} ${BOLD}scientific-abstract${RESET}    ${CYAN}For journal article abstracts${RESET}"
	@echo "  ${MAGENTA}•${RESET} ${BOLD}meeting-abstract${RESET}        ${CYAN}For conference/meeting abstracts${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}More Information:${RESET}"
	@echo ""
	@echo "  ${GREEN}make help${RESET}                  ${CYAN}→ Show detailed usage guide${RESET}"
	@echo "  ${GREEN}make templates${RESET}             ${CYAN}→ List all available templates${RESET}"
	@echo ""
	@echo "${CYAN}📚 Documentation: https://github.com/drpedapati/irl-template/wiki${RESET}"
	@echo ""
	@echo "${BOLD}${GREEN}Ready to start? Run: ${BOLD}make irl${RESET} ${CYAN}(interactive)${RESET}"
	@echo ""

# Prevent make from treating arguments as targets (must come before irl target)
.PHONY: -t --t meeting-abstract scientific-abstract irl-basic-template
-t --t meeting-abstract scientific-abstract irl-basic-template:
	@true
.PHONY: %
%:
	@true

# Main irl target - interactive or parse arguments
.PHONY: irl
irl:
	@PROJECT_NAME=""; \
	TEMPLATE=""; \
	TEMPLATE_FLAG=0; \
	ARGS="$(filter-out irl,$(MAKECMDGOALS))"; \
	for arg in $$ARGS; do \
		if [ "$$arg" = "-t" ] || [ "$$arg" = "--t" ]; then \
			TEMPLATE_FLAG=1; \
		elif [ "$$TEMPLATE_FLAG" = "1" ]; then \
			TEMPLATE="$$arg"; \
			TEMPLATE_FLAG=0; \
		elif [ -z "$$PROJECT_NAME" ]; then \
			PROJECT_NAME="$$arg"; \
		fi; \
	done; \
	if [ -z "$$PROJECT_NAME" ]; then \
		echo ""; \
		echo "${BOLD}${CYAN}Interactive IRL Project Setup${RESET}"; \
		echo ""; \
		printf "${CYAN}Project name:${RESET} "; \
		read PROJECT_NAME; \
		if [ -z "$$PROJECT_NAME" ]; then \
			echo "${BOLD}${YELLOW}⚠${RESET} Project name cannot be empty"; \
			exit 1; \
		fi; \
		echo ""; \
		echo "${BOLD}${YELLOW}Available Templates:${RESET}"; \
		echo ""; \
		echo "  ${MAGENTA}0${RESET} ${BOLD}No template${RESET} (start with empty main-plan.md)"; \
		TEMPLATE_NUM=1; \
		for template in $$(ls -1 $(TEMPLATE_DIR)/*.md 2>/dev/null | sed 's|.*/||' | sed 's|\.md$$||' | sort); do \
			TEMPLATE_DESC=""; \
			case "$$template" in \
				irl-basic-template) TEMPLATE_DESC="General purpose IRL template" ;; \
				scientific-abstract) TEMPLATE_DESC="For journal article abstracts" ;; \
				meeting-abstract) TEMPLATE_DESC="For conference/meeting abstracts" ;; \
				*) TEMPLATE_DESC="Template" ;; \
			esac; \
			echo "  ${MAGENTA}$$TEMPLATE_NUM${RESET} ${BOLD}$$template${RESET}      ${CYAN}$$TEMPLATE_DESC${RESET}"; \
			TEMPLATE_NUM=$$((TEMPLATE_NUM + 1)); \
		done; \
		echo ""; \
		printf "${CYAN}Select template [0-$$((TEMPLATE_NUM - 1))]:${RESET} "; \
		read SELECTION; \
		if [ "$$SELECTION" = "0" ] || [ -z "$$SELECTION" ]; then \
			TEMPLATE=""; \
		else \
			TEMPLATE_LIST=$$(ls -1 $(TEMPLATE_DIR)/*.md 2>/dev/null | sed 's|.*/||' | sed 's|\.md$$||' | sort); \
			TEMPLATE=$$(echo "$$TEMPLATE_LIST" | sed -n "$$SELECTION p"); \
			if [ -z "$$TEMPLATE" ]; then \
				echo "${BOLD}${YELLOW}⚠${RESET} Invalid selection, using no template"; \
				TEMPLATE=""; \
			fi; \
		fi; \
	fi; \
	if [ -d "$$PROJECT_NAME" ]; then \
		echo "${BOLD}${YELLOW}⚠${RESET} ${BOLD}Error:${RESET} Directory '${CYAN}$$PROJECT_NAME${RESET}' already exists"; \
		exit 1; \
	fi; \
	echo ""; \
	echo "${BOLD}${CYAN}Creating IRL project:${RESET} ${GREEN}$$PROJECT_NAME${RESET}"; \
	mkdir -p "$$PROJECT_NAME"; \
	echo "${CYAN}Copying template files...${RESET}"; \
	for file in $$(ls -A | grep -v "^$$PROJECT_NAME$$" | grep -v "^\.git$$"); do \
		if [ -d "$$file" ] && [ -d "$$file/.git" ]; then \
			cp -r "$$file" "$$PROJECT_NAME"/ 2>/dev/null && \
			rm -rf "$$PROJECT_NAME/$$file/.git" 2>/dev/null || true; \
		else \
			cp -r "$$file" "$$PROJECT_NAME"/ 2>/dev/null || true; \
		fi; \
	done; \
	cd "$$PROJECT_NAME" && \
		rm -rf .git && \
		git init -q && \
		git add -A && \
		git commit -q -m "Initial commit from IRL template"; \
	if [ -n "$$TEMPLATE" ]; then \
		echo "${CYAN}Setting up template:${RESET} ${GREEN}$$TEMPLATE${RESET}"; \
		if [ -f "$$PROJECT_NAME/$(TEMPLATE_DIR)/$$TEMPLATE.md" ]; then \
			cp "$$PROJECT_NAME/$(TEMPLATE_DIR)/$$TEMPLATE.md" "$$PROJECT_NAME/main-plan.md"; \
			echo "${BOLD}${GREEN}✓${RESET} Template applied"; \
		else \
			echo "${BOLD}${YELLOW}⚠${RESET} Warning: Template '$$TEMPLATE' not found, starting with empty plan"; \
		fi; \
	else \
		echo "${CYAN}Starting with empty plan document${RESET}"; \
	fi; \
	echo "${BOLD}${GREEN}✓${RESET} Project initialized"; \
	echo ""; \
	echo "${BOLD}${GREEN}✓${RESET} ${BOLD}Created IRL project:${RESET} ${CYAN}$$PROJECT_NAME${RESET}"; \
	if [ -n "$$TEMPLATE" ]; then \
		echo "${BOLD}${GREEN}✓${RESET} ${BOLD}Template used:${RESET} ${CYAN}$$TEMPLATE${RESET}"; \
	else \
		echo "${BOLD}${GREEN}✓${RESET} ${BOLD}Template:${RESET} ${CYAN}None (empty plan)${RESET}"; \
	fi; \
	echo ""; \
	echo "${BOLD}${YELLOW}Next steps:${RESET}"; \
	echo ""; \
	echo "  ${CYAN}cd${RESET} ${GREEN}$$PROJECT_NAME${RESET}"; \
	echo "  ${CYAN}# Edit main-plan.md to customize your plan${RESET}"; \
	echo "  ${CYAN}# Start your first iteration!${RESET}"; \
	echo ""

# Help target
.PHONY: help
help:
	@echo ""
	@echo "${BOLD}${CYAN}IRL Template Project Generator${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Usage:${RESET}"
	@echo ""
	@echo "  ${GREEN}make${RESET}                        ${CYAN}Show welcome message${RESET}"
	@echo "  ${GREEN}make irl${RESET}                    ${CYAN}Interactive setup (recommended)${RESET}"
	@echo "  ${GREEN}make irl project-name${RESET}       ${CYAN}Create project with default template${RESET}"
	@echo "  ${GREEN}make irl project-name -t scientific-abstract${RESET}"
	@echo "                                      ${CYAN}Create project with specific template${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Available Templates:${RESET}"
	@echo ""
	@echo "  ${MAGENTA}•${RESET} ${BOLD}irl-basic-template${RESET}      ${CYAN}General purpose IRL template${RESET}"
	@echo "  ${MAGENTA}•${RESET} ${BOLD}scientific-abstract${RESET}    ${CYAN}For writing scientific journal abstracts${RESET}"
	@echo "  ${MAGENTA}•${RESET} ${BOLD}meeting-abstract${RESET}        ${CYAN}For writing conference/meeting abstracts${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Examples:${RESET}"
	@echo ""
	@echo "  ${GREEN}make irl my-research${RESET}"
	@echo "  ${GREEN}make irl conference-abstracts${RESET}"
	@echo "  ${GREEN}make irl apa-2025 -t meeting-abstract${RESET}"
	@echo ""

# List available templates
.PHONY: templates
templates:
	@echo ""
	@echo "${BOLD}${CYAN}Available Templates:${RESET}"
	@echo ""
	@if [ -d "$(TEMPLATE_DIR)" ]; then \
		for template in $$(ls -1 $(TEMPLATE_DIR)/*.md 2>/dev/null | sed 's|.*/||' | sed 's|\.md$$||'); do \
			echo "  ${MAGENTA}•${RESET} ${BOLD}$$template${RESET}"; \
		done; \
	else \
		echo "  ${BOLD}${YELLOW}⚠${RESET} Error: Could not find template directory"; \
	fi
	@echo ""

# ==============================================================================
# CLI Development Commands
# ==============================================================================

LDFLAGS = -ldflags "-X github.com/drpedapati/irl-template/cmd.Version=$(V)"

# Build local binary (dev version)
.PHONY: build
build:
	@go build -o irl .
	@echo "${GREEN}✓${RESET} Built ./irl (dev)"

# Run quick tests
.PHONY: test
test: build
	@./irl --help > /dev/null && echo "${GREEN}✓${RESET} --help"
	@./irl doctor > /dev/null && echo "${GREEN}✓${RESET} doctor"
	@./irl config > /dev/null && echo "${GREEN}✓${RESET} config"

# Clean build artifacts
.PHONY: clean
clean:
	@rm -f irl irl-template irl-darwin-* irl-linux-*
	@echo "${GREEN}✓${RESET} Cleaned"

# Build all platform binaries with version
.PHONY: build-all
build-all: clean
	@if [ -z "$(V)" ]; then echo "Usage: make build-all V=x.y.z"; exit 1; fi
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o irl-darwin-arm64 .
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o irl-darwin-amd64 .
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o irl-linux-amd64 .
	@echo "${GREEN}✓${RESET} Built v$(V): darwin-arm64, darwin-amd64, linux-amd64"

# Create GitHub release: make release V=0.3.3 M="Release notes"
.PHONY: release
release: build-all
	@gh release create v$(V) irl-darwin-arm64 irl-darwin-amd64 irl-linux-amd64 --title "v$(V)" --notes "$(M)"
	@echo "${GREEN}✓${RESET} Released v$(V)"
	@echo "  Run: ${CYAN}make brew-update V=$(V)${RESET}"

# Update homebrew formula: make brew-update V=0.3.3
.PHONY: brew-update
brew-update:
	@if [ -z "$(V)" ]; then echo "Usage: make brew-update V=x.y.z"; exit 1; fi
	@./scripts/brew-update.sh $(V)

# Reinstall from tap (forces fresh fetch)
.PHONY: brew-reinstall
brew-reinstall:
	@brew uninstall irl 2>/dev/null || true
	@brew untap drpedapati/tap 2>/dev/null || true
	@brew tap drpedapati/tap
	@brew install drpedapati/tap/irl
	@echo "${GREEN}✓${RESET} Reinstalled"

# Full release: make full-release V=0.3.3 M="Release notes"
.PHONY: full-release
full-release: release brew-update brew-reinstall clean
	@echo "${GREEN}✓${RESET} ${BOLD}Released v$(V) and updated homebrew${RESET}"

# Show version info
.PHONY: version
version:
	@echo "Latest: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'none')"
	@echo "Installed: $$(irl version 2>/dev/null || echo 'not installed')"

# Dev help
.PHONY: dev-help
dev-help:
	@echo ""
	@echo "${BOLD}${CYAN}CLI Development Commands${RESET}"
	@echo ""
	@echo "  ${GREEN}make build${RESET}              Build local binary"
	@echo "  ${GREEN}make test${RESET}               Build and run tests"
	@echo "  ${GREEN}make clean${RESET}              Remove build artifacts"
	@echo "  ${GREEN}make version${RESET}            Show version info"
	@echo ""
	@echo "${BOLD}${YELLOW}Release Workflow${RESET}"
	@echo ""
	@echo "  ${GREEN}make release V=x.y.z M=\"msg\"${RESET}      Create GitHub release"
	@echo "  ${GREEN}make brew-update V=x.y.z${RESET}          Update homebrew formula"
	@echo "  ${GREEN}make brew-reinstall${RESET}               Reinstall from tap"
	@echo "  ${GREEN}make full-release V=x.y.z M=\"msg\"${RESET} All of the above"
	@echo ""
