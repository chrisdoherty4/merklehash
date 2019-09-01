package main

import (
	"container/list"
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

type pathHasher struct {
	path   string
	info   os.FileInfo
	hasher hash.Hash
}

func newPathHasher(path string, hasher hash.Hash) *pathHasher {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return &pathHasher{path: path, info: info, hasher: hasher}
}

// Hash uses the pathHasher's hasher to generate the hash of a file.
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

// MerkleTree is an interface to retrieve directory hashes.
// This structure forms the basis for the merkle hash.
type MerkleTree struct {
	pathHasher
	nodes []hashable
}

// NewMerkleTree creates and initialises a MerkleTree structure.
// MerkleTree's has can be retrieved via the MerkleTree.Hash() interface.
func NewMerkleTree(path string, hasher hash.Hash) *MerkleTree {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatalf("%v is not a directory", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	directory := MerkleTree{
		pathHasher: pathHasher{
			path:   path,
			info:   info,
			hasher: hasher,
		},
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

// algorithm defines a selectable algorithm from the command line interface.
type Algorithm struct {
	Ident string
	New   func() hash.Hash
}

// We're going to create an AlgorithmList so we can append a function todo
// retrieve the hash.Hash.
type AlgorithmList struct {
	list.List
}

var algorithms = AlgorithmList{}

// Bind a hidden find function to the list structure.
func (this *AlgorithmList) GetHasher(ident string) hash.Hash {
	for e := this.Front(); e != nil; e = e.Next() {
		alg := e.Value.(*Algorithm)
		if alg.Ident == ident {
			return alg.New()
		}
	}

	return nil
}

func init() {
	algorithms.PushBack(&Algorithm{
		"md5",
		func() hash.Hash { return md5.New() },
	})
	algorithms.PushBack(&Algorithm{
		"sha224",
		func() hash.Hash { return sha256.New224() },
	})
	algorithms.PushBack(&Algorithm{
		"sha256",
		func() hash.Hash { return sha256.New() },
	})
	algorithms.PushBack(&Algorithm{
		"sha384",
		func() hash.Hash { return sha512.New384() },
	})
	algorithms.PushBack(&Algorithm{
		"sha512",
		func() hash.Hash { return sha512.New() },
	})
	algorithms.PushBack(&Algorithm{
		"sha3-224",
		func() hash.Hash { return sha3.New224() },
	})
	algorithms.PushBack(&Algorithm{
		"sha3-256",
		func() hash.Hash { return sha3.New256() },
	})
	algorithms.PushBack(&Algorithm{
		"sha3-384",
		func() hash.Hash { return sha3.New384() },
	})
	algorithms.PushBack(&Algorithm{
		"sha3-512",
		func() hash.Hash { return sha3.New512() },
	})
}

func main() {
	// Define the command line interface
	alg := flag.String(
		"alg",
		"sha256",
		"Hashing algorithm to use with the merkle tree.",
	)

	raw := flag.Bool("raw", false, "Print only the hex hash.")

	// TODO: Tidy up help comment
	flag.Parse()

	if flag.NArg() != 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	hasher := algorithms.GetHasher(*alg)

	if hasher == nil {
		fmt.Fprintf(os.Stdout, "'%v' is not a valid algorithm. Value algorithms are:\n", *alg)
		for e := algorithms.Front(); e != nil; e = e.Next() {
			fmt.Fprintf(os.Stdout, "  %v\n", e.Value.(*Algorithm).Ident)
		}
		os.Exit(0)
	}

	// TODO: Add support for multiple directories.
	// TODO: Protection against huge file systems?

	// Create a new merkle tree and output the hex representation of it's hash.
	fmt.Fprintf(os.Stdout, "%x", NewMerkleTree(flag.Arg(0), hasher).Hash())

	if !*raw {
		fmt.Fprintf(os.Stdout, " %v", flag.Arg(0))
	}
	println()
}
