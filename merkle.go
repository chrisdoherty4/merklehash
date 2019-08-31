package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Hasher is an interface to retrieve arbitrary hashes of objects.
type Hasher interface {
	Hash() []byte
}

// PathHasher is an abstract path hasher.
type PathHasher struct {
	path     string
	fileInfo os.FileInfo
}

// NewFileHasher creates and initiliases a new FileHasher.
func NewPathHasher(path string) *PathHasher {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	pathFileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return &PathHasher{path: path, fileInfo: pathFileInfo}
}

// Hash retrieves the hash of a File.
func (this *PathHasher) Hash() []byte {
	fd, err := os.Open(this.path)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, fd); err != nil {
		log.Fatal(err)
	}

	return hasher.Sum(nil)
}

// DirectoryHasher is an interface to retrieve directory hashes.
// This structure forms the basis for the merkle hash.
type DirectoryHasher struct {
	PathHasher
	nodes []*PathHasher
}

// NewDirectoryHasher creates and initialises a DirectoryHasher structure.
func NewDirectoryHasher(path string) *DirectoryHasher {
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

	directory := DirectoryHasher{
		PathHasher: PathHasher{
			path:     path,
			fileInfo: dirFileInfo,
		},
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			directory.Add(&NewDirectoryHasher(fullPath).PathHasher)
		} else {
			directory.Add(NewPathHasher(fullPath))
		}
	}

	return &directory
}

// Add adds a node to a DirectoryHasher.
func (this *DirectoryHasher) Add(node *PathHasher) {
	this.nodes = append(this.nodes, node)
}

// Hash retrieves the merkle hash for a given directory.
func (this *DirectoryHasher) Hash() []byte {
	hasher := sha256.New()

	for _, node := range this.nodes {
		hasher.Write(node.Hash())
	}

	return hasher.Sum(nil)
}

func hashFile(path string) {
	f := NewPathHasher(path)
	fmt.Println(fmt.Sprintf("Hash of %v: %v", path, hex.EncodeToString(f.Hash())))
}

func hashDirectory(path string) {
	f := NewDirectoryHasher(path)
	fmt.Println(fmt.Sprintf("Hash of %v: %v", path, hex.EncodeToString(f.Hash())))
}

func main() {
	hashFile("test/test1.txt")
	hashDirectory("test")
}
