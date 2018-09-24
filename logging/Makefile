.DEFAULT_GOAL = research

.PHONY: bench
bench:
	go test -bench . -run skip ./...

.PHONY: test
test:
	go test -race -v ./...

.PHONY: research
research: test bench
