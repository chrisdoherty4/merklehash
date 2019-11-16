# merklehash

A go implementation of a [merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
for hashing arbitrary sized directories.

Caution should be exercised when hashing directories with large files or large quantities
of files as it could take some time.

## Install

**Requires**

* Golang 1.12

```bash
go get "github.com/chrisdoherty4/merklehash/cmd/merklehash"
go install "github.com/chrisdoherty4/merklehash/cmd/merklehash"

merklehash <directory>
```

## API

The `merklehash` package exposes a list of supported algorithms and their
string identifiers and a method fo creating a new `MerkleTree` structure. You
can utilize the merklehash package in your own go code.

### Example

```golang
package main

import "github.com/chrisdoherty4/merklehash/pkg/merklehash"

func main() {
  path := "/directory/to/hash"
  merkleHash := merklehash.New(path, merkletree.GetHasher('sha256'))

  fmt.Fprintf(os.Stdout, "%x\n", merkleHash.Hash())
}
```

## To do

* Add support for specifying multiple directories.
* Overridable protection against huge file systems.
* Resolve symlinks when traversing directory structures.
* Improve/expand help output.
* Package for Linux platforms.
* Package in a Windows installer.
* Complete test code.

## Known issues

* Symlinks are not followed and the app errors instead.
