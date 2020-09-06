package merkletree

import (
	"context"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// ErrCancelled is triggered when the calling context is closed.
var ErrCancelled = errors.New("Cancelled")

// HashFactory creates a new instance of hash.Hash
type HashFactory func() hash.Hash

// New hashes a directory by iterating over all files and nested files combining
// their individual digests into 1. Symlinks are not followed and there is no
// protection against a recursive symlink.
func New(
	ctx context.Context,
	dirPath string,
	factory HashFactory,
) ([]byte, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%v is not a directory", dirPath)
	}

	var fileCount int32
	results := make(chan hashResult)
	ctx, cancel := context.WithCancel(ctx)

	err = filepath.Walk(
		dirPath,
		func(filepath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			go func(index int32) {
				hash, err := hashFile(filepath, factory)

				result := hashResult{
					Index: index,
				}

				switch {
				case err != nil:
					result.Err = err
				default:
					result.Data = hash
				}

				select {
				case <-ctx.Done():
					return
				case results <- result:
				}
			}(fileCount)

			fileCount++

			return nil
		},
	)

	if err != nil {
		cancel()
		return nil, err
	}

	hashes := make([][]byte, fileCount)

	// Collect hashes storing them in a slice to maintain ordering. This
	// ensures that when we iterate over the collected hashes and add them to
	// our final hash we have determinism as the Goroutines launched in the
	// Walk() function above are scheduled in an non-deterministic order.
	err = func() error {
		for {
			if fileCount == 0 {
				return nil
			}

			select {
			case <-ctx.Done():
				return ErrCancelled
			default:
			}

			select {
			case <-ctx.Done():
				return ErrCancelled
			case result := <-results:
				if result.Err != nil {
					return err
				}

				hashes[result.Index] = result.Data

				fileCount--
			}
		}
	}()

	if err != nil {
		cancel()
		return nil, err
	}

	hasher := factory()
	for _, hash := range hashes {
		hasher.Write(hash)
	}

	return hasher.Sum(nil), nil
}

type hashResult struct {
	// Err holds an error should something go wrong with hashing.
	Err error

	// Index is used to track file ordering.
	Index int32

	// Data is the hash data in bytes
	Data []byte
}

func hashFile(filepath string, factory HashFactory) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := factory()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
