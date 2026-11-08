DIST := backend/internal/web/dist

# Dev tools are pinned in their own module so the application's go.mod holds only what it uses.
TOOL := go tool -modfile=../tools/go.mod

.PHONY: test test-backend lint gen build dev e2e

# Unit and HTTP tests, then the E2E tests against the built binary.
test: test-backend e2e

test-backend:
	cd backend && go test ./...

# The only way to run the E2E tests: it builds first, so they never see a stale binary.
e2e: build
	cd frontend && npx playwright test

lint: frontend/node_modules
	cd backend && $(TOOL) golangci-lint run
	cd frontend && npm run lint

# Writes frontend/src/generated from the Go structs (D12).
gen:
	cd backend && $(TOOL) tygo generate

# One binary that embeds the built front end. .gitkeep stays so the embed compiles on a clean clone.
build: frontend/node_modules
	find $(DIST) -mindepth 1 ! -name .gitkeep -delete
	cd frontend && npm run build
	cd backend && go build -o ../bin/arkham ./cmd/arkham

# Go server on :8080 and Vite dev server proxying /api to it.
dev: frontend/node_modules
	cd backend && go build -o ../bin/arkham-dev ./cmd/arkham
	bin/arkham-dev & server=$$!; \
	trap 'kill $$server 2>/dev/null' INT TERM EXIT; \
	cd frontend && npm run dev

frontend/node_modules: frontend/package-lock.json
	cd frontend && npm ci
	touch frontend/node_modules
