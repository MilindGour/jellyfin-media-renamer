BINARY_NAME=jmr
.DEFAULT_GOAL := build

clean:
	@echo
	@echo "[MAKEFILE] Cleaning output directories"
	rm -rf ./dist
	rm -rf ./frontend/.svelte-kit
	rm -rf ./frontend/build

build: clean
	@echo
	@echo "[MAKEFILE] Building all"
	$(MAKE) buildbe
	$(MAKE) buildfe
	mv ./frontend/build ./dist/frontend

buildbe:
	@echo
	@echo "[MAKEFILE] Building backend"
	go build -o ./dist/backend/${BINARY_NAME}
	cp service/jmr.service ./dist/backend/
	cp service/start-all.sh ./dist/
	cp config.json ./dist/backend/config.sample.json

buildfe:
	@echo
	@echo "[MAKEFILE] Building frontend"
	cd frontend && npm install && npm run build

test:
	@echo
	@echo "[MAKEFILE] Testing"
	go test ./...

deploy:
	@echo
	@echo "[MAKEFILE] Deploying"
	/var/lib/jenkins/scripts/deploy-jmr.sh dist

watch-dev:
	@echo
	@echo "[MAKEFILE] Starting Watch"
	go run . --dev
