.PHONY: all
all: test reoncetest lint

# runs all tests without coverage or the race detector
.PHONY: quick
quick:
	go test ./...

# runs all tests against the package with race detection and coverage
# percentage
.PHONY: test
test:
	go test -race -cover ./...

# test the `reoncetest` build tag
.PHONY: reoncetest
reoncetest:
	go test -race -cover -tags reonce ./...

.PHONY: clean
clean:
	@[ ! -d ./build ] || rm -r ./build
