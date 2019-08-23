package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

type Hasher interface {
	Hash() (hash []byte)
}

type Node struct {
	Hasher
	path string
}

type FileHash struct {
	Node
}

func (f *FileHash) Hash() (hash []byte) {
	fd, err := os.Open(f.path)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	h := sha256.New()

	if _, err := io.Copy(h, fd); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

func main() {
	f := FileHash{"./test/test1.txt"}
	fmt.Println(fmt.Sprintf("%x", f.Hash()))
}
