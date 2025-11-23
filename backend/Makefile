BINARY_NAME=jellyfin-media-renamer
.DEFAULT_GOAL := build

clean:
	@echo
	@echo "[MAKEFILE] Cleaning output directories"
	rm -rf ./dist

build: clean
	@echo
	@echo "[MAKEFILE] Building backend"
	go build -o ./dist/${BINARY_NAME}

test:
	@echo
	@echo "[MAKEFILE] Testing"
	go test ./...

watch:
	@echo
	@echo "[MAKEFILE] Starting Watch"
	go run . --dev
