.PHONY: clean

all: dist
	go build -o dist/program ./...

dist: clean
	mkdir dist

clean: 
	rm -rf dist

test:
	go test -coverprofile cover.out -v ./...
