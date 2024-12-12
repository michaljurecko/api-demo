.PHONY: build

dc-build:
	cd .devcontainer; docker compose build

clean:
	./scripts/clean.sh

shell:
	./scripts/shell.sh

run-raw:
	go run ./cmd/demo || true

run:
	./scripts/run-pretty-log.sh ./cmd/demo || true

lint:
	bash ./scripts/lint.sh

fix:
	bash ./scripts/fix.sh

tests:
	bash ./scripts/tests.sh ./internal/pkg

deps-tidy:
	go mod tidy && go mod vendor

deps-upgrade:
	go-mod-upgrade && go mod tidy && go mod vendor

gen-model:
	go generate ./internal/pkg/app/demo/model/...

gen-api:
	cd api; buf generate

gen-wire:
	cd internal/pkg/app/demo/cmd; wire

buf-lint:
	cd api; buf lint

buf-update:
	cd api; buf dep update
