.PHONY: bench
bench:
	go test -bench . -run skip ./...

.PHONY: test
test:
	go test -race -v ./...

.PHONY: run
run: test bench
