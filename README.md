# merklehash
A go implementation of a [merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
for hashing arbitrary sized directories.

## CLI
Install
```
go install "github.com/chrisdoherty4/merklehash"
```
and run
```
merklehash <directory>
```

## API
The `merklehash` package exposes a list of supported algorithms and their
string identifiers and a method fo creating a new `MerkleTree` structure. You
can utilize the merklehash package in your own go code.

**Example**
```
package main

import "github.com/chrisdoherty4/merklehash/merklehash"

func main() {
  path := "/directory/to/hash"
  merkleTree := merklehash.New(path, merkletree.Algorithms.GetHasher('sha256'))

  fmt.Fprintf(os.Stdout, "%x\n", merkleTree.Hash())
}
```

## To do
* Add support for specifying multiple directories.
* Overridable protection against huge file systems.
* Resolve symlinks when traversing directory structures.
* Improve/expand help output.
* Package for Linux platforms.
* Package in a Windows installer.
* Add a CI.
* Complete test code.

## Known issues
* Symlinks are not followed and the app errors instead.
