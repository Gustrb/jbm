all: build

build: clean outdir
	go build -o ./bin/jbm ./src/

outdir:
	mkdir -p bin

clean:
	rm -rf ./bin

test:
	go test ./...
