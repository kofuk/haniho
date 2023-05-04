.PHONY: all
all: codegen
	CGO_ENABLED=0 go build -o haniho .

.PHONY: run
run: codegen
	go run . --server

.PHONY: codegen
codegen:
	go generate ./...

.PHONY: clean
clean:
	$(RM) haniho generator/data.go
