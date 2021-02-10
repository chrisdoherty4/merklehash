package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chrisdoherty4/merklehash/merkletree"
	"github.com/spf13/cobra"
)

var (
	// The long help for the command
	longHelp string

	// A map of supported algorithms.
	algs map[string]merkletree.HashFactory

	// MerkleHashCmd represents the root command for the merkle hasher.
	merkleHashCmd *cobra.Command

	// Optional argument
	algorithm string
)

func init() {
	algs = make(map[string]merkletree.HashFactory)
	algs = map[string]merkletree.HashFactory{
		"md5":        md5.New,
		"sha1":       sha1.New,
		"sha224":     sha256.New224,
		"sha256":     sha256.New,
		"sha384":     sha512.New384,
		"sha512":     sha512.New,
		"sha512/224": sha512.New512_224,
		"sha512/256": sha512.New512_256,
	}

	algsSlice := []string{}
	for alg := range algs {
		algsSlice = append(algsSlice, alg)
	}

	sort.Strings(algsSlice)
	algsSlice = mapSlice(algsSlice, func(v string) string {
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
		Use:     "merklehash <directory path>",
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

			fmt.Printf("%x\n", hash)
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

func mapSlice(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
