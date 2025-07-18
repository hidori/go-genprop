AUTHOR = hidori
PROJECT = genprop
IMAGE_NAME = $(AUTHOR)/$(PROJECT):latest

DOCKER_LINT_CMD = docker run --rm -v $(PWD):$(PWD) -w $(PWD) golangci/golangci-lint:latest-alpine

.PHONY: lint
lint:
	$(DOCKER_LINT_CMD) golangci-lint run

.PHONY: format
format:
	$(DOCKER_LINT_CMD) golangci-lint run --fix

.PHONY: test
test:
	go test -v -cover -race ./generator
	go run ./cmd/genprop/main.go -- ./example/example.go > ./example/example_prop.go
	go run ./cmd/example/main.go

.PHONY: build
build:
	mkdir -p ./bin
	go build -o ./bin/genprop ./cmd/genprop/main.go

.PHONY: run
run: build
	./bin/genprop -- ./example/example.go > ./example/example_prop.go
	go run ./cmd/example/main.go

.PHONY: clean
clean:
	rm -rf ./bin/
	rm -f ./example/example_prop.go

.PHONY: container/rmi
container/rmi:
	docker rmi -f $(IMAGE_NAME)

.PHONY: dev
dev: clean format test build run

.PHONY: ci
ci: lint test build

.PHONY: container/build
container/build:
	docker build -f ./Dockerfile -t $(IMAGE_NAME) .

.PHONY: container/rebuild
container/rebuild:
	docker build -f ./Dockerfile -t $(IMAGE_NAME) --no-cache .

.PHONY: container/run
container/run: container/build
	docker run --rm -it -v $(PWD):$(PWD) -w $(PWD) $(IMAGE_NAME) ./example/example.go > ./example/example_prop.go
	go run ./cmd/example/main.go

.PHONY: version/patch
version/patch: test lint
	git fetch
	git checkout main
	git pull
	docker run --rm hidori/semver -i patch `cat ./meta/version.txt` > ./meta/version.txt
	git add ./meta/version.txt
	git commit -m 'Updated version.txt'
	git push
	git tag v`cat ./meta/version.txt`
	git push origin --tags
