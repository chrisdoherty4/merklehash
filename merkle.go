package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const sha256Len = 64

// Hasher is an interface to retrieve arbitrary hashes of objects.
type Hasher interface {
	Hash() []byte
}

// NodeHasher is an abstract path hasher.
type NodeHasher struct {
	Hasher
	path string
}

// FileHasher is an interface to retrieve the hash of a file.
type FileHasher struct {
	NodeHasher
}

// NewFileHasher creates and initiliases a new FileHasher.
func NewFileHasher(path string) FileHasher {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	return FileHasher{NodeHasher: NodeHasher{path: path}}
}

// Hash retrieves the hash of a File.
func (this FileHasher) Hash() []byte {
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

// DirectoryHasher is an interface to retrieve directory hashes.
// This structure forms the basis for the merkle hash.
type DirectoryHasher struct {
	NodeHasher
	nodes []NodeHasher
}

// NewDirectoryHasher creates and initialises a DirectoryHasher structure.
func NewDirectoryHasher(path string) DirectoryHasher {
	// Make the path absolute and ensure it exists. We check it's existence
	// by ensuring there's no error when running ioutil.ReadDir.
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	directory := DirectoryHasher{NodeHasher: NodeHasher{path: path}}
	directory.nodes = make([]NodeHasher, len(files))

	for _, file := range files {
		if file.IsDir() {
			directory.Add(NewDirectoryHasher(file.Name()).NodeHasher)
		} else {
			directory.Add(NewFileHasher(file.Name()).NodeHasher)
		}
	}

	return directory
}

// Add adds a node to a DirectoryHasher.
func (this DirectoryHasher) Add(node NodeHasher) {
	this.nodes = append(this.nodes, node)
}

// Hash retrieves the merkle hash for a given directory.
func (this DirectoryHasher) Hash() []byte {
	hashes := make([]byte, len(this.nodes)*sha256Len)

	fmt.Println(fmt.Sprintf("Hash len: %v", len(hashes)))

	for _, node := range this.nodes {
		fmt.Println(fmt.Sprintf("%v", node.path))
	}

	return sha256.New().Sum(hashes)
}

func main() {
	f := NewDirectoryHasher("./test")
	_ = f.Hash()
}
