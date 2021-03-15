.PHONY: help

help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


build_linux: ## Build server for linux machine
	GOARCH=amd64 GOOS=linux go build

migrate:  ## apply migration script (do not uses config.yaml)
	cd migrations && goose sqlite3 ../app.db up