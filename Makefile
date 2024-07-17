.PHONY: test
test:
	go test ./...

.PHONY: mod/download
mod/download:
	go mod download

.PHONY: mod/tidy
mod/tidy:
	go mod tidy

.PHONY: tag/patch
tag/patch: test
	git fetch
	git checkout main
	git pull
	docker run --rm hidori/semver -i patch `cat ./version.txt` > ./version.txt
	git add ./version.txt
	git commit -m 'Updated version.txt'
	git push
	git tag v`cat ./version.txt`
	git push origin --tags
