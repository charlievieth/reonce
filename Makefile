.PHONY: quick
quick: # runs all tests without coverage or the race detector
	go test ./...

.PHONY: test
test: # runs all tests against the package with race detection and coverage percentage
	go test -race -cover ./...

.PHONY: reoncetest # test the `reoncetest` build tag
reoncetest:
	go test -race -cover -tags reonce ./...

.PHONY: lint
lint:
	./scripts/lint

.PHONY: all
all: test reoncetest lint

.PHONY: cover
cover: # runs all tests against the package, generating a coverage report and opening it in the default browser
	@mkdir -p ./build && \
	go test -race -covermode=atomic -coverprofile=./build/cover.out ./... && \
	go tool cover -html ./build/cover.out -o ./build/cover.html && \
	which open && open ./build/cover.html

.PHONY: clean
clean:
	rm -rf ./build
