package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
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

func NewFileHash(path string) FileHash {
	return FileHash{Node: Node{path: path}}
}

func (this FileHash) Hash() (hash []byte) {
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

func (this DirectoryHash) Hash() (hash []byte) {
	return []byte{}
}

func main() {
	f := NewDirectoryHash("./test")
	fmt.Println(fmt.Sprintf("%x", f.Hash()))
}
