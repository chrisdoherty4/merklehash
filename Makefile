.PHONY: all
all: test build

.PHONY: build
build:
	go build -o merklehash ./cmd/merklehash

.PHONY: test
test:
	go test -v ./merkletree

.PHONY: clean
clean:
	rm ./merklehash
