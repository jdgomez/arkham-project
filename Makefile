GO_VERSION := $(shell sed -n "s/^go //p" backend/go.mod)
GOLANGCI_LINT_VERSION := v2.15.0
GO_LINT := cd backend && GOTOOLCHAIN=go$(GO_VERSION) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: test test-backend lint gen build dev e2e

# Unit and HTTP tests, then the E2E tests against the built binary.
test: test-backend e2e

test-backend:
	cd backend && go test ./...

e2e: build
	cd frontend && npx playwright test

lint:
	$(GO_LINT) run
	cd frontend && npm run lint

# Writes frontend/src/generated from the Go structs (D12).
gen:
	cd backend && go tool tygo generate

# One binary that embeds the built front end.
build: frontend/node_modules
	cd frontend && npm run build
	cd backend && go build -o ../bin/arkham ./cmd/arkham

# Go server on :8080 and Vite dev server proxying /api to it.
dev: frontend/node_modules
	cd backend && go run ./cmd/arkham & server=$$!; \
	trap "kill $$server" EXIT; \
	cd frontend && npm run dev

frontend/node_modules: frontend/package-lock.json
	cd frontend && npm ci
	touch frontend/node_modules
