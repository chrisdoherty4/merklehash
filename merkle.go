package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const sha256Len = 64

type Hasher interface {
	Hash() []byte
}

type Node struct {
	Hasher
	path string
}

type FileHash struct {
	Node
}

func NewFileHash(path string) FileHash {
	return FileHash{Node: Node{path: path}}
}

func (this FileHash) Hash() []byte {
	fd, err := os.Open(this.path)
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

type DirectoryHash struct {
	Node
	nodes []Node
}

func NewDirectoryHash(path string) DirectoryHash {
	directory := DirectoryHash{Node: Node{path: path}}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	directory.nodes = make([]Node, len(files))

	for _, file := range files {
		if file.IsDir() {
			directory.Add(NewDirectoryHash(file.Name()).Node)
		} else {
			directory.Add(NewFileHash(file.Name()).Node)
		}
	}

	return directory
}

func (this DirectoryHash) Add(node Node) {
	this.nodes = append(this.nodes, node)
}

func (this DirectoryHash) Hash() []byte {
	hashes := make([]byte, len(this.nodes)*sha256Len)

	fmt.Print("Hashes slice len: ")
	fmt.Println(len(hashes))

	for i, _ := range this.nodes {
		println(i * sha256Len)
	}

	return sha256.New().Sum(hashes)
}

func main() {
	f := NewDirectoryHash("./test")
	fmt.Println(fmt.Sprintf("%x", f.Hash()))
}
