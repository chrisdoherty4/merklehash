package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

type hashable interface {
	Hash() []byte
}

// fileHasher is an abstract path hasher.
type fileHasher struct {
	path     string
	fileInfo os.FileInfo
	hasher   hash.Hash
}

func newPathHasher(path string, hasher hash.Hash) *fileHasher {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	pathFileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return &fileHasher{path: path, fileInfo: pathFileInfo, hasher: hasher}
}

func (this *fileHasher) Hash() []byte {
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

// MerkleTree is an interface to retrieve directory hashes.
// This structure forms the basis for the merkle hash.
type MerkleTree struct {
	fileHasher
	nodes []hashable
}

// NewMerkleTree creates and initialises a MerkleTree structure.
func NewMerkleTree(path string, hasher hash.Hash) *MerkleTree {
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

	directory := MerkleTree{
		fileHasher: *newPathHasher(path, hasher),
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			directory.Add(NewMerkleTree(fullPath, hasher))
		} else {
			directory.Add(newPathHasher(fullPath, hasher))
		}
	}

	return &directory
}

// Add adds a node to a MerkleTree.
func (this *MerkleTree) Add(node hashable) {
	this.nodes = append(this.nodes, node)
}

// Hash retrieves the merkle hash for a given directory.
func (this *MerkleTree) Hash() []byte {
	for _, node := range this.nodes {
		this.hasher.Write(node.Hash())
	}

	return this.hasher.Sum(nil)
}

// Define all supported hashing algorithms.
var algorithms = map[string]func() hash.Hash{
	"md5":      func() hash.Hash { return md5.New() },
	"sha256":   func() hash.Hash { return sha256.New() },
	"sha224":   func() hash.Hash { return sha256.New224() },
	"sha384":   func() hash.Hash { return sha512.New384() },
	"sha512":   func() hash.Hash { return sha512.New() },
	"sha3-224": func() hash.Hash { return sha3.New224() },
	"sha3-256": func() hash.Hash { return sha3.New256() },
	"sha3-384": func() hash.Hash { return sha3.New384() },
	"sha3-512": func() hash.Hash { return sha3.New512() },
}

var alg = flag.String(
	"alg",
	"sha256",
	"Hashing algorithm to use with the merkle tree.",
)
var raw = flag.Bool("raw", false, "Print only the hex hash.")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Invalid invocation")
		os.Exit(0)
	}

	hasher, ok := algorithms[*alg]
	if !ok {
		fmt.Println("Available algs...")
		os.Exit(1)
	}

	// Create a new merkle tree and output the hex representation of it's hash.
	fmt.Print(fmt.Sprintf("%x", NewMerkleTree(flag.Arg(0), hasher()).Hash()))

	if !*raw {
		fmt.Print(fmt.Sprintf(" %v", flag.Arg(0)))
	}

	println()
}
