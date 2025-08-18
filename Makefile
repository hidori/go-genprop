AUTHOR = hidori
PROJECT = genprop
IMAGE_NAME = $(AUTHOR)/$(PROJECT):latest

DOCKER_LINT_CMD = docker run --rm -v $(PWD):$(PWD) -w $(PWD) golangci/golangci-lint:latest-alpine

.PHONY: lint
lint:
	$(DOCKER_LINT_CMD) golangci-lint config verify
	$(DOCKER_LINT_CMD) golangci-lint run ./internal/... ./public/...

.PHONY: format
format:
	find ./internal ./public -name "*.go" -exec go tool goimports -w {} +
	$(DOCKER_LINT_CMD) golangci-lint config verify
	$(DOCKER_LINT_CMD) golangci-lint run --fix ./internal/... ./public/...

.PHONY: test
test:
	go test -v -cover -race ./internal/... ./public/...

.PHONY: build
build:
	mkdir -p ./bin
	go build -o ./bin/genprop ./cmd/genprop/main.go

.PHONY: example/generate
example/generate:
	go run ./cmd/genprop/main.go -- ./example/basic/user.go > example/basic/user_prop.go
	go run ./cmd/genprop/main.go -- ./example/private-setter/user.go > example/private-setter/user_prop.go
	go run ./cmd/genprop/main.go -initialism="id,url,api,json,xml,sql,http,uuid" -- ./example/initialism/api_client.go > example/initialism/api_client_prop.go
	go run ./cmd/genprop/main.go -new-func -- ./example/new/user.go > example/new/user_prop.go
	go run ./cmd/genprop/main.go -new-func -validation-func="validateField" -validation-tag="validate" -- ./example/validate/account.go > example/validate/account_prop.go

.PHONY: example/run
example/run: example/generate
	go run ./cmd/example/basic/main.go
	go run ./cmd/example/private-setter/main.go
	go run ./cmd/example/initialism/main.go
	go run ./cmd/example/new/main.go
	go run ./cmd/example/validate/main.go

.PHONY: clean
clean:
	rm -rf ./bin/
	rm -rf ./tmp/

.PHONY: container/rmi
container/rmi:
	docker rmi -f $(IMAGE_NAME)

.PHONY: container/build
container/build:
	docker build -f ./Dockerfile -t $(IMAGE_NAME) .

.PHONY: container/rebuild
container/rebuild:
	docker build -f ./Dockerfile -t $(IMAGE_NAME) --no-cache .

.PHONY: container/run
container/run: container/build
	docker run --rm -it -v $(PWD):$(PWD) -w $(PWD) $(IMAGE_NAME) ./example/basic/user.go > ./example/basic/user_prop.go
	go run ./cmd/example/basic/main.go

.PHONY: version/patch
version/patch: test lint
	git fetch
	git checkout main
	git pull
	docker run --rm hidori/semver -i patch `cat ./public/meta/version.txt` > ./public/meta/version.txt
	git add ./public/meta/version.txt
	git commit -m 'Updated version.txt'
	git push
	git tag v`cat ./public/meta/version.txt`
	git push origin --tags
