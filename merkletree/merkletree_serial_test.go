package merkletree

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
)

// NewSerial generates a directory digest in serial
func NewSerial(
	ctx context.Context,
	dirpath string,
	factory HashFactory,
) ([]byte, error) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %v %v", dirpath, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hasher := factory()

	var (
		fn func() ([]byte, error)
		p  string
	)

	for _, file := range files {
		p = path.Join(dirpath, file.Name())

		if file.IsDir() {
			fn = func() ([]byte, error) { return NewSerial(ctx, p, factory) }
		} else {
			fn = func() ([]byte, error) { return calculateFileDigest(p, factory) }
		}

		d, err := fn()
		if err != nil {
			return nil, err
		}

		hasher.Write(d)
	}

	return hasher.Sum(nil), nil
}
