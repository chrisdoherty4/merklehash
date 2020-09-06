# merklehash

A go implementation of a [merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
for hashing arbitrary sized directories.

Caution should be exercised when hashing directories with large files or large 
quantities of files as it could take some time. Should a platform support
recursive symlinks one must be careful to ensure they do not exist within
the directory being hashed as it'll result in an infinite loop.

## Install

**Requires**

* Golang 1.12

```bash
go get "github.com/chrisdoherty4/merklehash/cmd/merklehash"
go install "github.com/chrisdoherty4/merklehash/cmd/merklehash"

merklehash <directory>
```

## API

The `merkletree` package exposes a single function, `New()`. The function
accepts a `context.Context` that can be cancelled by the caller as desired.

### Example

```golang
package main

import "github.com/chrisdoherty4/merklehash/merkletree"

func main() {
  path := "/directory/to/hash"
  hash := merkletree.New(context.Background(), path, sha256.New)

  fmt.Printf("%x\n", hash)
}
```
