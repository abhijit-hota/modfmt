package main

import (
	"fmt"
	"log"
	"os"
	"slices"
)

const gomodName = "go.mod"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	updatedContents, err := MergeRequires(gomodName)
	if err != nil {
		return fmt.Errorf("failed to merge requires: %w", err)
	}

	// check if we want to replace the contents of go.mod
	if !slices.Contains(os.Args, "--replace") {
		// replace not found, so print updated contents to stdout
		fmt.Println(string(updatedContents))
		return nil
	}

	// get current file info
	info, err := os.Stat(gomodName)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// write updated contents to go.mod
	if err = os.WriteFile(gomodName, updatedContents, info.Mode()); err != nil {
		return fmt.Errorf("failed to write updated go.mod: %w", err)
	}

	return nil
}
