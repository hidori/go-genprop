.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint golangci-lint run

.PHONY: format
format:
	docker run --rm -v $$PWD:$$PWD -w $$PWD golangci/golangci-lint golangci-lint run --fix

.PHONY: run
run:
	go run ./cmd/genprop/main.go -- ./example/example.go > ./example/example.prop.go
	go run ./cmd/genprop/main.go -- ./example/example2.go > ./example/example2.prop.go

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
