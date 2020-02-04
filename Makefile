.PHONY: lint
lint:
	echo TODO

.PHONY: test
test:
	go test ./...

.PHONY: test-all
test-all:
	go test ./...
	docker run -v $(PWD):/work -w /work --rm golang:1.13 go test ./...
