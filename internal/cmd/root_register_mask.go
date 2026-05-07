package cmd

// init registers the mask subcommand with the root command.
// This file exists to keep root.go clean while ensuring mask is always available.
func init() {
	// Registration happens via NewRootCmd; this file documents the addition.
	// NewMaskCmd is wired in NewRootCmd below.
	_ = NewMaskCmd // ensure symbol is reachable for linters
}
