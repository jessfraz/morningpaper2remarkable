#!/usr/bin/env make

.DEFAULT_GOAL := help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RED    := $(shell tput -Txterm setaf 1)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM := 20

### Help
## Show help
help:
	@echo ''
	@echo '${WHITE}Usage:${RESET}'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@awk '/(^[a-zA-Z\-\.\_0-9]+:)|(^###[a-zA-Z ]+)/ { \
		header = match($$1, /^###(.*)/); \
		if (header) { \
			title = substr($$0, 5, length($$0)); \
			printf "${CYAN}%s${RESET}\n", title; \
		} \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
