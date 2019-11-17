package main

import (
	"fmt"
	"os"

	"github.com/chrisdoherty4/merklehash/pkg/merkletree"
	"github.com/spf13/cobra"
)

var (
	algorithm string

	// MerkleHashCmd represents the root command for the merkle hasher.
	MerkleHashCmd = &cobra.Command{
		Use:   "merklehash <directory path>",
		Short: "MerkleHash is a directory hasher.",
		Long: `MerkleHash is a hashing tool for generating digests of arbitrary depth 
directory hierarchies.

MerkleHash will output the hash of the directory followed by the directory that 
was hashed.`,
		Example: "  merklehash /path/to/directory",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !merkletree.Supports(algorithm) {
				fmt.Fprintf(os.Stdout, "'%v' is not a valid algorithm. Value algorithms are:\n", algorithm)
				for _, ident := range merkletree.GetAlgorithms() {
					fmt.Printf("  %v\n", ident)
				}
				os.Exit(0)
			}

			fmt.Printf(
				"%x\n",
				merkletree.New(args[0], merkletree.GetHasher(algorithm)).Hash(),
			)
		},
	}

	supportsCmd = &cobra.Command{
		Use:   "algorithms",
		Short: "List supported algorithms.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, ident := range merkletree.GetAlgorithms() {
				fmt.Printf("%v\n", ident)
			}
		},
	}
)

func init() {
	MerkleHashCmd.Flags().StringVarP(&algorithm, "algorithm", "a", "sha256", "Hashing algorithm to use")
	MerkleHashCmd.AddCommand(supportsCmd)
}
