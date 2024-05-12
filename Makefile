build:
	@go build -o bin/out

run: build
	@./bin/out
