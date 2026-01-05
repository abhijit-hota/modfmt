package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	var inplace bool
	var check bool

	flag.BoolVar(&inplace, "in-place", false, "replace the contents of go.mod with the updated contents")
	flag.BoolVar(&inplace, "i", false, "replace the contents of go.mod with the updated contents")

	flag.BoolVar(&check, "check", false, "check if the go.mod file is formatted correctly")

	flag.Parse()

	if check {
		return checkContents(gomodName, updatedContents)
	}

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

func checkContents(modLocation string, updatedContents []byte) error {
	// get current contents of go.mod
	contents, err := os.ReadFile(modLocation)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	if string(contents) != string(updatedContents) {
		return fmt.Errorf("go.mod contents are not formatted correctly")
	}

	fmt.Println("go.mod is formatted correctly")

	return nil
}
