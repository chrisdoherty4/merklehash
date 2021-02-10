package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chrisdoherty4/merklehash/merkletree"
)

var (
	// The long help for the command
	longHelp string

	// A map of supported algorithms.
	// todo add
	//		crc32.New
	//		crc64.New
	algs map[string]merkletree.HashFactory = map[string]merkletree.HashFactory{
		"md5":        md5.New,
		"sha1":       sha1.New,
		"sha224":     sha256.New224,
		"sha256":     sha256.New,
		"sha384":     sha512.New384,
		"sha512":     sha512.New,
		"sha512/224": sha512.New512_224,
		"sha512/256": sha512.New512_256,
		"fnv/128":    fnv.New128,
		"fnv/128a":   fnv.New128a,
		"fnv/32":     func() hash.Hash { return fnv.New32() },
		"fnv/32a":    func() hash.Hash { return fnv.New32a() },
		"fnv/64":     func() hash.Hash { return fnv.New64() },
		"fnv/64a":    func() hash.Hash { return fnv.New64a() },
		"crc32ieee":  func() hash.Hash { return crc32.NewIEEE() },
		"adler32":    func() hash.Hash { return adler32.New() },
	}

	// MerkleHashCmd represents the root command for the merkle hasher.
	merkleHashCmd *cobra.Command

	// Optional argument
	algorithm string
)

func init() {
	algsSlice := []string{}
	for alg := range algs {
		algsSlice = append(algsSlice, alg)
	}

	sort.Strings(algsSlice)
	algsSlice = mapStrings(algsSlice, func(v string) string {
		return fmt.Sprintf("  * %v", v)
	})

	longHelp = fmt.Sprintf(`
merklehash is a hashing tool for generating digests of arbitrary depth directory
hierarchies.

Supported algorithms include:
%v
`,
		strings.Join(algsSlice, "\n"),
	)

	merkleHashCmd = &cobra.Command{
		Use:     "merklehash [options] <directory path>",
		Short:   "merklehash is a directory hasher.",
		Long:    longHelp,
		Example: "  merklehash /path/to/directory",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			factory, ok := algs[algorithm]

			if !ok {
				fmt.Fprintf(os.Stderr, "Unsupported algorithm %v", algorithm)
				os.Exit(1)
			}

			hash, err := merkletree.New(
				context.Background(),
				args[0],
				factory,
			)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			fmt.Printf("%v %x\n", algorithm, hash)
		},
	}

	merkleHashCmd.Flags().StringVarP(
		&algorithm,
		"algorithm",
		"a",
		"sha256",
		"hashing algorithm",
	)
}

func mapStrings(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
