build:
	go build -o /tmp/go-fiber-template main.go

run: build
	/tmp/go-fiber-template

watch:
	reflex -s -r '\.go$$' make run

test:
	./tools/runTests.sh

test-verbose:
	./tools/runTests.sh -v

docs:
	./tools/generateDocs.sh

sqlc:
	./tools/generateSQLC.sh

.PHONY: build run watch test test-verbose docs sqlc
