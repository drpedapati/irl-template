# IRL Template Project Generator
# Usage:
#   make                    - Show this welcome message
#   make IRL                - Create a new IRL project with default settings
#   make IRL NAME=my-project - Create a project with custom name
#   make IRL TEMPLATE=scientific-abstract NAME=my-abstract - Use specific template

TEMPLATE_DIR := 01-plans/templates
DEFAULT_TEMPLATE := irl-basic-template
DEFAULT_NAME := irl-project

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
	@echo "  ${GREEN}make IRL${RESET}                    ${CYAN}→ Create project 'irl-project' with default template${RESET}"
	@echo "  ${GREEN}make IRL NAME=my-project${RESET}    ${CYAN}→ Create project with custom name${RESET}"
	@echo "  ${GREEN}make IRL TEMPLATE=meeting-abstract NAME=conference${RESET}"
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
	@echo "${BOLD}${GREEN}Ready to start? Run: ${BOLD}make IRL${RESET}"
	@echo ""

# Main IRL target
.PHONY: IRL
IRL: $(or $(NAME),$(DEFAULT_NAME))
	@echo ""
	@echo "${BOLD}${GREEN}✓${RESET} ${BOLD}Created IRL project:${RESET} ${CYAN}$(or $(NAME),$(DEFAULT_NAME))${RESET}"
	@echo "${BOLD}${GREEN}✓${RESET} ${BOLD}Template used:${RESET} ${CYAN}$(or $(TEMPLATE),$(DEFAULT_TEMPLATE))${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Next steps:${RESET}"
	@echo ""
	@echo "  ${CYAN}cd${RESET} ${GREEN}$(or $(NAME),$(DEFAULT_NAME))${RESET}"
	@echo "  ${CYAN}# Edit 01-plans/main-plan.md to customize your plan${RESET}"
	@echo "  ${CYAN}# Start your first iteration!${RESET}"
	@echo ""

# Create new IRL project
$(DEFAULT_NAME): $(or $(NAME),$(DEFAULT_NAME))
	@true

# Pattern rule for named projects
%:
	@if [ -d "$@" ]; then \
		echo "${BOLD}${YELLOW}⚠${RESET} ${BOLD}Error:${RESET} Directory '${CYAN}$@${RESET}' already exists"; \
		exit 1; \
	fi
	@echo ""
	@echo "${BOLD}${CYAN}Creating IRL project:${RESET} ${GREEN}$@${RESET}"
	@mkdir -p "$@"
	@echo "${CYAN}Copying template files...${RESET}"
	@cp -r . "$@"/ 2>/dev/null || \
		(echo "${BOLD}${YELLOW}⚠${RESET} ${BOLD}Error:${RESET} Could not copy template files" && exit 1)
	@cd "$@" && \
		rm -rf .git && \
		git init -q && \
		git add -A && \
		git commit -q -m "Initial commit from IRL template"
	@echo "${CYAN}Setting up template:${RESET} ${GREEN}$(or $(TEMPLATE),$(DEFAULT_TEMPLATE))${RESET}"
	@if [ -f "$@/$(TEMPLATE_DIR)/$(or $(TEMPLATE),$(DEFAULT_TEMPLATE)).md" ]; then \
		cp "$@/$(TEMPLATE_DIR)/$(or $(TEMPLATE),$(DEFAULT_TEMPLATE)).md" "$@/01-plans/main-plan.md"; \
		echo "${BOLD}${GREEN}✓${RESET} Template applied"; \
	else \
		echo "${BOLD}${YELLOW}⚠${RESET} Warning: Template '$(or $(TEMPLATE),$(DEFAULT_TEMPLATE))' not found, using default"; \
		cp "$@/$(TEMPLATE_DIR)/$(DEFAULT_TEMPLATE).md" "$@/01-plans/main-plan.md" 2>/dev/null || true; \
	fi
	@echo "${BOLD}${GREEN}✓${RESET} Project initialized"

# Help target
.PHONY: help
help:
	@echo ""
	@echo "${BOLD}${CYAN}IRL Template Project Generator${RESET}"
	@echo ""
	@echo "${BOLD}${YELLOW}Usage:${RESET}"
	@echo ""
	@echo "  ${GREEN}make${RESET}                        ${CYAN}Show welcome message${RESET}"
	@echo "  ${GREEN}make IRL${RESET}                    ${CYAN}Create project with default name 'irl-project'${RESET}"
	@echo "  ${GREEN}make IRL NAME=my-project${RESET}    ${CYAN}Create project with custom name${RESET}"
	@echo "  ${GREEN}make IRL TEMPLATE=scientific-abstract NAME=my-abstract${RESET}"
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
	@echo "  ${GREEN}make IRL${RESET}"
	@echo "  ${GREEN}make IRL NAME=conference-abstracts${RESET}"
	@echo "  ${GREEN}make IRL TEMPLATE=meeting-abstract NAME=apa-2025${RESET}"
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
