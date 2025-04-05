package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const gomodName = "go.mod"

var ErrFailedToParseFlags = fmt.Errorf("failed to parse flags")

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		if errors.Is(err, ErrFailedToParseFlags) {
			os.Exit(2)
		}

		log.Fatal(err)
	}
}

func run(args []string, output io.Writer) error {
	updatedContents, err := MergeRequires(gomodName)
	if err != nil {
		return fmt.Errorf("failed to merge requires: %w", err)
	}

	var inplace bool

	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	fs.SetOutput(output)

	fs.BoolVar(&inplace, "in-place", false, "replace the contents of go.mod with the updated contents")
	fs.BoolVar(&inplace, "i", false, "replace the contents of go.mod with the updated contents")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToParseFlags, err)
	}

	// check if we want to replace the contents of go.mod
	if inplace {
		return updateInplace(gomodName, updatedContents)
	}

	fmt.Println(string(updatedContents))
	return nil
}

func updateInplace(modLocation string, updatedContents []byte) error {
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
