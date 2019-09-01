package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type hashable interface {
	Hash() []byte
}

// pathHasher is an abstract path hasher.
type pathHasher struct {
	path     string
	fileInfo os.FileInfo
	hasher   hash.Hash
}

func newPathHasher(path string, algorithm hash.Hash) *pathHasher {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	pathFileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return &pathHasher{path: path, fileInfo: pathFileInfo, hasher: algorithm}
}

func (this *pathHasher) Hash() []byte {
	fd, err := os.Open(this.path)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	if _, err := io.Copy(this.hasher, fd); err != nil {
		log.Fatal(err)
	}

	return this.hasher.Sum(nil)
}

// MerkleHasher is an interface to retrieve directory hashes.
// This structure forms the basis for the merkle hash.
type MerkleHasher struct {
	pathHasher
	nodes []hashable
}

// NewMerkleHasher creates and initialises a MerkleHasher structure.
func NewMerkleHasher(path string, algorithm hash.Hash) *MerkleHasher {
	// Make the path absolute and ensure it exists.
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	dirFileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if !dirFileInfo.IsDir() {
		log.Fatalf("%v is not a directory", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	directory := MerkleHasher{
		pathHasher: *newPathHasher(path, algorithm),
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			directory.Add(NewMerkleHasher(fullPath, algorithm))
		} else {
			directory.Add(newPathHasher(fullPath, algorithm))
		}
	}

	return &directory
}

// Add adds a node to a MerkleHasher.
func (this *MerkleHasher) Add(node hashable) {
	this.nodes = append(this.nodes, node)
}

// Hash retrieves the merkle hash for a given directory.
func (this *MerkleHasher) Hash() []byte {
	for _, node := range this.nodes {
		this.hasher.Write(node.Hash())
	}

	return this.hasher.Sum(nil)
}

func main() {
	path := "test"
	fmt.Println(fmt.Sprintf("%v %v",
		path,
		hex.EncodeToString(NewMerkleHasher(path, sha256.New()).Hash()),
	))
}
