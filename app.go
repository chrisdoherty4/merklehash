package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chrisdoherty4/gomerkle/merklehash"
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

	hasher := merklehash.Algorithms.GetHasher(*alg)

	if hasher == nil {
		fmt.Fprintf(os.Stdout, "'%v' is not a valid algorithm. Value algorithms are:\n", *alg)
		for e := merklehash.Algorithms.Front(); e != nil; e = e.Next() {
			fmt.Fprintf(os.Stdout, "  %v\n", e.Value.(*merklehash.Algorithm).Ident)
		}
		os.Exit(0)
	}

	// TODO: Add support for multiple directories.
	// TODO: Protection against huge file systems?
	// TODO: Resolve directory symlinks when going through directories.

	// Create a new merkle hash and output the hex representation of it's hash.
	fmt.Fprintf(os.Stdout, "%x", merklehash.New(flag.Arg(0), hasher).Hash())

	if !*raw {
		fmt.Fprintf(os.Stdout, " %v", flag.Arg(0))
	}
	println()
}