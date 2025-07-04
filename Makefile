BINARY_NAME=jmr
.DEFAULT_GOAL := build

clean:
	rm -rf ./dist

build:
	go build -o ./dist/${BINARY_NAME}
	cp -r service ./dist/
	cp config.json ./dist/config.sample.json

test:
	go test ./...

deploy:
	/var/lib/jenkins/scripts/deploy-jmr.sh dist
