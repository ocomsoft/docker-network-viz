// Package main provides the entry point for the docker-network-viz CLI tool.
// This tool visualizes Docker network topology in a tree-style format,
// showing networks, containers, and their connectivity relationships.
package main

import (
	"fmt"
	"os"

	"git.o.ocom.com.au/go/docker-network-viz/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
