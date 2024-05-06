all: build

build: outdir
	go build -o bin/jbm ./src/main.go

outdir:
	mkdir -p bin
