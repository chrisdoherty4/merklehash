package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chrisdoherty4/merklehash/pkg/merkletree"
)

func main() {
	// Define the command line interface
	alg := flag.String(
		"alg",
		"sha256",
		"Hashing algorithm to use with the merkle tree.",
	)

	raw := flag.Bool("raw", false, "Print only the hex hash.")

	// TODO: Tidy up help comment adding supported algorithms.
	flag.Parse()

	if flag.NArg() != 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// TODO: Change to use the merkletree.Supports function
	if !merkletree.Supports(*alg) {
		fmt.Fprintf(os.Stdout, "'%v' is not a valid algorithm. Value algorithms are:\n", *alg)
		for ident := range merkletree.GetAlgorithms() {
			fmt.Fprintf(os.Stdout, "  %v\n", ident)
		}
		os.Exit(0)
	}

	// Create a new merkle hash and output the hex representation of it's hash.
	fmt.Fprintf(
		os.Stdout,
		"%x",
		merkletree.New(flag.Arg(0), merkletree.GetHasher(*alg)).Hash(),
	)

	if !*raw {
		fmt.Fprintf(os.Stdout, " %v", flag.Arg(0))
	}
	println()
}
