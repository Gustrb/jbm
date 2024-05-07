all: build

build: outdir
	go build -o ./bin/jbm ./src/

outdir:
	mkdir -p bin

test:
	go test ./...
