package merkletree

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData string

var testDataDir testData

func (td testData) path(p string) string {
	return path.Join(string(td), p)
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	testDataDir = testData(path.Join(path.Dir(filename), "testdata"))
}

func TestDirectoryHashOfFilesOnly(t *testing.T) {
	hash, err := New(context.Background(), testDataDir.path("test-1"), sha256.New)

	assert.Nil(t, err, err)
	assert.Equal(
		t,
		"b90b3ade27194c7a225346df4764aee1dffc3c5cd57558ad5f4f63299f12b615",
		fmt.Sprintf("%x", hash),
	)
}

func TestDirectoryWithSubDir(t *testing.T) {
	hash, err := New(context.Background(), testDataDir.path("test-2"), sha256.New)

	assert.Nil(t, err, err)
	assert.Equal(
		t,
		"85d46d34fb3e3ac1f1b86a9a9ab3a20b6aa5c9b961ad93da5a9723d4c2fc3029",
		fmt.Sprintf("%x", hash),
	)
}

func TestDirectoryWithMultipleSubDirs(t *testing.T) {
	hash, err := New(context.Background(), testDataDir.path("test-3"), sha256.New)

	assert.Nil(t, err, err)
	assert.Equal(
		t,
		"8c2888c54d6c8119d5d08da63fa05bb0db8787f5e06b25e546e61fff090b82f6",
		fmt.Sprintf("%x", hash),
	)
}

var result []byte

func BenchmarkMerkleTree(b *testing.B) {
	var r []byte
	for n := 0; n < b.N; n++ {
		r, _ = New(context.Background(), testDataDir.path("test-3"), sha256.New)
	}
	result = r
}
