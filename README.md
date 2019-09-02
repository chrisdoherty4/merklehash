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
