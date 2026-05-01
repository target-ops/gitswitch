// Package main is the gitswitch entrypoint.
//
// gitswitch (Go rewrite, Phase 1) — see README.md and the
// `next` branch PR for context. The Python implementation in
// src/ remains the canonical 0.x release on `main` until 1.0
// ships from this branch.
package main

import (
	"fmt"
	"os"

	"github.com/target-ops/gitswitch/internal/cmd"
)

// Version is injected at build time via -ldflags "-X main.Version=...".
var Version = "1.0.0-dev"

func main() {
	root := cmd.NewRootCommand(Version)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
