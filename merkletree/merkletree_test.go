package merkletree_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chrisdoherty4/merklehash/merkletree"
)

func TestDirectoryHashOfFilesOnly(t *testing.T) {
	testdir, err := filepath.Abs(filepath.Join("testdata", "test-1"))
	require.NoError(t, err)

	hash, err := merkletree.New(context.Background(), testdir, sha256.New)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"b90b3ade27194c7a225346df4764aee1dffc3c5cd57558ad5f4f63299f12b615",
		fmt.Sprintf("%x", hash),
	)
}

func TestDirectoryWithSubDir(t *testing.T) {
	testdir, err := filepath.Abs(filepath.Join("testdata", "test-2"))
	require.NoError(t, err)

	hash, err := merkletree.New(context.Background(), testdir, sha256.New)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"b734b659c21e1aed8ac690f4be08e03aa337dcf06a08caaf474231cbd14989dc",
		fmt.Sprintf("%x", hash),
	)
}

func TestDirectoryWithMultipleSubDirs(t *testing.T) {
	testdir, err := filepath.Abs(filepath.Join("testdata", "test-3"))
	require.NoError(t, err)

	hash, err := merkletree.New(context.Background(), testdir, sha256.New)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"a09d42fcf44f468f70e260fc318e62aaa6638ca20c0ab271c359566669e7e29d",
		fmt.Sprintf("%x", hash),
	)
}

func TestCanceledContext(t *testing.T) {
	testdir, err := filepath.Abs(filepath.Join("testdata", "test-3"))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = merkletree.New(ctx, testdir, sha256.New)
	assert.Error(t, err)
}

var (
	result  []byte
	testdir string
)

func init() {
	t, err := filepath.Abs(filepath.Join("testdata", "test-3"))
	if err != nil {
		panic(err)
	}
	testdir = t
}

func BenchmarkMerkletreeNew(b *testing.B) {
	var r []byte
	for n := 0; n < b.N; n++ {
		r, _ = merkletree.New(context.Background(), testdir, sha256.New)
	}
	result = r
}

func BenchmarkMerkletreeNewSerial(b *testing.B) {
	var r []byte
	for n := 0; n < b.N; n++ {
		r, _ = merkletree.NewSerial(context.Background(), testdir, sha256.New)
	}
	result = r
}
