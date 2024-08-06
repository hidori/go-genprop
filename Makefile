IMAGE_NAME = hidori/genprop

.PHONY: test
test:
	go test ./...
	go run ./cmd/example/main.go

.PHONY: lint
lint:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint golangci-lint run

.PHONY: format
format:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint golangci-lint run --fix

.PHONY: run
run:
	go run ./cmd/genprop/main.go -- ./example/example.go > ./example/example.prop.go

.PHONY: build
build:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} .

.PHONY: rebuild
rebuild:
	docker build -f ./Dockerfile -t ${IMAGE_NAME} --no-cache .

.PHONY: rmi
rmi:
	docker rmi -f ${IMAGE_NAME}

.PHONY: mod/download
mod/download:
	go mod download

.PHONY: mod/tidy
mod/tidy:
	go mod tidy

.PHONY: mod/update
mod/update:
	go get -u ./...

.PHONY: version/patch
version/patch: test
	git fetch
	git checkout main
	git pull
	docker run --rm hidori/semver -i patch `cat ./version.txt` > ./version.txt
	git add ./version.txt
	git commit -m 'Updated version.txt'
	git push
	git tag v`cat ./version.txt`
	git push origin --tags
