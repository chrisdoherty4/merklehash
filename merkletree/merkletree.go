package merkletree

import (
	"context"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// HashFactory creates a new instance of hash.Hash
type HashFactory func() hash.Hash

// digestResult stores a hashing result suitable for communicating over a channel.
type digestResult struct {
	// Index is used to track file ordering so we can produce deterministic digests.
	index int
	// Data is the hash data in bytes
	digest []byte
	// Err holds an error should something go wrong with hashing.
	err error
}

type digestFunc func(p string) ([]byte, error)

// New hashes a directory by iterating over all files and nested files combining
// their individual digests into 1. Symlinks are not followed and there is no
// protection against a recursive symlink.
func New(
	ctx context.Context,
	dirpath string,
	factory HashFactory,
) ([]byte, error) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %v %v", dirpath, err)
	}

	var fileCount int
	results := make(chan digestResult)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		fn digestFunc
		p  string
	)

	for _, file := range files {
		p = path.Join(dirpath, file.Name())

		if file.IsDir() {
			fn = func(path string) ([]byte, error) { return New(ctx, path, factory) }
		} else {
			fn = func(path string) ([]byte, error) { return calculateFileDigest(path, factory) }
		}

		go func(index int, path string, fn digestFunc) {
			hash, err := fn(path)

			result := digestResult{
				index:  index,
				digest: hash,
				err:    err,
			}

			select {
			case <-ctx.Done():
				return
			case results <- result:
			}
		}(fileCount, p, fn)
		fileCount++
	}

	digests := make([][]byte, fileCount)

	// Collect hashes storing them in a slice to maintain ordering so that when we produce the
	// final hash it's deterministic.
	err = func() error {
		for ; fileCount != 0; fileCount-- {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case result := <-results:
				if result.err != nil {
					return err
				}

				digests[result.index] = result.digest
			}
		}

		return nil
	}()

	if err != nil {
		return nil, err
	}

	hasher := factory()
	for _, d := range digests {
		hasher.Write(d)
	}

	return hasher.Sum(nil), nil
}

func calculateFileDigest(filepath string, factory HashFactory) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v %v", filepath, err)
	}
	defer file.Close()

	hasher := factory()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, fmt.Errorf("could not create digest: %v", err)
	}

	return hasher.Sum(nil), nil
}
