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
	if slices.Contains(os.Args, "--in-place") {
		return inplace(gomodName, updatedContents)
	}

	// print updated contents to stdout
	fmt.Println(string(updatedContents))
	return nil
}

func inplace(modLocation string, updatedContents []byte) error {
	// get current file info
	info, err := os.Stat(modLocation)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// write updated contents to go.mod
	if err = os.WriteFile(modLocation, updatedContents, info.Mode()); err != nil {
		return fmt.Errorf("failed to write updated go.mod: %w", err)
	}
	
	return nil
}
