all:
	go build -o merklehash cmd/merklehash/*.go

clean:
	rm ./merklehash