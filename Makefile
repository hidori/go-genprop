AUTHOR = hidori
PROJECT = genprop
IMAGE_NAME = ${AUTHOR}/${PROJECT}:latest

.PHONY: lint
lint:
	docker run --rm -v ${PWD}:${PWD} -w ${PWD} golangci/golangci-lint:latest-alpine golangci-lint run

.PHONY: format
format:
	docker run --rm -v ${PWD}:${PWD} -w ${PWD} golangci/golangci-lint:latest-alpine golangci-lint run --fix

.PHONY: test
test:
	go test -cover ./generator
	go run ./cmd/genprop/main.go -- ./example/example.go > ./example/example_prop.go
	go run ./cmd/example/main.go

.PHONY: build
build:
	go build -o ./bin/genprop ./cmd/genprop/main.go

.PHONY: run
run: build
	./bin/genprop -- ./example/example.go > ./example/example_prop.go
	go run ./cmd/example/main.go

.PHONY: container/clean
container/rmi:
	docker rmi -f ${IMAGE_NAME}

.PHONY: container/build
container/build:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} .

.PHONY: container/rebuild
container/rebuild:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} --no-cache .

.PHONY: container/run
container/run: container/build
	docker run --rm -it -v ${PWD}:${PWD} -w ${PWD} ${IMAGE_NAME} ./example/example.go > ./example/example_prop.go
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
