all: build

build: clean outdir
	go build -o ./bin/jbm ./src/

outdir:
	mkdir -p bin

clean:
	rm -rf ./bin

test: build-test
	go test ./... -v

build-test:
	javac ./tests/fixtures/*.java
