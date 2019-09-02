# merklehash
A go implementation of a [merkle tree](https://en.wikipedia.org/wiki/Merkle_tree) for hashing arbitrary sized directories.

## Usage
Install
```
go install "github.com/chrisdoherty4/merklehash"
```
and run
```
merklehash <directory>
```
## Known issues
* Symlinks are not followed and the app errors instead.
