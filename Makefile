IMAGE_NAME = hidori/genprop:latest

.PHONY: lint
lint:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint:latest-alpine golangci-lint run

.PHONY: format
format:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint:latest-alpine golangci-lint run --fix

.PHONY: test
test:
	go test ./...
	go run ./cmd/genprop/main.go -- ./example/example.go > ./example/example.prop.go
	go run ./cmd/example/main.go

.PHONY: build
build: test lint
	go build -o ./bin/genprop ./cmd/genprop/main.go

.PHONY: run
run: build
	./bin/genprop ./example/example.go > ./example/example.prop.go

.PHONY: mod/download
mod/download:
	go mod downloadmake

.PHONY: mod/tidy
mod/tidy:
	go mod tidy

.PHONY: mod/update
mod/update:
	go get -u ./...

.PHONY: container/build
container/build:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} .

.PHONY: container/rebuild
container/rebuild:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} --no-cache .

.PHONY: container/rmi
container/rmi:
	docker rmi -f ${IMAGE_NAME}

.PHONY: version/patch
version/patch: test lint
	git fetch
	git checkout main
	git pull
	docker run --rm hidori/semver -i patch `cat ./version.txt` > ./version.txt
	git add ./version.txt
	git commit -m 'Updated version.txt'
	git push
	git tag v`cat ./version.txt`
	git push origin --tags
