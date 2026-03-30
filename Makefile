# Build klse cli tools
.PHONY: build
build:
	go build -o bin/klse.exe cmd/klse/main.go

# Run klse cli tools
.PHONY: run
run: build
	./bin/klse.exe